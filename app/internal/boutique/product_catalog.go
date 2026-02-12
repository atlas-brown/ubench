package boutique

import (
	"context"
	"strings"
	"fmt"
)

const (
	debug_product_catalog = false
)

var allProducts []Product

var CatalogSize = 1000

func InitAllProducts(ctx context.Context, products []Product) {
	allProducts = products
	fmt.Println("InitAllProducts: ", len(allProducts))
}

func GetAllProducts(ctx context.Context) {
	for _, product := range allProducts {
		fmt.Println(product)
	}
}

func AddProduct(ctx context.Context, product Product) string {
	allProducts = append(allProducts, product)
	return product.Id
}

func AddProducts(ctx context.Context, products []Product) {
	if len(allProducts) < CatalogSize {
		rest := CatalogSize - len(allProducts)
		if len(products) < rest {
			rest = len(products)
		}
		for i := 0; i < rest; i++ {
			allProducts = append(allProducts, products[i])
		}
	}
	return
}

func GetProduct(ctx context.Context, Id string) Product {
	if debug_product_catalog { fmt.Println("GetProduct: ", Id) }

	var product Product
	for _, p := range allProducts {
		if p.Id == Id {
			product = p
			break
		}
	}
	return product
}

func SearchProducts(ctx context.Context, name string) []Product {
	if debug_product_catalog { fmt.Println("SearchProducts: ", name) }

	var products []Product
	for _, p := range allProducts {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(name)) {
			products = append(products, p)
		}
	}
	return products
}

func FetchCatalog(ctx context.Context, catalogSize int) []Product {
	if debug_product_catalog { fmt.Println("FetchCatalog: ", catalogSize) }
	var products []Product
	if catalogSize < len(allProducts) {
		products = allProducts[:catalogSize]
	} else {
		products = allProducts
	}
	return products
}
