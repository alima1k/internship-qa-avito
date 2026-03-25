package getSellerAdvertisements

import (
	"sync"
	"task2/internal/constants/statusCode"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
	"task2/internal/models"
	"task2/internal/utils"
	"testing"

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

// Перед запуском всех тестов создаем объявление продавцу и сохраняем его id и dto для дальнейших проверок
func (s *TestSuite) SetupSuite() {
	dto := builder.NewCreateAdvertisementBuilder().WithSellerId(utils.GenerateSellerId()).Build()
	id := utils.ExtractUUIDFromResponse(advertisements.CreateAdvertisement(s.T(), statusCode.Ok, dto))
	if id == "" {
		s.T().Fatal("Не удалось получить UUID")
	}
	s.createdId = id
	s.dto = dto
}

// После всех тестов удаляем ранее созданное объявление
func (s *TestSuite) TearDownSuite() {
	advertisements.DeleteAdvertisementById(s.T(), statusCode.Ok, s.createdId)
}

// Проверяем, что у продавца можно получить объявления
// Тест-кейс №44
func (s *TestSuite) TestGetSellerAdvertisementPositive() {
	s.Run("Успешное получение всех объявлений продавца", func() {
		advertisementResponse := advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, s.dto.SellerId)
		assert.GreaterOrEqual(s.T(), len(advertisementResponse), 1,
			"У продавца не отображаются объявления")
	})
}

// Проверяем, что при строке в pathVariable выдает 400
// Тест-кейс №46
func (s *TestSuite) TestGetSellerAdvertisementInvalidSellerId() {
	s.Run("Ошибка 400 при невалидном id продавца", func() {
		stringSellerId := "string"
		advertisementResponse := advertisements.GetSellerAdvertisementsExpectError(s.T(), statusCode.BadRequest, stringSellerId)
		asserts.AssertResponseError(s.T(), advertisementResponse, "400",
			"передан некорректный идентификатор продавца")
	})
}

// Проверяем граничные значения
// Тест-кейсы №47-52
// Bugs: Bug-9
func (s *TestSuite) TestGetSellerAdvertisementsSellerIdValidation() {
	validSellerIds := []int{111111, 999999, 999998, 111112}
	invalidSellerIds := []int{111110, 1000000}

	s.Run("Проверка валидных граничных значений", func() {
		for _, id := range validSellerIds {
			getSellerAds := advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, id)
			assert.IsType(s.T(), []models.AdvertisementResponse{}, getSellerAds, "Ответ не является массивом")
		}
	})

	s.Run("Проверка невалидных граничных значений", func() {
		for _, id := range invalidSellerIds {
			advertisements.GetSellerAdvertisementsExpectError(s.T(), statusCode.BadRequest, id)
		}
	})
}

// Проверяем идемпотентность запросов
// Тест-кейсы №53
func (s *TestSuite) TestGetSellerAdvertisementsIdempotency() {
	id := s.dto.SellerId
	var firstResponse, secondResponse []models.AdvertisementResponse
	s.Run("Первый запрос на получение объявлений продавца", func() {
		firstResponse = advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, id)
		firstAds := findAdvertisementById(firstResponse, s.createdId)
		assert.NotNil(s.T(), firstAds, "В первом ответе отсутствует созданное объявление")
	})
	s.Run("Второй запрос на получение объявлений продавца", func() {
		secondResponse = advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, id)
		secondAds := findAdvertisementById(secondResponse, s.createdId)
		assert.NotNil(s.T(), secondAds, "Во втором ответе отсутствует созданное объявление")
	})
	s.Run("Сравнение первого и второго ответов", func() {
		assert.Equal(s.T(), firstResponse, secondResponse, "Ответы не совпадают")
	})
}

// Проверяем, что эндпоинт корректно отвечает на 10 параллельных запросов
// Тест-кейс №54
func (s *TestSuite) TestGetSellerAdvertisementsConcurrency() {
	s.Run("Создание 10-ти параллельных запросов", func() {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				response := advertisements.GetSellerAdvertisementsBySellerId(s.T(), statusCode.Ok, s.dto.SellerId)
				assert.GreaterOrEqual(s.T(), len(response), 1, "У продавца не отобразилось объявление")
			}()
		}
		wg.Wait()
	})
}

// Вспомогательный метод для поиска среди тела ответа созданное ранее объявление
func findAdvertisementById(ads []models.AdvertisementResponse, id string) *models.AdvertisementResponse {
	for i := range ads {
		if ads[i].Id == id {
			return &ads[i]
		}
	}
	return nil
}
