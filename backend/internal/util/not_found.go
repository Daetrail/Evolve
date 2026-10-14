package util

import "net/http"

func NotFound(w http.ResponseWriter, _ *http.Request) {
	SendError(w, http.StatusNotFound, "route not found")
}