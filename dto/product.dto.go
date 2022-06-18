package dto

import "encoding/json"

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryId int     `json:"categoryId" binding:"required"`
}

type GetProductsRequest struct {
	Name       string `form:"name"`
	CategoryId int    `form:"categoryId"`
}

type BreadCrumb struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b BreadCrumb) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(b)
	return bytes, err
}

type GetProductResponse struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	BreadCrumbs []*BreadCrumb `json:"breadCrumbs"`
}
