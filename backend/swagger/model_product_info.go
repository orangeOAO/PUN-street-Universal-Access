/*
 * PUN street Universal Access - OpenAPI 3.0
 *
 * pua
 *
 * API version: v1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ProductInfo struct {
	ProductId int64 `json:"product_id"`

	Name string `json:"name"`

	StoreId int64 `json:"store_id"`

	Description string `json:"description"`

	Picture string `json:"picture"`

	Price int64 `json:"price"`

	Stock int64 `json:"stock"`

	Status int64 `json:"status"`
}
