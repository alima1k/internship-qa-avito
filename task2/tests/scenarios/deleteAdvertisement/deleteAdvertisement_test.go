package deleteAdvertisement

import (
	"fmt"
	"task2/internal/constants/statusCode"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
	"task2/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	createdId string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupTest() {
	dto := builder.NewCreateAdvertisementBuilder().Build()
	createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
	id := utils.ExtractUUIDFromResponse(createdResponse)
	if id == "" {
		s.T().Fatal("Не удалось получить UUID")
	}
	s.createdId = id
}

// Проверяем успешное удаление объявления
// Тест-кейс №55
func (s *TestSuite) TestDeleteAdvertisementPositive() {
	s.Run("Успешное удаление существующего объявления", func() {
		deleteResponse := advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		assert.Empty(s.T(), deleteResponse, "Тело ответа не пустое")
	})
}

// Проверяем 400 ошибку при попытке в pathVariable вместо UUID вписать строку
// Тест-кейс №56
func (s *TestSuite) TestDeleteAdvertisementInvalidIdFormat() {
	s.Run("Ошибка 400 при удалении объявления по невалидному id", func() {
		invalidId := utils.RandomString(5)
		errorResponse := advertisements.DeleteAdvertisementById(s.T(), statusCode.BadRequest, invalidId)
		asserts.AssertResponseError(s.T(), errorResponse, "400", "переданный id айтема некорректный")
	})
}

// Проверяем ошибку 404 при попытке удаления объявления по рандомному UUID
// Тест-кейс №57
// Bugs: Bug-8
func (s *TestSuite) TestDeleteAdvertisementByNonExistingId() {
	s.Run("Ошибка 404 при удалении объявления по несуществующему id", func() {
		id := uuid.NewString()
		errorResponse := advertisements.DeleteAdvertisementById(s.T(), statusCode.NotFound, id)
		asserts.AssertResponseError(s.T(), errorResponse, "404", fmt.Sprintf("item %s not found", id))
	})
}

// Проверяем, что удаление второй раз не меняет состояние системы
// Тест-кейс №58
// Bugs: Bug-8
func (s *TestSuite) TestDeleteAdvertisementIdempotency() {
	s.Run("Удаление объявления с 200 в первый раз", func() {
		deleteResponse := advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		assert.Empty(s.T(), deleteResponse, "Тело ответа успешного удаления не пустое")
	})
	s.Run("Удаление объявления с 404 во второй раз", func() {
		deleteResponse := advertisements.DeleteAdvertisementById(s.T(), statusCode.NotFound, s.createdId)
		asserts.AssertResponseError(s.T(), deleteResponse, "404", "")
	})
}
