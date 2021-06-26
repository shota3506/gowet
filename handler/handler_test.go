package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shota3506/gowet/database"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	d := database.NewFakeClient()
	h := NewHandler(d)

	t.Run("github.com/tenntenn/greeting", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"http://example.com/github.com/tenntenn/greeting",
			nil,
		)
		recorder := httptest.NewRecorder()

		h.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("example.com", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"http://example.com/example.com",
			nil,
		)
		recorder := httptest.NewRecorder()

		h.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
