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

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK)
	require.NotEmpty(t, responseRecorder.Body)
}

func TestMainHandleWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=tomsk", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)

	assert.Equal(t, responseRecorder.Body.String(), "wrong city value")
}

func TestMainHandleWrongCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=test&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)

	assert.Equal(t, responseRecorder.Body.String(), "wrong count value")
}

func TestMainHandlerWhenMissingCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)

	assert.Equal(t, responseRecorder.Body.String(), "count missing")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	urlPath := fmt.Sprintf("/cafe?city=moscow&count=%d", totalCount+1)

	req := httptest.NewRequest("GET", urlPath, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK)
	require.NotEmpty(t, responseRecorder.Body)

	listCafe := responseRecorder.Body.String()
	expectedSlice := strings.Split(listCafe, ",")
	expectedCount := len(expectedSlice)

	assert.Equal(t, totalCount, expectedCount)
}
