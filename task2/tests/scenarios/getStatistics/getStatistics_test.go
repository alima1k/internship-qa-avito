package getStatistics

import (
	"fmt"
	"sync"
	"task2/internal/constants/statusCode"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
	"task2/internal/models"
	"task2/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	dto       models.CreateAdvertisementRequest
	createdId string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	dto := builder.NewCreateAdvertisementBuilder().WithSellerId(utils.GenerateSellerId()).Build()
	id := utils.ExtractUUIDFromResponse(advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto))
	if id == "" {
		s.T().Fatal("Не удалось получить UUID")
	}
	s.createdId = id
	s.dto = dto
}

func (s *TestSuite) TearDownSuite() {
	advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, s.createdId)
	s.createdId = ""
}

// Проверяем, что статистика созданного объявления корректно выдается
// Тест-кейс №39
func (s *TestSuite) TestGetStatisticAdvertisementPositive() {
	s.Run("Успешный вывод статистики", func() {
		getStatistic := advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		asserts.AssertStatisticFieldsMatch(s.T(), s.dto.Statistics, getStatistic)
	})
}

// Проверяем 400 ошибку при попытке передать не UUID
// Тест-кейс №40
func (s *TestSuite) TestGetStatisticAdvertisementInvalidId() {
	s.Run("Ошибка 400 при передаче в path-параметр не UUID", func() {
		getStatistic := advertisements.GetStatisticAdvertisementByIdExpectError(s.T(), statusCode.BadRequest, "string")
		assert.NotEmpty(s.T(), getStatistic, "Пустое тело ответа")
		asserts.AssertResponseError(s.T(), getStatistic, "400",
			"передан некорректный идентификатор объявления")
	})
}

// Проверяем 404 ошибку при попытке передать несуществующий id
// Тест-кейс №41
func (s *TestSuite) TestGetStatisticAdvertisementNotFound() {
	s.Run("Ошибка 404 при передаче в path-параметр несуществующий id", func() {
		id := uuid.NewString()
		getStatistic := advertisements.GetStatisticAdvertisementByIdExpectError(s.T(), statusCode.NotFound, id)
		assert.NotEmpty(s.T(), getStatistic, "Пустое тело ответа")
		asserts.AssertResponseError(s.T(), getStatistic, "404",
			fmt.Sprintf("statistic %s not found", id))
	})
}

// Проверяем идемпотентность запросов
// Тест-кейс №42
func (s *TestSuite) TestGetStatisticIdempotency() {
	var firstResponse, secondResponse models.Statistics
	s.Run("Первый запрос на получение статистики", func() {
		firstResponse = advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		asserts.AssertStatisticFieldsMatch(s.T(), s.dto.Statistics, firstResponse)
	})
	s.Run("Второй запрос на получение статистики", func() {
		secondResponse = advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		asserts.AssertStatisticFieldsMatch(s.T(), s.dto.Statistics, secondResponse)
	})
	s.Run("Сравнение двух ответов", func() {
		assert.Equal(s.T(), firstResponse, secondResponse, "Ответы не совпадают")
	})
}

// Проверяем, что эндпоинт корректно отвечает на 10 параллельных запросов
// Тест-кейс №43
func (s *TestSuite) TestGetStatisticConcurrency() {
	s.Run("Создание 10-ти параллельных запросов", func() {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				advertisements.GetStatisticAdvertisementById(s.T(), statusCode.Ok, s.createdId)
			}()
		}
		wg.Wait()
	})
}
