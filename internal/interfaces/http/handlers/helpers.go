package handlers

import "net/http"

func noCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
}
