package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/boutique"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net"
	"net/http"
	"runtime"
	"os"
	"github.com/goccy/go-json"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func addProduct(ctx context.Context, req *boutique.AddProductRequest) *boutique.AddProductResponse {
	productId := boutique.AddProduct(ctx, req.Product)
	resp := boutique.AddProductResponse{ProductId: productId}
	return &resp
}

func getProduct(ctx context.Context, req *boutique.GetProductRequest) *boutique.GetProductResponse {
	product := boutique.GetProduct(ctx, req.ProductId)
	resp := boutique.GetProductResponse{Product: product}
	return &resp
}

func searchProducts(ctx context.Context, req *boutique.SearchProductsRequest) *boutique.SearchProductsResponse {
	products := boutique.SearchProducts(ctx, req.Query)
	resp := boutique.SearchProductsResponse{Products: products}
	return &resp
}

func fetchCatalog(ctx context.Context, req *boutique.FetchCatalogRequest) *boutique.FetchCatalogResponse {
	products := boutique.FetchCatalog(ctx, req.CatalogSize)
	resp := boutique.FetchCatalogResponse{Catalog: products}
	return &resp
}

func addProducts(ctx context.Context, req *boutique.AddProductsRequest) *boutique.AddProductsResponse {
	boutique.AddProducts(ctx, req.Products)
	resp := boutique.AddProductsResponse{OK: "OK"}
	return &resp
}

func loadProducts(ctx context.Context) []boutique.Product {

	// List directory
	var products []boutique.Product
	catalogJSON, err := os.ReadFile("/app/cmd/boutique/products.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(catalogJSON, &products)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded products: ", len(products))
	return products
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/add_product", wrappers.Wrapper[boutique.AddProductRequest, boutique.AddProductResponse](addProduct))
	http.HandleFunc("/add_products", wrappers.Wrapper[boutique.AddProductsRequest, boutique.AddProductsResponse](addProducts))
	http.HandleFunc("/ro_get_product", wrappers.Wrapper[boutique.GetProductRequest, boutique.GetProductResponse](getProduct))
	http.HandleFunc("/ro_search_products", wrappers.Wrapper[boutique.SearchProductsRequest, boutique.SearchProductsResponse](searchProducts))
	http.HandleFunc("/ro_fetch_catalog", wrappers.Wrapper[boutique.FetchCatalogRequest, boutique.FetchCatalogResponse](fetchCatalog))
	boutique.InitAllProducts(context.Background(), loadProducts(context.Background()))
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
