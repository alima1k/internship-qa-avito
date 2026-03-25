package createAdvertisement

import (
	consts "task2/internal/constants/statusCode"
	"task2/internal/helpers/asserts"
	"task2/internal/helpers/builder"
	"task2/internal/managers/advertisements"
)

const (
	badRequestStatus  = "400"
	validationMessage = "Ошибка валидации"
)

// Структура для работы с негативными тестами. Содержит ожидаемые поля ошибок
// expectedStatus - значения, ожидаемое в поле status
// expectedMessage - значение, ожидаемое в result.message
// expectedFieldError - пары значений названия поля и ошибки в этом поле
type InvalidTests []struct {
	name               string
	dto                map[string]interface{}
	expectedStatus     string
	expectedMessage    string
	expectedFieldError map[string]string
}

// Проверяем, что валидация на null работает на всех полях сначала по отдельности, потом все вместе
// Тест-кейсы №4-9, 20
// Bugs: Bug-2
func (s *TestSuite) TestCreateAdvertisementNullValidation() {
	tests := InvalidTests{
		{
			name:               "Ошибка создания объявления c sellerId = null",
			dto:                builder.NewInvalidAdvertisementBuilder().SellerIdAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"sellerID": "sellerID не может быть null"},
		},
		{
			name:               "Ошибка создания объявления c name = null",
			dto:                builder.NewInvalidAdvertisementBuilder().NameAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"name": "name не может быть null"},
		},
		{
			name:               "Ошибка создания объявления c price = null",
			dto:                builder.NewInvalidAdvertisementBuilder().PriceAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"price": "price не может быть null"},
		},
		{
			name:               "Ошибка создания объявления c likes = null",
			dto:                builder.NewInvalidAdvertisementBuilder().LikesAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"likes": "likes не может быть null"},
		},
		{
			name:               "Ошибка создания объявления c viewCount = null",
			dto:                builder.NewInvalidAdvertisementBuilder().ViewCountAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"viewCount": "viewCount не может быть null"},
		},
		{
			name:               "Ошибка создания объявления c contacts = null",
			dto:                builder.NewInvalidAdvertisementBuilder().ContactsAsNull().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"contacts": "contacts не может быть null"},
		},
		{
			name: "Ошибка создания объявления cо всеми полями = null",
			dto: builder.NewInvalidAdvertisementBuilder().ContactsAsNull().SellerIdAsNull().PriceAsNull().
				NameAsNull().LikesAsNull().ViewCountAsNull().Build(),
			expectedStatus:  badRequestStatus,
			expectedMessage: validationMessage,
			expectedFieldError: map[string]string{
				"contacts":  "contacts не может быть null",
				"sellerID":  "sellerID не может быть null",
				"price":     "price не может быть null",
				"name":      "name не может быть null",
				"likes":     "likes не может быть null",
				"viewCount": "viewCount не может быть null",
			},
		},
	}

	s.createAdvertisementExpectError(tests)
}

// Проверяем, что валидация на отрицательные числа полей price, likes, viewCount, contacts работает
// №16-19, 21
// Bugs: Bug-3
func (s *TestSuite) TestCreateAdvertisementInvalidNegativeValues() {
	tests := InvalidTests{
		{
			name:               "Ошибка создания объявления с отрицательным полем price",
			dto:                builder.NewInvalidAdvertisementBuilder().PriceAsNegative().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"price": "price не может быть отрицательным"},
		},
		{
			name:               "Ошибка создания объявления с отрицательным полем likes",
			dto:                builder.NewInvalidAdvertisementBuilder().LikesAsNegative().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"likes": "likes не может быть отрицательным"},
		},
		{
			name:               "Ошибка создания объявления с отрицательным полем viewCount",
			dto:                builder.NewInvalidAdvertisementBuilder().ViewCountAsNegative().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"viewCount": "viewCount не может быть отрицательным"},
		},
		{
			name:               "Ошибка создания объявления с отрицательным полем contacts",
			dto:                builder.NewInvalidAdvertisementBuilder().ContactsAsNegative().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"contacts": "contacts не может быть отрицательным"},
		},
		{
			name: "Ошибка создания объявления с отрицательными полями price, likes, viewCount, contacts",
			dto: builder.NewInvalidAdvertisementBuilder().LikesAsNegative().ViewCountAsNegative().
				PriceAsNegative().ContactsAsNegative().Build(),
			expectedStatus:  badRequestStatus,
			expectedMessage: validationMessage,
			expectedFieldError: map[string]string{
				"contacts":  "contacts не может быть отрицательным",
				"price":     "price не может быть отрицательным",
				"likes":     "likes не может быть отрицательным",
				"viewCount": "viewCount не может быть отрицательным",
			},
		},
	}
	s.createAdvertisementExpectError(tests)
}

