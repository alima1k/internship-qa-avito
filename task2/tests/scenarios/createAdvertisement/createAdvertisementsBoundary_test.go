package createAdvertisement

import "task2/internal/helpers/builder"

// Проверяем создание объявления с валидными граничными значениями
// Тест-кейсы № 25 - 27, 29, 30
func (s *TestSuite) TestCreateAdvertisementValidBoundary() {
	test := ValidTests{
		{
			name: "Успешное создание объявления с userId = 999998",
			dto:  builder.NewCreateAdvertisementBuilder().WithSellerId(999998).Build(),
		},
		{
			name: "Успешное создание объявления с userId = 999999",
			dto:  builder.NewCreateAdvertisementBuilder().WithSellerId(999999).Build(),
		},
		{
			name: "Успешное создание объявления с userId = 111111",
			dto:  builder.NewCreateAdvertisementBuilder().WithSellerId(111111).Build(),
		},
		{
			name: "Успешное создание объявления с userId = 111112",
			dto:  builder.NewCreateAdvertisementBuilder().WithSellerId(111112).Build(),
		},
		{
			name: "Успешное создание объявления с одним символом в поле name",
			dto:  builder.NewCreateAdvertisementBuilder().WithName("a").Build(),
		},
	}
	s.createAdvertisement(test)
}

// Проверяем, что объявление не создается при невалидных граничных значениях
// Тест-кейсы №24, 28, 31
// Bugs: Bug-6
func (s *TestSuite) TestCreateAdvertisementInvalidBoundary() {
	tests := InvalidTests{
		{
			name:               "Ошибка создания объявления с полем пустым name",
			dto:                builder.NewInvalidAdvertisementBuilder().NameAsEmpty().Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"name": "поле name обязательно"},
		},
		{
			name:               "Ошибка создания объявления с sellerId = 111110",
			dto:                builder.NewInvalidAdvertisementBuilder().WithSellerId(111110).Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"userId": "значение поля от 111111 до 999999"},
		},
		{
			name:               "Ошибка создания объявления с sellerId = 1000000",
			dto:                builder.NewInvalidAdvertisementBuilder().WithSellerId(1000000).Build(),
			expectedStatus:     badRequestStatus,
			expectedMessage:    validationMessage,
			expectedFieldError: map[string]string{"userId": "значение поля от 111111 до 999999"},
		},
	}
	s.createAdvertisementExpectError(tests)
}
