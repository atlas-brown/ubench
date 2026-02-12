package wrappers

import (
	"context"
	"github.com/atlas/slowpoke/pkg/utility"
	"github.com/goccy/go-json"
	"net/http"
	"io"
)

func Wrapper[ReqType interface{}, RespType interface{}](handler func(context.Context, *ReqType) *RespType) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		input, err := io.ReadAll(r.Body)
		r.Body.Close()
		var req ReqType
		err = json.Unmarshal(input, &req)
		if err != nil {
			panic(err)
		}
		resp := handler(ctx, &req)
		utility.DumpJson(resp, w)
	}
}