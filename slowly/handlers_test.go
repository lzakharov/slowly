// +build unit

package slowly

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSlow(t *testing.T) {
	type want struct {
		status      int
		contentType string
		data        string
	}

	test := func(ctx context.Context, body io.Reader, want want) func(*testing.T) {
		return func(t *testing.T) {
			r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/", body)
			w := httptest.NewRecorder()
			Slow(w, r)

			if w.Code != want.status {
				t.Errorf("wrong status code: got %v want %v", w.Code, want.status)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != want.contentType {
				t.Errorf("wrong content type: got %v want %v", contentType, want.contentType)
			}

			data := w.Body.String()
			if !strings.EqualFold(data, want.data) {
				t.Errorf("wrong response data: got '%v' want '%v'", data, want.data)
			}
		}
	}

	tests := []struct {
		name string
		ctx  context.Context
		body io.Reader
		want want
	}{
		{
			name: "happy path",
			ctx:  context.Background(),
			body: bytes.NewBufferString(`{"timeout": 3000}`),
			want: want{
				status:      http.StatusOK,
				contentType: "application/json",
				data:        `{"status":"ok"}`,
			},
		},
		{
			name: "empty request",
			ctx:  context.Background(),
			body: bytes.NewBufferString(``),
			want: want{
				status:      http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				data:        "EOF\n",
			},
		},
		{
			name: "cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			body: bytes.NewBufferString(`{"timeout": 6000}`),
			want: want{
				status:      http.StatusOK,
				contentType: "",
				data:        ``,
			},
		},
	}

	for _, tt := range tests {
		name := tt.name
		t.Run(name, test(tt.ctx, tt.body, tt.want))
	}
}
