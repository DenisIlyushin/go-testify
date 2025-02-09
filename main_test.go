package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type test struct {
	name           string
	count          int
	city           string
	expectedStatus int
	expectedBody   string
}

// Генерирует строку с параметрами для запроса
func generateUri(count int, city string) string {
	return fmt.Sprintf("/cafe?count=%d&city=%s", count, city)
}

// Получает ответ тестового сервера
func getResponse(req *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

// Реализует тесты при успешном ответе
func TestMainHandler(t *testing.T) {
	tests := []test{
		{
			name:           "Valid city and count",
			count:          len(cafeList["moscow"]),
			city:           "moscow",
			expectedStatus: http.StatusOK,
			expectedBody:   strings.Join(cafeList["moscow"], ","),
		},
		{
			name:           "More cafes requested than available",
			count:          len(cafeList["moscow"]) + 1,
			city:           "moscow",
			expectedStatus: http.StatusOK,
			expectedBody:   strings.Join(cafeList["moscow"], ","),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				generateUri(tt.count, tt.city),
				nil,
			)

			responseRecorder := getResponse(req)

			require.Equal(t, tt.expectedStatus, responseRecorder.Code)
			assert.ElementsMatch(t, cafeList["moscow"], strings.Split(responseRecorder.Body.String(), ","))
		})
	}
}

// Реализует тесты при ответе с ошибками
func TestMainHandlerForBadRequest(t *testing.T) {
	tests := []test{
		{
			name:           "Unsupported city",
			count:          1,
			city:           "ankh-morpork",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong city value",
		},
		{
			name:           "Invalid count value",
			count:          -1,
			city:           "moscow",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong count value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				generateUri(tt.count, tt.city),
				nil,
			)

			responseRecorder := getResponse(req)

			require.Equal(t, tt.expectedStatus, responseRecorder.Code)
			assert.Equal(t, tt.expectedBody, responseRecorder.Body.String())
		})
	}
}

func TestMainHandlerWhenCountIsMissing(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/cafe?city=moscow", // Missing city parameter
		nil,
	)

	responseRecorder := getResponse(req)

	require.Equal(
		t,
		http.StatusBadRequest, // Expect 400 if city is required
		responseRecorder.Code,
		"Сервис должен возвращать код ответа 400 при отсутствии параметра city",
	)

	assert.Equal(
		t,
		"count missing", // Adjust expected error message based on your handler logic
		responseRecorder.Body.String(),
	)
}
