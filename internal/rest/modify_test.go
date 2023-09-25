package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-messages/internal/rest/mocks"
	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_modifyMessage(t *testing.T) {
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
			name: "Modified - text (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					s.On("EditMessageText", mock.Anything, int64(42), int64(55), "Hi!").Return(int64(60), nil)
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
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "jane"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "something")
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
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
			name: "Bad body - text (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.MessageReadMarkV1{
						Read: true,
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Only sender can edit message text (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "john"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "john", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 55,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Empty text without attachments (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Timestamp is not match (412)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "Hi!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					s.On("EditMessageText", mock.Anything, int64(42), int64(55), "Hi!").Return(int64(0),
						&ErrTimestampIsNotMatchTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusPreconditionFailed,
		},
		{
			name: "Not modified (304)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "Hello",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					s.On("EditMessageText", mock.Anything, int64(42), int64(55), "Hello").Return(int64(0),
						&ErrMessageNotModifiedTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotModified,
		},
		{
			name: "Message deleted (410)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedMessageTextV1{
						Text: "Hello",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.editedMessageText.v1+json")
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
					s.On("EditMessageText", mock.Anything, int64(42), int64(55), "Hello").Return(int64(0),
						&ErrMessageDeletedTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusGone,
		},
		{
			name: "Modified - read (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.MessageReadMarkV1{
						Read: true,
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.messageReadMark.v1+json")
					r.Header.Set("If-Match", "55")
					r = r.WithContext(context.WithValue(r.Context(), authenticatedUserId{}, "john"))
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "42")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("PersonalMessageV1", mock.Anything, "john", int64(42)).Return(
						&models.PersonalMessageV1{
							Id:        42,
							Timestamp: 55,
							From:      "jane",
							To:        "john",
							Text:      "Hello",
						},
						nil)
					s.On("SetMessageReadState", mock.Anything, int64(42), int64(55), true).Return(int64(60), nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"ETag": "60",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Only receiver can mark message as read (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.MessageReadMarkV1{
						Read: true,
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/{id}", &buf)
					r.Header.Set("Content-Type", "application/vnd.messageReadMark.v1+json")
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
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.modifyMessage(tt.args.w, tt.args.r)

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

type ErrMessageDeletedTest struct{}
type ErrTimestampIsNotMatchTest struct{}
type ErrMessageNotModifiedTest struct{}

func (e *ErrMessageDeletedTest) Error() string {
	return "message deleted"
}

func (e *ErrMessageDeletedTest) ImplementsMessageDeletedError() {
}

func (e *ErrTimestampIsNotMatchTest) Error() string {
	return "message timestamp is not match"
}

func (e *ErrTimestampIsNotMatchTest) ImplementsTimestampIsNotMatchError() {
}

func (e *ErrMessageNotModifiedTest) Error() string {
	return "message has not been modified"
}

func (e *ErrMessageNotModifiedTest) ImplementsMessageNotModifiedError() {
}
