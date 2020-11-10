package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	type want struct {
		status      int
		contentType string
		data        string
	}

	test := func(status int, v interface{}, want want) func(*testing.T) {
		return func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, status, v)

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
		name   string
		status int
		v      interface{}
		want   want
	}{
		{
			name:   "number",
			status: http.StatusOK,
			v:      42,
			want: want{
				status:      http.StatusOK,
				contentType: "application/json",
				data:        `42`,
			},
		},
		{
			name:   "nil",
			status: http.StatusAccepted,
			v:      nil,
			want: want{
				status:      http.StatusAccepted,
				contentType: "application/json",
				data:        `null`,
			},
		},
		{
			name:   "map",
			status: http.StatusBadRequest,
			v:      map[string]interface{}{"n": 42},
			want: want{
				status:      http.StatusBadRequest,
				contentType: "application/json",
				data:        `{"n":42}`,
			},
		},
		{
			name:   "unsupported complex number",
			status: http.StatusBadRequest,
			v:      1 + 2i,
			want: want{
				status:      http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				data:        "json: unsupported type: complex128\n",
			},
		},
	}

	for _, tt := range tests {
		name := tt.name
		t.Run(name, test(tt.status, tt.v, tt.want))
	}
}
