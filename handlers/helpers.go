package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// ExtractID extracts an integer ID from either query param "id" or the last path segment
func ExtractID(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		return strconv.Atoi(idStr)
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return 0, errors.New("missing ID")
	}

	return strconv.Atoi(parts[len(parts)-1])
}
