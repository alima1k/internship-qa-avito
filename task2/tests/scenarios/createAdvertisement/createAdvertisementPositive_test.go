package createAdvertisement

import (
	"encoding/json"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
	"task2/internal/models"
	"task2/internal/utils"
	"testing"

	consts "task2/internal/constants/statusCode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Структура для работы с позитивными тестами
type ValidTests []struct {
	name string
	dto  models.CreateAdvertisementRequest
}

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
		advertisements.DeleteAdvertisementById(s.T(), consts.Ok, id)
	}
	s.createdIds = nil
}

// Проверяем создание объявления с корректными полями
// Тест-кейсы №1-3
// Bugs: Bug-4
func (s *TestSuite) TestCreateAdvertisementPositive() {
	test := ValidTests{
		{
			name: "Успешное создание объявления с корректными полями",
			dto:  builder.NewCreateAdvertisementBuilder().Build(),
		},
		{
			name: "Успешное создание объявления с ценой и полями статистики = 0",
			dto: builder.NewCreateAdvertisementBuilder().WithPrice(0).WithContacts(0).
				WithViewCount(0).WithLikes(0).Build(),
		},
		{
			name: "Успешное создание объявления со спецсимволами в поле name",
			dto:  builder.NewCreateAdvertisementBuilder().WithName("@#$%^&*(){}[]+=-").Build(),
		},
	}
	s.createAdvertisement(test)
}

// Проверяем, что текстовое поле name защищено от различных атак
// Тест-кейсы № 32-33
func (s *TestSuite) TestSecurityScriptAndSQL() {
	tests := ValidTests{
		{
			name: "Проверка безопасности системы от XSS-атак",
			dto:  builder.NewCreateAdvertisementBuilder().WithName("<script>alert()</script>").Build(),
		},
		{
			name: "Проверка безопасности системы от SQL-инъекции",
			dto:  builder.NewCreateAdvertisementBuilder().WithName("'; DROP TABLE advertisements; --").Build(),
		},
	}
	s.createAdvertisement(tests)
}

// Проверяем, что при одинаковых полях создается 2 объявления с разным id
// Тест-кейс № 65
func (s *TestSuite) TestUniqueId() {
	var firstId, secondId string
	s.Run("Создание объявлений", func() {
		dto := builder.NewCreateAdvertisementBuilder().Build()
		firstResponse := advertisements.CreateAdvertisement(s.T(), 200, dto)
		secondResponse := advertisements.CreateAdvertisement(s.T(), 200, dto)

		firstId = utils.ExtractUUIDFromResponse(firstResponse)
		secondId = utils.ExtractUUIDFromResponse(secondResponse)
	})
	s.Run("Сравнение двух id", func() {
		assert.NotEqual(s.T(), firstId, secondId, "id совпадают при одинаковом теле")
	})
}

// Тест проверки ответа всегда падает, так как сервер отдает невалидное тело ответа на успешное создание
// Bugs: Bug-1
func (s *TestSuite) createAdvertisement(tests ValidTests) {
	for _, t := range tests {
		var createdAdvertisementResponse models.AdvertisementResponse
		var createdAdvertisementRaw string

		s.Run(t.name, func() {
			createdAdvertisementRaw = advertisements.CreateAdvertisement(s.T(), consts.Ok, t.dto)

			err := json.Unmarshal([]byte(createdAdvertisementRaw), &createdAdvertisementResponse)
			assert.NoError(s.T(), err, "Неверный формат ответа")
			s.createdIds = append(s.createdIds, createdAdvertisementResponse.Id)
		})

		s.Run("Проверяем ответ сервера", func() {
			assert.NotEmpty(s.T(), createdAdvertisementResponse.Id, "Поле id отсутствует. Получено: %s",
				createdAdvertisementRaw)
			asserts.ResponseMatchesDto(s.T(), t.dto, createdAdvertisementResponse)
		})
	}
}
