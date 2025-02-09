package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateUri(count int, city string) string {
	return fmt.Sprintf("/cafe?count=%d&city=%s", count, city)
}

func getResponse(req *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		generateUri(len(cafeList["moscow"]), "moscow"),
		nil)

	responseRecorder := getResponse(req)

	require.Equal(
		t,
		http.StatusOK,
		responseRecorder.Code,
		"Сервис должен возвращать код ответа 200")
	require.NotEmpty(
		t,
		responseRecorder.Body.String(),
		"Сервис должен возвращать не пустое тело ответа")
}

func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		generateUri(len(cafeList["Ankh-Morpork"]), "Ankh-Morpork"),
		nil)

	responseRecorder := getResponse(req)

	require.Equal(
		t,
		http.StatusBadRequest,
		responseRecorder.Code,
		"Сервис должен возвращать код ответа 400")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}

func TestMainHandlerWhenCafeCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		generateUri(len(cafeList["moscow"])+1, "moscow"),
		nil)

	responseRecorder := getResponse(req)

	require.Equal(
		t,
		http.StatusOK,
		responseRecorder.Code,
		"Сервис должен возвращать код ответа 200")
	// 	assert.Equal(
	//	t,
	//	strings.Join(cafeList["moscow"], ","),
	//	responseRecorder.Body.String(),
	//	"Сервис возвращает неверный список кафе - количество или порядок кафе в списке не верен")
	assert.ElementsMatch(
		t,
		cafeList["moscow"],
		strings.Split(responseRecorder.Body.String(), ","),
		"Сервис должен вернуть все доступные кафе")
}
