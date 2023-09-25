package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-messages/internal/rest/mocks"
	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_deleteMessageData(t *testing.T) {
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
		wantStatus  int
	}{
		{
			name: "Data deleted (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "jane", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 55,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					s.On("DeleteMessageData", mock.Anything, int64(42), int64(55)).Return(int64(60), nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"ETag": "60",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid id (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "bad-id")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Invalid timestamp (412)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "something")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusPreconditionFailed,
		},
		{
			name: "Message not found (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "jane", int64(42)).Return(nil, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Timestamp is not match (412)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "jane", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 56,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					s.On("DeleteMessageData", mock.Anything, int64(42), int64(55)).Return(int64(0),
						&ErrTimestampIsNotMatchTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusPreconditionFailed,
		},
		{
			name: "Message deleted (410)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "jane", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 55,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					s.On("DeleteMessageData", mock.Anything, int64(42), int64(55)).Return(int64(0),
						&ErrMessageDeletedTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusGone,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/{id}", nil)
					r.Header.Set("If-Match", "55")
					r.Header.Set("request-id", "test-request-id")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "jane", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 55,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					s.On("DeleteMessageData", mock.Anything, int64(42), int64(55)).Return(int64(0),
						errors.New("test error"))
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
			s.deleteMessageData(tt.args.w, tt.args.r)

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
		})
	}
}
