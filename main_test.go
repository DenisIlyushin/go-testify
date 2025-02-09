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
			name:           "Unsupported city",
			count:          1,
			city:           "ankh-morpork",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong city value",
		},
		{
			name:           "More cafes requested than available",
			count:          len(cafeList["moscow"]) + 1,
			city:           "moscow",
			expectedStatus: http.StatusOK,
			expectedBody:   strings.Join(cafeList["moscow"], ","),
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

			response := getResponse(req)

			require.Equal(t, tt.expectedStatus, response.Code)

			switch tt.expectedStatus {
			case http.StatusOK:
				assert.ElementsMatch(t, cafeList["moscow"], strings.Split(response.Body.String(), ","))
			case http.StatusBadRequest:
				assert.Equal(t, tt.expectedBody, response.Body.String())
			default:
				assert.Equal(t, "", response.Body.String())
			}
		})
	}
}

type edgeTest struct {
	name           string
	uri            string
	expectedStatus int
	expectedBody   string
}

// Проверяет ответ при пограничных значениях,
func TestMainHandlerWithEdgeConditions(t *testing.T) {
	edgeTests := []edgeTest{
		{
			name:           "City count parameter is missing",
			uri:            "/cafe?city=moscow",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "count missing",
		},
		{
			name:           "City parameter is missing",
			uri:            "/cafe?count=3",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong city value",
		},
		{
			name:           "Negative count value",
			uri:            "/cafe?count=-1&city=moscow",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong count value",
		},
		{
			name:           "Non-numeric count value",
			uri:            "/cafe?count=ahaha&city=moscow",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong count value",
		},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.uri, nil)

			response := getResponse(req)

			require.Equal(t, tt.expectedStatus, response.Code)
			assert.Equal(t, tt.expectedBody, response.Body.String())
		})
	}

}
