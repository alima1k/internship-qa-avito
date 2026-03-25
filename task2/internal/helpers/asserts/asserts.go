package asserts

import (
	"encoding/json"
	"task2/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ResponseMatchesDto проверяет, что поля объекта объявления совпадают с ожидаемым dto
func ResponseMatchesDto(t *testing.T, dto models.CreateAdvertisementRequest, response models.AdvertisementResponse) {
	assert.Equal(t, dto.SellerId, response.SellerId, "SellerId не совпадает")
	assert.Equal(t, dto.Name, response.Name, "Name не совпадает")
	assert.Equal(t, dto.Price, response.Price, "Price не совпадает")

	AssertStatisticFieldsMatch(t, dto.Statistics, response.Statistics)
}

// AssertStatisticFieldsMatch проверяет поля статистики объявления (contacts, likes, viewCount)
func AssertStatisticFieldsMatch(t *testing.T, dtoStats models.Statistics, responseStats models.Statistics) {
	assert.Equal(t, dtoStats.Contacts, responseStats.Contacts, "Contacts не совпадает")
	assert.Equal(t, dtoStats.Likes, responseStats.Likes, "Likes не совпадает")
	assert.Equal(t, dtoStats.ViewCount, responseStats.ViewCount, "ViewCount не совпадает")
}

// AssertResponseError проверяет, что ответ API с ошибкой имеет ожидаемый статус и сообщение.
// Возвращает объект ErrorResponse для дальнейших проверок
func AssertResponseError(t *testing.T, raw string, expectedStatus string, expectedError string) models.ErrorResponse {
	var errResp models.ErrorResponse

	assert.NoError(t, json.Unmarshal([]byte(raw), &errResp), "Неверный формат ответа")
	assert.Equal(t, expectedStatus, errResp.Status, "Показывается неверный статус")
	assert.Equal(t, expectedError, errResp.Result.Message, "Показывается неверное сообщение")

	return errResp
}

// AssertResponseFieldError проверяет, что в ответе API с ошибкой присутствуют корректные сообщения для каждого поля
// fieldErrorExpected - map вида "поле":"ожидаемое сообщение"
func AssertResponseFieldError(t *testing.T, raw string, expectedStatus string, expectedError string, fieldErrorExpected map[string]string) {
	errResponse := AssertResponseError(t, raw, expectedStatus, expectedError)

	for k, v := range fieldErrorExpected {
		assert.Equalf(t, v, errResponse.Result.Messages[k], "Поле статистики %s не соответствует ожидаемому", k)
	}
}
