package builder

/*
InvalidAdvertisementBuilder для работы с невалидными данными в реквесте. Для гибкости используется map
*/
type InvalidAdvertisementBuilder struct {
	data map[string]interface{}
}

const (
	str            = "string"
	positiveNumber = 10
	negativeNumber = -10
)

func NewInvalidAdvertisementBuilder() *InvalidAdvertisementBuilder {
	return &InvalidAdvertisementBuilder{
		data: map[string]interface{}{
			"sellerId": 333333,
			"name":     "Машина",
			"price":    positiveNumber,
			"statistics": map[string]interface{}{
				"likes":     positiveNumber,
				"viewCount": positiveNumber,
				"contacts":  positiveNumber,
			},
		},
	}
}

func (b *InvalidAdvertisementBuilder) SellerIdAsString() *InvalidAdvertisementBuilder {
	b.data["sellerId"] = str
	return b
}

func (b *InvalidAdvertisementBuilder) SellerIdAsNull() *InvalidAdvertisementBuilder {
	b.data["sellerId"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) WithoutSellerId() *InvalidAdvertisementBuilder {
	delete(b.data, "sellerId")
	return b
}

func (b *InvalidAdvertisementBuilder) WithSellerId(value int) *InvalidAdvertisementBuilder {
	b.data["sellerId"] = value
	return b
}

func (b *InvalidAdvertisementBuilder) PriceAsString() *InvalidAdvertisementBuilder {
	b.data["price"] = str
	return b
}

func (b *InvalidAdvertisementBuilder) PriceAsNull() *InvalidAdvertisementBuilder {
	b.data["price"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) WithoutPrice() *InvalidAdvertisementBuilder {
	delete(b.data, "price")
	return b
}

func (b *InvalidAdvertisementBuilder) PriceAsNegative() *InvalidAdvertisementBuilder {
	b.data["price"] = negativeNumber
	return b
}

func (b *InvalidAdvertisementBuilder) NameAsNumber() *InvalidAdvertisementBuilder {
	b.data["name"] = 12345
	return b
}

func (b *InvalidAdvertisementBuilder) NameAsNull() *InvalidAdvertisementBuilder {
	b.data["name"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) NameAsEmpty() *InvalidAdvertisementBuilder {
	b.data["name"] = ""
	return b
}

func (b *InvalidAdvertisementBuilder) LikesAsString() *InvalidAdvertisementBuilder {
	b.getStatistics()["likes"] = str
	return b
}

func (b *InvalidAdvertisementBuilder) LikesAsNull() *InvalidAdvertisementBuilder {
	stats := b.getStatistics()
	stats["likes"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) LikesAsNegative() *InvalidAdvertisementBuilder {
	b.getStatistics()["likes"] = negativeNumber
	return b
}

func (b *InvalidAdvertisementBuilder) ViewCountAsString() *InvalidAdvertisementBuilder {
	b.getStatistics()["viewCount"] = str
	return b
}

func (b *InvalidAdvertisementBuilder) ViewCountAsNull() *InvalidAdvertisementBuilder {
	stats := b.getStatistics()
	stats["viewCount"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) ViewCountAsNegative() *InvalidAdvertisementBuilder {
	b.getStatistics()["viewCount"] = negativeNumber
	return b
}

func (b *InvalidAdvertisementBuilder) ContactsAsString() *InvalidAdvertisementBuilder {
	b.getStatistics()["contacts"] = str
	return b
}

func (b *InvalidAdvertisementBuilder) ContactsAsNull() *InvalidAdvertisementBuilder {
	stats := b.getStatistics()
	stats["contacts"] = nil
	return b
}

func (b *InvalidAdvertisementBuilder) ContactsAsNegative() *InvalidAdvertisementBuilder {
	b.getStatistics()["contacts"] = negativeNumber
	return b
}

func (b *InvalidAdvertisementBuilder) getStatistics() map[string]interface{} {
	existingStats, exists := b.data["statistics"]
	if exists {
		if statsMap, ok := existingStats.(map[string]interface{}); ok {
			return statsMap
		}
	}
	newStats := make(map[string]interface{})
	b.data["statistics"] = newStats

	return newStats
}

func (b *InvalidAdvertisementBuilder) Build() map[string]interface{} {
	return b.data
}
