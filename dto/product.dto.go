package dto

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryId int     `json:"categoryId" binding:"required"`
}

type GetProductsRequest struct {
	Name       string `form:"name"`
	CategoryId int    `form:"categoryId"`
}
