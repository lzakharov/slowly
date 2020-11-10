// +build integration

package it

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lzakharov/slowly/slowly"
)

var (
	addr    string
	timeout time.Duration
	app     *slowly.App
	client  *http.Client
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestApp(t *testing.T) {
	type want struct {
		status      int
		contentType string
		data        string
	}

	test := func(ctx context.Context, body io.Reader, want want) func(*testing.T) {
		return func(t *testing.T) {
			var url = "http://" + addr + "/api/slow"

			resp, err := client.Post(url, "application/json", body)
			if err != nil {
				t.Errorf("failed to make a request: %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != want.status {
				t.Errorf("wrong status code: got %v want %v", resp.StatusCode, want.status)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != want.contentType {
				t.Errorf("wrong content type: got %v want %v", contentType, want.contentType)
			}

			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read the response body: %v", err)
			}
			data := string(raw)
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
			body: bytes.NewBufferString(`{"timeout": 1000}`),
			want: want{
				status:      http.StatusOK,
				contentType: "application/json",
				data:        `{"status":"ok"}`,
			},
		},
		{
			name: "sad path",
			ctx:  context.Background(),
			body: bytes.NewBufferString(`{"timeout": 3000}`),
			want: want{
				status:      http.StatusBadRequest,
				contentType: "application/json",
				data:        `{"error":"timeout too long"}`,
			},
		},
	}

	for _, tt := range tests {
		name := tt.name
		t.Run(name, test(tt.ctx, tt.body, tt.want))
	}
}

func setup() {
	var ok bool
	addr, ok = os.LookupEnv("SLOWLY_TEST_ADDR")
	if !ok {
		addr = "localhost:9090"
	}
	timeout = 2 * time.Second

	app = slowly.NewApp(addr, timeout)
	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	client = &http.Client{Timeout: 5 * time.Second}

}

func shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := app.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
	cancel()
}