// Проверяем 400 ошибку при попытке вместо нужного типа отправить другой
// Тест-кейсы №10-15,22
// Bugs: Bug-5
func (s *TestSuite) TestInvalidTypeFields() {
	tests := InvalidTests{
		{
			name:               "Ошибка создания объявления со строкой в поле sellerID",
			dto:                builder.NewInvalidAdvertisementBuilder().SellerIdAsString().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"sellerID": "Неверный формат поля sellerID"},
		},
		{
			name:               "Ошибка создания объявления со строкой в поле price",
			dto:                builder.NewInvalidAdvertisementBuilder().PriceAsString().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"price": "Неверный формат поля price"},
		}, {
			name:               "Ошибка создания объявления со строкой в поле likes",
			dto:                builder.NewInvalidAdvertisementBuilder().LikesAsString().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"likes": "Неверный формат поля likes"},
		},
		{
			name:               "Ошибка создания объявления со строкой в поле viewCount",
			dto:                builder.NewInvalidAdvertisementBuilder().ViewCountAsString().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"viewCount": "Неверный формат поля viewCount"},
		},
		{
			name:               "Ошибка создания объявления со строкой в поле contacts",
			dto:                builder.NewInvalidAdvertisementBuilder().ContactsAsString().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"contacts": "Неверный формат поля contacts"},
		},
		{
			name:               "Ошибка создания объявления с числом в поле name",
			dto:                builder.NewInvalidAdvertisementBuilder().NameAsNumber().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"name": "Неверный формат поля name"},
		},
		{
			name: "Ошибка создания объявления со всеми неверными типами",
			dto: builder.NewInvalidAdvertisementBuilder().NameAsNumber().PriceAsString().SellerIdAsString().
				LikesAsString().ViewCountAsString().ContactsAsString().Build(),
			expectedStatus:  badRequestStatus,
			expectedMessage: validationMessage,
			expectedFieldError: map[string]string{
				"name":      "Неверный формат поля name",
				"contacts":  "Неверный формат поля contacts",
				"sellerID":  "Неверный формат поля sellerID",
				"price":     "Неверный формат поля price",
				"likes":     "Неверный формат поля likes",
				"viewCount": "Неверный формат поля viewCount",
			},
		},
	}
	s.createAdvertisementExpectError(tests)
}

// Проверяем ошибку при отсутствии обязательных полей
// Тест-кейс №23
// Bugs: Bug-2
func (s *TestSuite) TestCreateAdvertisementWithoutRequiredFields() {
	tests := InvalidTests{
		{
			name:            "Ошибка создания объявления без полей price и sellerID",
			dto:             builder.NewInvalidAdvertisementBuilder().WithoutPrice().WithoutSellerId().Build(),
			expectedStatus:  badRequestStatus,
			expectedMessage: validationMessage,
			expectedFieldError: map[string]string{
				"price":    "поле price обязательно",
				"sellerID": "поле sellerId обязательно"},
		},
	}

	s.createAdvertisementExpectError(tests)
}

func (s *TestSuite) createAdvertisementExpectError(tests InvalidTests) {
	for _, t := range tests {
		s.Run(t.name, func() {
			raw := advertisements.CreateAdvertisementFromMap(s.T(), consts.BadRequest, t.dto)

			asserts.AssertResponseFieldError(s.T(), raw, t.expectedStatus, t.expectedMessage, t.expectedFieldError)
		})
	}
}
