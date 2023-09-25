package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-messages/internal/rest/mocks"
	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_getMessageData(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		testService testService
		args        args
		wantHeaders map[string]string
		wantBody    *models.PersonalMessageV1
		wantStatus  int
	}{
		{
			name: "OK (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, mock.Anything, int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 67,
							Deleted:   true,
						},
						nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/vnd.personalMessage.v1+json",
			},
			wantBody: &models.PersonalMessageV1{
				Id:        42,
				Timestamp: 67,
				Deleted:   true,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Not found - bad id (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "bad-id")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Not found (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, mock.Anything, int64(42)).Return(nil, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/{id}", nil)
					r.Header.Set("request-id", "test-request-id")
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, mock.Anything, int64(42)).Return(nil, errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.getMessageData(tt.args.w, tt.args.r)

			for k, v := range tt.wantHeaders {
				require.Equal(t, v, func() string {
					h := tt.args.w.Result().Header
					if h == nil {
						return ""
					}
					v := h[k]
					if len(v) == 0 {
						return ""
					}
					return v[0]
				}())
			}

			require.Equal(t, tt.wantStatus, tt.args.w.Code)

			if tt.wantBody == nil {
				return
			}

			var body *models.PersonalMessageV1
			decoded := models.PersonalMessageV1{}
			err := json.NewDecoder(tt.args.w.Body).Decode(&decoded)

			if err != nil && err != io.EOF {
				t.Fatal(err)
			}

			if err == nil {
				body = &decoded
			}

			require.Equal(t, body, tt.wantBody)
		})
	}
}
