package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-messages/internal/rest/mocks"
	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_syncMessages(t *testing.T) {
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
		wantBody    *models.MessageUpdatesV1
		wantStatus  int
	}{
		{
			name: "Updates received - default parameters (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "application/vnd.messageUpdates.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("MessageUpdatesV1", mock.Anything, mock.Anything, int64(0), 50).Return(
						&models.MessageUpdatesV1{
							Total: 3,
							Messages: []*models.MessageUpdateInfoV1{
								{Id: 100, Timestamp: 200},
								{Id: 110, Timestamp: 215},
								{Id: 100, Timestamp: 240},
							},
						},
						nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/vnd.messageUpdates.v1+json",
			},
			wantBody: &models.MessageUpdatesV1{
				Total: 3,
				Messages: []*models.MessageUpdateInfoV1{
					{Id: 100, Timestamp: 200},
					{Id: 110, Timestamp: 215},
					{Id: 100, Timestamp: 240},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Updates received - specified parameters (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?after=99&limit=20", nil)
					r.Header.Set("Accept", "application/vnd.messageUpdates.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("MessageUpdatesV1", mock.Anything, mock.Anything, int64(99), 20).Return(
						&models.MessageUpdatesV1{
							Total: 3,
							Messages: []*models.MessageUpdateInfoV1{
								{Id: 100, Timestamp: 200},
								{Id: 110, Timestamp: 215},
								{Id: 100, Timestamp: 240},
							},
						},
						nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/vnd.messageUpdates.v1+json",
			},
			wantBody: &models.MessageUpdatesV1{
				Total: 3,
				Messages: []*models.MessageUpdateInfoV1{
					{Id: 100, Timestamp: 200},
					{Id: 110, Timestamp: 215},
					{Id: 100, Timestamp: 240},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Incorrect parameters (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?after=99&limit=200", nil)
					r.Header.Set("Accept", "application/vnd.messageUpdates.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Requested media type is not supported (406)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?after=99&limit=200", nil)
					r.Header.Set("Accept", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotAcceptable,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "application/vnd.messageUpdates.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("MessageUpdatesV1", mock.Anything, mock.Anything, int64(0), 50).Return(nil, errors.New("test error"))
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
			s.syncMessages(tt.args.w, tt.args.r)

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

			var body *models.MessageUpdatesV1
			decoded := models.MessageUpdatesV1{}
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
