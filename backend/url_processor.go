package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type processURLRequest struct {
	URL       string `json:"url"`
	Operation string `json:"operation"`
}

type processURLResponse struct {
	ProcessedURL string `json:"processed_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// ProcessURLHandler godoc
// @Summary Process a URL (canonical/redirection/all)
// @Tags url
// @Accept json
// @Produce json
// @Param payload body processURLRequest true "Payload"
// @Success 200 {object} processURLResponse
// @Failure 400 {object} errorResponse
// @Router /process-url [post]
func ProcessURLHandler(w http.ResponseWriter, r *http.Request) {
	var req processURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	req.Operation = strings.TrimSpace(req.Operation)

	if req.URL == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "`url` is required"})
		return
	}
	if req.Operation == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "`operation` is required"})
		return
	}
	if !isValidOperation(req.Operation) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "`operation` must be one of: canonical, redirection, all"})
		return
	}

	parsed, err := url.Parse(req.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid URL (must include scheme and host)"})
		return
	}

	processed, err := processURL(parsed, req.Operation)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, processURLResponse{ProcessedURL: processed})
}

func processURL(u *url.URL, op string) (string, error) {
	switch op {
	case "canonical":
		out := cloneURL(u)
		applyCanonical(out)
		return out.String(), nil

	case "redirection":
		out := cloneURL(u)
		applyRedirection(out)
		return out.String(), nil

	case "all":
		out := cloneURL(u)
		applyCanonical(out)
		applyRedirection(out)
		return out.String(), nil

	default:
		return "", errors.New("unsupported operation")
	}
}

func applyCanonical(u *url.URL) {

	u.RawQuery = ""
	u.ForceQuery = false
	u.Fragment = ""

	u.Path = strings.TrimRight(u.Path, "/")
	if u.Path == "" {
		u.Path = "/"
	}
}

func applyRedirection(u *url.URL) {

	u.Host = "www.byfood.com"

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.ToLower(u.Path)
	u.RawQuery = strings.ToLower(u.RawQuery)
	u.Fragment = strings.ToLower(u.Fragment)
}

func isValidOperation(op string) bool {
	return op == "canonical" || op == "redirection" || op == "all"
}

func cloneURL(u *url.URL) *url.URL {
	clone := *u
	return &clone
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
