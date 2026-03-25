package getAdvertisement

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
	dto := builder.NewCreateAdvertisementBuilder().Build()
	createdResponse := advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto)
	id := utils.ExtractUUIDFromResponse(createdResponse)
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

// Проверяем успешноем получение объявления
// Тест-кейс №34
func (s *TestSuite) TestGetAdvertisementPositive() {
	s.Run("Успешный запрос на получение объявления по существующему id", func() {
		advertisementResponse := advertisements.GetAdvertisementById(s.T(), statusCode.Ok, s.createdId)
		asserts.ResponseMatchesDto(s.T(), s.dto, advertisementResponse)
	})
}

// Проверяем 400 ошибку при попытке получения по строке(не UUID)
// Тест-кейс №35
func (s *TestSuite) TestGetAdvertisementInvalidIdFormat() {
	s.Run("Ошибка 400 при выдаче объявления по невалидному UUID", func() {
		str := utils.RandomString(5)
		errorResponse := advertisements.GetAdvertisementByIdExpectError(s.T(), statusCode.BadRequest, str)
		asserts.AssertResponseError(s.T(), errorResponse, "400", fmt.Sprintf("ID айтема не UUID: %s", str))
	})
}

// Проверяем ошибку 404 при попытке получения объявления по рандомному UUID
// Тест-кейс №36
func (s *TestSuite) TestGetAdvertisementByNonExistingId() {
	s.Run("Ошибка 404 при поиске объявления по несуществующему id", func() {
		id := uuid.NewString()
		errorResponse := advertisements.GetAdvertisementByIdExpectError(s.T(), statusCode.NotFound, id)
		asserts.AssertResponseError(s.T(), errorResponse, "404", fmt.Sprintf("item %s not found", id))
	})
}

// Проверяем идемпотентность запросов
// Тест-кейс №37
func (s *TestSuite) TestGetAdvertisementIdempotency() {
	id := s.createdId
	var secondResponse, firstResponse models.AdvertisementResponse
	s.Run("Первый запрос на получение объявления", func() {
		firstResponse = advertisements.GetAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.ResponseMatchesDto(s.T(), s.dto, firstResponse)
	})
	s.Run("Второй запрос на получение объявления", func() {
		secondResponse = advertisements.GetAdvertisementById(s.T(), statusCode.Ok, id)
		asserts.ResponseMatchesDto(s.T(), s.dto, secondResponse)
	})
	s.Run("Сравнение первого и второго ответов", func() {
		assert.Equal(s.T(), firstResponse, secondResponse, "Ответы не совпадают")
	})
}

// Проверяем, что эндпоинт корректно отвечает на 10 параллельных запросов
// Тест-кейс №38
func (s *TestSuite) TestGetAdvertisementConcurrency() {
	s.Run("Создание 10-ти параллельных запросов", func() {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				advertisements.GetAdvertisementById(s.T(), statusCode.Ok, s.createdId)
			}()
		}
		wg.Wait()
	})
}
