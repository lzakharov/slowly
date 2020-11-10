package slowly

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lzakharov/slowly/pkg/web"
)

// Slow is a handler that launches emulation of a long running process. If the request context ends before the process
// terminates, the handler stops the process and returns immediately.
func Slow(w http.ResponseWriter, r *http.Request) {
	payload := SlowRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeout := time.Duration(payload.Timeout) * time.Millisecond
	process := time.NewTimer(timeout)

	select {
	case <-process.C:
	case <-r.Context().Done():
		process.Stop()
		return
	}

	web.WriteJSON(w, http.StatusOK, SlowResponse{Status: StatusOK})
}
