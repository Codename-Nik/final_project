package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"final_project/pkg/handlers"
	"final_project/pkg/utils"
)

func TestGetCurrentTime(t *testing.T) {
	currentTime := utils.GetCurrentTime()

	if currentTime == "" {
		t.Error("GetCurrentTime() вернула пустую строку")
	}

	if len(currentTime) != 19 {
		t.Errorf("Неверная длина формата времени: ожидалось 19, получено %d", len(currentTime))
	}

	if !strings.Contains(currentTime, "-") || !strings.Contains(currentTime, ":") {
		t.Error("Неверный формат времени: отсутствуют разделители")
	}
}

func TestGetCurrentTimeNotEmpty(t *testing.T) {
	time1 := utils.GetCurrentTime()
	time.Sleep(1 * time.Second)
	time2 := utils.GetCurrentTime()

	if time1 == time2 {
		t.Error("Время должно различаться между вызовами")
	}
}

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handlers.HomeHandler(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Result().StatusCode)
	}
}

func TestHomeHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handlers.HomeHandler(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Ожидался статус 404, получен %d", w.Result().StatusCode)
	}
}

func TestTimeHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/time", nil)
	w := httptest.NewRecorder()

	handlers.TimeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
}

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "С параметром name",
			url:            "/index?name=Тест",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Без параметра",
			url:            "/index",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "С пустым параметром",
			url:            "/index?name=",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "С именем admin",
			url:            "/index?name=admin",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handlers.IndexHandler(w, req)

			if w.Result().StatusCode != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d",
					tt.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func TestAPIHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()

	handlers.APIDataHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
	}
}

func TestAPIResponseStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()

	handlers.APIDataHandler(w, req)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)

	if err != nil {
		t.Errorf("Ошибка парсинга JSON: %v", err)
	}

	if _, exists := response["data"]; !exists {
		t.Error("Поле 'data' отсутствует в JSON ответе")
	}

	if _, exists := response["meta"]; !exists {
		t.Error("Поле 'meta' отсутствует в JSON ответе")
	}

	if data, ok := response["data"].(map[string]interface{}); ok {
		requiredFields := []string{"message", "timestamp", "status", "server", "version"}
		for _, field := range requiredFields {
			if _, exists := data[field]; !exists {
				t.Errorf("Поле '%s' отсутствует в data ответе", field)
			}
		}
	}
}

func TestContentType(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()

	handlers.APIDataHandler(w, req)

	if contentType := w.Header().Get("Content-Type"); contentType == "" {
		t.Error("Заголовок Content-Type отсутствует")
	}

	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		t.Errorf("Неверный Content-Type: %s", contentType)
	}
}

func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()

	handlers.APIDataHandler(w, req)

	if cors := w.Header().Get("Access-Control-Allow-Origin"); cors != "*" {
		t.Errorf("Ожидался заголовок CORS *, получен %s", cors)
	}
}

func BenchmarkGetCurrentTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GetCurrentTime()
	}
}

func BenchmarkHomeHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handlers.HomeHandler(w, req)
	}
}

func BenchmarkAPIDataHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		w := httptest.NewRecorder()
		handlers.APIDataHandler(w, req)
	}
}
