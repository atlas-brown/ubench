package invoke

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	pb "github.com/atlas/slowpoke/pkg/pb"
	"github.com/atlas/slowpoke/pkg/utility"
	"google.golang.org/grpc"
)

var (
	grpcConns    map[string]*grpc.ClientConn = make(map[string]*grpc.ClientConn)
	grpcConnLock sync.RWMutex
)

var HTTPClient = &http.Client{
	Transport: &http.Transport{MaxConnsPerHost: 0, MaxIdleConnsPerHost: 500, MaxIdleConns: 500},
	Timeout:   60 * time.Second,
}

func performRequest[T interface{}](ctx context.Context, req *http.Request, res *T, app string, method string, argBytes []byte) {
	resp, err := HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	utility.Assert(resp.StatusCode == http.StatusOK)
	defer resp.Body.Close()
	utility.ParseJson(resp.Body, res)
}

func Invoke[T interface{}](ctx context.Context, app string, method string, input interface{}, request http.Request) T {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	var res T
	// Use kubernete native DNS addr
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/%s", app, "default", "80", method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	// Forward x-request-id if present
	if rid := request.Header.Get("x-request-id"); rid != "" {
		req.Header.Set("x-request-id", rid)
	}
	if err != nil {
		panic(err)
	}
	performRequest[T](ctx, req, &res, app, method, buf)
	return res
}

func InitGRPCConn(app string) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("%s.default.svc.cluster.local:%s", app, "80"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func InvokeGRPC(ctx context.Context, app string, method string, input interface{}) string {
	var grpcConn *grpc.ClientConn

	grpcConnLock.RLock()
	if _, ok := grpcConns[app]; ok {
		grpcConn = grpcConns[app]
	}
	grpcConnLock.RUnlock()

	if grpcConn == nil {
		grpcConnLock.Lock()
		if _, ok := grpcConns[app]; !ok {
			grpcConns[app] = InitGRPCConn(app)
		}
		grpcConn = grpcConns[app]
		grpcConnLock.Unlock()
	}

	client := pb.NewSimpleClient(grpcConn)
	resp_, err := client.SimpleRPC(ctx, &pb.SimpleRequest{Endpoint: method})
	if err != nil {
		panic(err)
	}
	return resp_.Resp
}
