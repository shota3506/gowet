package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shota3506/gowet/database"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Run("github.com/tenntenn/greeting", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mock := database.NewMockDB(mockCtrl)
		gomock.InOrder(
			mock.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string) {}).
				Return(nil, database.NewNotFoundError(errors.New("some error"))),
			mock.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string) {}).
				Return(nil, database.NewNotFoundError(errors.New("some error"))),
			mock.EXPECT().
				Set(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _, _ string) {}).
				Return(nil),
		)

		h := NewHandler(mock)

		req := httptest.NewRequest(
			http.MethodGet,
			"http://example.com/github.com/tenntenn/greeting",
			nil,
		)
		recorder := httptest.NewRecorder()

		h.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("github.com/tenntenn/greeting cached", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mock := database.NewMockDB(mockCtrl)
		mock.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ string) {}).
			Return([]byte("{}"), nil)

		h := NewHandler(mock)

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
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mock := database.NewMockDB(mockCtrl)
		gomock.InOrder(
			mock.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string) {}).
				Return(nil, database.NewNotFoundError(errors.New("some error"))),
			mock.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string) {}).
				Return(nil, database.NewNotFoundError(errors.New("some error"))),
		)

		h := NewHandler(mock)

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
