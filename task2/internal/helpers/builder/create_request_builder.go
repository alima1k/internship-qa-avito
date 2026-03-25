package builder

import "task2/internal/models"

/*
CreateAdvertisementBuilder для удобной работы с CreateAdvertisementRequest
*/
type CreateAdvertisementBuilder struct {
	request models.CreateAdvertisementRequest
}

func NewCreateAdvertisementBuilder() *CreateAdvertisementBuilder {
	positiveNumber := 10
	sellerId := 333333
	advName := "Машина"
	return &CreateAdvertisementBuilder{
		request: models.CreateAdvertisementRequest{
			SellerId: sellerId,
			Name:     advName,
			Price:    positiveNumber,
			Statistics: models.Statistics{
				Likes:     positiveNumber,
				ViewCount: positiveNumber,
				Contacts:  positiveNumber,
			},
		}}
}

func (b *CreateAdvertisementBuilder) WithSellerId(sellerId int) *CreateAdvertisementBuilder {
	b.request.SellerId = sellerId
	return b
}

func (b *CreateAdvertisementBuilder) WithName(name string) *CreateAdvertisementBuilder {
	b.request.Name = name
	return b
}

func (b *CreateAdvertisementBuilder) WithPrice(price int) *CreateAdvertisementBuilder {
	b.request.Price = price
	return b
}

func (b *CreateAdvertisementBuilder) WithLikes(likes int) *CreateAdvertisementBuilder {
	b.request.Statistics.Likes = likes
	return b
}

func (b *CreateAdvertisementBuilder) WithViewCount(viewCount int) *CreateAdvertisementBuilder {
	b.request.Statistics.ViewCount = viewCount
	return b
}

func (b *CreateAdvertisementBuilder) WithContacts(contacts int) *CreateAdvertisementBuilder {
	b.request.Statistics.Contacts = contacts
	return b
}

func (b *CreateAdvertisementBuilder) Build() models.CreateAdvertisementRequest {
	return b.request
}
