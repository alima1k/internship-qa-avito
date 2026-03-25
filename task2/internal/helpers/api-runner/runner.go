package api_runner

import (
	"net/http"

	"github.com/steinfletcher/apitest"
)

// Класс предоставляет обёртку над библиотекой apitest для упрощённого создания HTTP-запросов с UserAgent в тестах
const (
	defaultUserAgent = "Chrome/143.0.0.0"
)

type ApiTest struct {
	userAgent string
}

func New() *ApiTest {
	return &ApiTest{
		userAgent: defaultUserAgent,
	}
}

func GetRunner() *ApiTest {
	return New()
}

func (at *ApiTest) Create() *apitest.APITest {
	return apitest.New().
		EnableNetworking().
		Intercept(func(request *http.Request) {
			request.Header.Set("User-Agent", at.userAgent)
		})
}
