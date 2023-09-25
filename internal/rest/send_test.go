package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-messages/internal/rest/mocks"
	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_sendNewMessage(t *testing.T) {
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
			name: "Message sent (201)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonalMessageV1{
						To:   "john",
						Text: "Hello!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPersonalMessage.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateNewPersonalMessageV1", mock.Anything, mock.Anything, mock.Anything).Return(int64(123), int64(456), nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Location": "/123",
				"ETag":     "456",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Incorrect message data (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonalMessageV1{
						Text: "Hello!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPersonalMessage.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Unsupported message data (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonalMessageV1{
						To:   "john",
						Text: "Hello!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonalMessageV1{
						To:   "john",
						Text: "Hello!",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPersonalMessage.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateNewPersonalMessageV1", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), int64(0),
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
			s.sendNewMessage(tt.args.w, tt.args.r)

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
