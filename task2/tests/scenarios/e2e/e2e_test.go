package e2e

import (
	"task2/internal/constants/statusCode"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
	"task2/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	createdIds []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// В конце тест-сьюта удаляем все созданные объявления
func (s *TestSuite) TearDownSuite() {
	for _, id := range s.createdIds {
		advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, id)
	}
	s.createdIds = nil
}

/*
Получение объявления и статистики после создания
Тест-кейс №59
*/
func (s *TestSuite) TestCreateAndGetAdvertisementWithStatistics() {
	dto := builder.NewCreateAdvertisementBuilder().Build()
	var id string

	s.Run("Создание объявления с валидными полями", func() {
		createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
		id = utils.ExtractUUIDFromResponse(createdResponse)
		s.createdIds = append(s.createdIds, id)
	})

	s.Run("Успешный поиск по id", func() {
		getResponse := advertisements.GetAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.ResponseMatchesDto(s.T(), dto, getResponse)
	})

	s.Run("Валидная выдача статистики", func() {
		getStatsResponse := advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.AssertStatisticFieldsMatch(s.T(), dto.Statistics, getStatsResponse)
	})
}

/*
Получение у продавца созданного объявления
Тест-кейс №60
*/
func (s *TestSuite) TestCreateAdvertisementAppearsForSeller() {
	dto := builder.NewCreateAdvertisementBuilder().Build()
	var id string
	s.Run("Создание объявления с валидными полями", func() {
		createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
		id = utils.ExtractUUIDFromResponse(createdResponse)
		s.createdIds = append(s.createdIds, id)
	})
	s.Run("Успешный поиск по объявлениям продавца", func() {
		var exist bool
		var countAds int
		getSellerAdvResponse := advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, dto.SellerId)
		for _, ad := range getSellerAdvResponse {
			if id == ad.Id {
				exist = true
				countAds++
			}
		}
		assert.True(s.T(), exist, "Среди объявлений продавца не нашлось созданного")
		assert.Equal(s.T(), 1, countAds, "Объявление отображается больше 1 раза у продавца")
	})
}

/*
Получение удаленного объявления возвращает 404
Тест-кейс №61
*/
func (s *TestSuite) TestGetDeletedAdvertisementReturns404() {
	dto := builder.NewCreateAdvertisementBuilder().Build()
	var id string
	s.Run("Создание объявления с валидными полями", func() {
		createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
		id = utils.ExtractUUIDFromResponse(createdResponse)
	})
	s.Run("Успешный поиск по id", func() {
		getResponse := advertisements.GetAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.ResponseMatchesDto(s.T(), dto, getResponse)
	})
	s.Run("Удаление объявления", func() {
		body := advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, id)
		assert.Empty(s.T(), body, "Тело ответа не пустое")
	})
	s.Run("Ошибка 404 поиск по удаленному id", func() {
		advertisements.GetAdvertisementByIdExpectError(s.T(), statusCode.NotFound, id)
	})
}

/*
Получение статистики удаленного объявления возвращает 404
Тест-кейс №62
Bugs: Bug-7
*/
func (s *TestSuite) TestGetDeletedAdvertisementStatisticsReturns404() {
	dto := builder.NewCreateAdvertisementBuilder().Build()
	var id string
	s.Run("Создание объявления с валидными полями", func() {
		createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
		id = utils.ExtractUUIDFromResponse(createdResponse)
	})
	s.Run("Успешное получение статистики по id", func() {
		getStats := advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.AssertStatisticFieldsMatch(s.T(), dto.Statistics, getStats)
	})
	s.Run("Удаление объявления", func() {
		body := advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, id)
		assert.Empty(s.T(), body, "Тело ответа не пустое")
	})
	s.Run("Ошибка 404 получения статистики по удаленному id", func() {
		advertisements.GetStatisticAdvertisementByIdExpectError(s.T(), statusCode.NotFound, id)
	})
}
