package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	client "task2/internal/client/http"
	"task2/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	StatusCodeNotEqual = "HTTP status code должен быть %d"
)

// CreateAdvertisement используется для отправки валидных данных
func CreateAdvertisement(t *testing.T, expectedStatusCode int, request models.CreateAdvertisementRequest) string {
	createAdvertisementRawResponse := client.HttpCreateAdvertisement(t, request)

	assert.Equalf(t, expectedStatusCode, createAdvertisementRawResponse.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, createAdvertisementRawResponse))
}

// CreateAdvertisementFromMap используется для отправки невалидных или частично заполненных данных
func CreateAdvertisementFromMap(t *testing.T, expectedStatusCode int, req map[string]interface{}) string {
	createAdvertisementRawResponse := client.HttpCreateAdvertisement(t, req)

	assert.Equalf(t, expectedStatusCode, createAdvertisementRawResponse.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, createAdvertisementRawResponse))
}

/*
GetAdvertisementById получает объявление по id и гарантирует, что массив содержит
ровно один элемент согласно бизнес-логике
*/
func GetAdvertisementById(t *testing.T, expectedStatusCode int, id string) models.AdvertisementResponse {
	getAdvertisementRawResponse := client.HttpGetAdvertisement(t, id)

	assert.Equalf(t, expectedStatusCode, getAdvertisementRawResponse.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	var getResponse []models.AdvertisementResponse
	err := json.Unmarshal(readBody(t, getAdvertisementRawResponse), &getResponse)
	assert.NoError(t, err)

	assert.Len(t, getResponse, 1, "В массиве больше одного элемента")

	return getResponse[0]
}

// GetAdvertisementByIdExpectError используется для проверки ошибок
func GetAdvertisementByIdExpectError(t *testing.T, expectedStatusCode int, id string) string {
	getAdvertisementRawResponse := client.HttpGetAdvertisement(t, id)

	assert.Equalf(t, expectedStatusCode, getAdvertisementRawResponse.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, getAdvertisementRawResponse))
}

// DeleteAdvertisementById удаляет объявление по id и проверяет статус ответа
func DeleteAdvertisementById(t *testing.T, expectedStatusCode int, id string) string {
	deleteRawResponse := client.HttpDeleteAdvertisement(t, id)

	assert.Equalf(t, expectedStatusCode, deleteRawResponse.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, deleteRawResponse))
}

// GetSellerAdvertisementsBySellerId успешно получает все объявления конкретного продавца по sellerId
func GetSellerAdvertisementsBySellerId(t *testing.T, expectedStatusCode int, sellerId int) []models.AdvertisementResponse {
	getAdvertisementsRaw := client.HttpGetSellerAdvertisements(t, sellerId)

	assert.Equalf(t, expectedStatusCode, getAdvertisementsRaw.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	var getAdvertisementsResponse []models.AdvertisementResponse

	err := json.Unmarshal(readBody(t, getAdvertisementsRaw), &getAdvertisementsResponse)
	assert.NoError(t, err)

	return getAdvertisementsResponse
}

// GetSellerAdvertisementsExpectError - метод для возвращения ошибки получения объявлений продавца.
// Используется interface{}, так как ожидается ошибка как при строковом значении, так и при числовом
func GetSellerAdvertisementsExpectError(t *testing.T, expectedStatusCode int, sellerId interface{}) string {
	getAdvertisementsRaw := client.HttpGetSellerAdvertisements(t, sellerId)

	assert.Equalf(t, expectedStatusCode, getAdvertisementsRaw.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, getAdvertisementsRaw))
}

/*
GetStatisticAdvertisementById получает статистику по id.
Проверяет статус ответа, десериализует в Statistics и гарантирует, что массив содержит
ровно один элемент согласно бизнес-логике
*/
func GetStatisticAdvertisementById(t *testing.T, expectedStatusCode int, id string) models.Statistics {
	getStatisticResponseRaw := client.HttpGetStatisticByAdsId(t, id)

	assert.Equalf(t, expectedStatusCode, getStatisticResponseRaw.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	var getStatistic []models.Statistics

	err := json.Unmarshal(readBody(t, getStatisticResponseRaw), &getStatistic)
	assert.NoError(t, err)

	assert.NotZero(t, len(getStatistic), "Массив в ответе пустой")
	assert.Len(t, getStatistic, 1, "В массиве больше одного элемента")

	return getStatistic[0]
}

// GetStatisticAdvertisementByIdExpectError получает статистику объявления по id и ожидает ошибку
func GetStatisticAdvertisementByIdExpectError(t *testing.T, expectedStatusCode int, id string) string {
	getStatisticResponseRaw := client.HttpGetStatisticByAdsId(t, id)

	assert.Equalf(t, expectedStatusCode, getStatisticResponseRaw.StatusCode,
		StatusCodeNotEqual, expectedStatusCode)

	return string(readBody(t, getStatisticResponseRaw))
}

func readBody(t *testing.T, resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return body
}
