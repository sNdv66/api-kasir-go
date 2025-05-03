package services

import (
    "api-kasir/utils"
    "encoding/json"
    "errors"
)

type Product struct {
	ID     string `json:"id"`
	Stock  int    `json:"stock"`
	Name   string `json:"name"`
}

func GetProductByID(productID string) (Product, error) {
	req, err := utils.NewRequest("GET", "products?id=eq."+productID, nil)
	if err != nil {
		return Product{}, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return Product{}, err
	}

	if len(products) == 0 {
		return Product{}, errors.New("product not found")
	}

	return products[0], nil
}