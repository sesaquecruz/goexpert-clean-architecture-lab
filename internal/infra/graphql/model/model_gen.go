// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CrateOrderInput struct {
	Price float64 `json:"Price"`
	Tax   float64 `json:"Tax"`
}

type Order struct {
	ID         string  `json:"Id"`
	Price      float64 `json:"Price"`
	Tax        float64 `json:"Tax"`
	FinalPrice float64 `json:"FinalPrice"`
}
