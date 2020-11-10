// +build unit

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	type want struct {
		status int
	}

	test := func(
		timeout time.Duration,
		next http.Handler,
		callback func(http.ResponseWriter),
		want want,
	) func(*testing.T) {
		return func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			WithTimeout(timeout, callback)(next).ServeHTTP(w, r)

			if w.Code != want.status {
				t.Errorf("wrong status code: got %v want %v", w.Code, want.status)
			}
		}
	}

	tests := []struct {
		name     string
		timeout  time.Duration
		next     http.Handler
		callback func(http.ResponseWriter)
		want     want
	}{
		{
			name:    "no deadline",
			timeout: time.Second,
			next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name:    "deadline",
			timeout: time.Second,
			next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				time.Sleep(2 * time.Second)
			}),
			callback: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		name := tt.name
		t.Run(name, test(tt.timeout, tt.next, tt.callback, tt.want))
	}
}
