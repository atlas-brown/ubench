package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/atlas/slowpoke/pkg/invoke"
	"github.com/atlas/slowpoke/pkg/synthetic"
	"google.golang.org/grpc/metadata"

	// "github.com/goccy/go-json"
	"math/rand"
	"sync"
)

func pickDynamicService(calledServices []synthetic.CalledService) string {

	type s_p_pair struct {
		service_name string
		probability  int
	}

	sum_prob := 0
	var service_prob []s_p_pair
	for _, service := range calledServices {
		if service.Probability != 0 {
			sum_prob += service.Probability
			service_prob = append(service_prob, s_p_pair{service.Service, sum_prob})
		}
	}

	// if only one service is available, add a dummy service to make the service prob as p/100
	// A single dynamic service case
	if len(service_prob) == 1 {
		sum_prob = 100
		service_prob = append(service_prob, s_p_pair{"", sum_prob})
	}

	// Dynamic pattern: randomly pick one value
	picked_service := ""
	if len(service_prob) != 0 {
		rand_value := rand.Intn(sum_prob-1) + 1
		for _, p := range service_prob {
			if rand_value <= p.probability {
				picked_service = p.service_name
				break
			}
		}
	}
	return picked_service
}

func getRequestIDFromIncomingGRPC(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	vals := md.Get("x-request-id")
	if len(vals) == 0 {
		return "", false
	}
	return vals[0], true
}

func execParallel(calledServices []synthetic.CalledService, request *http.Request, grpc_ctx *context.Context) map[string]string {

	// pick dynamic service
	picked_service := pickDynamicService(calledServices)

	// forward requests
	respMap := make(map[string]string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, service := range calledServices {
		if service.Probability != 0 && service.Service != picked_service {
			continue
		}
		for i := 0; i < service.TrafficForwardRatio; i++ {
			wg.Add(1)
			go func(service synthetic.CalledService) {
				defer wg.Done()
				var resp string
				if service.Protocol == "grpc" {
					// 1) Choose a base ctx safely
					var ctx context.Context
					if grpc_ctx != nil && *grpc_ctx != nil {
						ctx = *grpc_ctx
					} else if request != nil {
						ctx = request.Context()
					} else {
						ctx = context.Background()
					}

					// 2) Propagate x-request-id (prefer incoming gRPC -> fallback HTTP)
					if rid, ok := getRequestIDFromIncomingGRPC(ctx); ok {
						ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
					} else if request != nil {
						if rid := request.Header.Get("x-request-id"); rid != "" {
							ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
						}
					}
					resp = invoke.InvokeGRPC(ctx, service.Service, service.Endpoint, "")
				} else {
					respRaw := invoke.Invoke[Response](request.Context(), service.Service, service.Endpoint, "", *request)
					resp = fmt.Sprintf("CPUResp: %s | NETResp : %s", respRaw.CPUResp, respRaw.NetworkResp)
				}
				key := fmt.Sprintf("%s [%s,%d]", service.Service, service.Endpoint, i)
				mu.Lock()
				respMap[key] = resp
				mu.Unlock()
			}(service)
		}
	}
	wg.Wait()
	return respMap
}

func execSequential(calledServices []synthetic.CalledService, request *http.Request, grpc_ctx *context.Context) map[string]string {

	// pick dynamic service
	picked_service := pickDynamicService(calledServices)

	// forward requests
	respMap := make(map[string]string)
	for _, service := range calledServices {
		if service.Probability != 0 && service.Service != picked_service {
			continue
		}
		for i := 0; i < service.TrafficForwardRatio; i++ {
			var resp string
			if service.Protocol == "grpc" {
				// 1) Choose a base ctx safely
				var ctx context.Context
				if grpc_ctx != nil && *grpc_ctx != nil {
					ctx = *grpc_ctx
				} else if request != nil {
					ctx = request.Context()
				} else {
					ctx = context.Background()
				}

				// 2) Propagate x-request-id (prefer incoming gRPC -> fallback HTTP)
				if rid, ok := getRequestIDFromIncomingGRPC(ctx); ok {
					ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
				} else if request != nil {
					if rid := request.Header.Get("x-request-id"); rid != "" {
						ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
					}
				}
				resp = invoke.InvokeGRPC(ctx, service.Service, service.Endpoint, "")
			} else {
				respRaw := invoke.Invoke[Response](request.Context(), service.Service, service.Endpoint, "", *request)
				resp = fmt.Sprintf("CPUResp: %s | NETResp : %s", respRaw.CPUResp, respRaw.NetworkResp)
			}
			key := fmt.Sprintf("%s [%s,%d]", service.Service, service.Endpoint, i)
			respMap[key] = resp
		}
	}
	return respMap
}

func execNetwork(request *http.Request, grpc_ctx *context.Context, endpoint *synthetic.Endpoint) map[string]string {
	if endpoint.NetworkComplexity == nil {
		return map[string]string{"nil": "No network complexity"}
	}

	if endpoint.NetworkComplexity.ForwardRequests == "asynchronous" {
		return execParallel(endpoint.NetworkComplexity.CalledServices, request, grpc_ctx)
	} else {
		return execSequential(endpoint.NetworkComplexity.CalledServices, request, grpc_ctx)
	}
}
