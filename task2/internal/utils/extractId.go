package utils

import (
	"encoding/json"
	"strings"
)

type StatusResponse struct {
	Status string `json:"status"`
}

// ExtractUUIDFromResponse получает id созданного объявления из невалидного ответа API на создание
func ExtractUUIDFromResponse(raw string) string {
	var resp StatusResponse
	err := json.Unmarshal([]byte(raw), &resp)
	if err != nil {
		return ""
	}
	parts := strings.Split(resp.Status, " - ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
