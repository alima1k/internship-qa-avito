package models

type AdvertisementResponse struct {
	Id         string     `json:"id"`
	SellerId   int        `json:"sellerId"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Statistics Statistics `json:"statistics"`
	CreatedAt  string     `json:"createdAt"`
}

type ErrorResponse struct {
	Result Result `json:"result"`
	Status string `json:"status"`
}

type Result struct {
	Messages map[string]string `json:"messages"`
	Message  string            `json:"message"`
}
