package http

import (
	"fmt"
	"net/http"
	"task2/internal/constants/path"
	api_runner "task2/internal/helpers/api-runner"
	"testing"
)

// HttpCreateAdvertisement interface{} здесь, чтобы отправлять как валидную структуру запроса, так и невалидную
func HttpCreateAdvertisement(t *testing.T, request interface{}) *http.Response {
	return api_runner.GetRunner().Create().Post(path.CreateAdvPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).End().Response
}

func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return api_runner.GetRunner().Create().Get(fmt.Sprintf(path.GetAdvByIdPath, id)).
		Expect(t).End().Response
}

func HttpDeleteAdvertisement(t *testing.T, id string) *http.Response {
	return api_runner.GetRunner().Create().Delete(fmt.Sprintf(path.DeleteAdvByIdPath, id)).
		Expect(t).End().Response
}

// interface{} здесь, чтобы обрабатывать как строку, так и integer
func HttpGetSellerAdvertisements(t *testing.T, sellerId interface{}) *http.Response {
	return api_runner.GetRunner().Create().Get(fmt.Sprintf(path.GetSellerAdvPath, sellerId)).
		Expect(t).End().Response
}

func HttpGetStatisticByAdsId(t *testing.T, id string) *http.Response {
	return api_runner.GetRunner().Create().Get(fmt.Sprintf(path.GetStatisticByAdvIdPath, id)).
		Expect(t).End().Response
}
