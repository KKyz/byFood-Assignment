package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type urlErrResp struct {
	Error string `json:"error"`
}

func TestProcessURLHandler_EdgeCases(t *testing.T) {
	router := chi.NewRouter()
	router.Post("/process-url", ProcessURLHandler)

	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantSubstr string
		wantErr    string
	}{
		{
			name:       "all operation",
			body:       `{"url":"https://BYFOOD.com/food-EXPeriences?query=abc/","operation":"all"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"processed_url":"https://www.byfood.com/food-experiences"`,
		},
		{
			name:       "canonical keeps case and removes query",
			body:       `{"url":"https://BYFOOD.com/food-EXPeriences?query=abc/","operation":"canonical"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"processed_url":"https://BYFOOD.com/food-EXPeriences"`,
		},
		{
			name:       "redirection forces www and lowercases (keeps query)",
			body:       `{"url":"https://BYFOOD.com/food-EXPeriences?query=ABC/","operation":"redirection"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"processed_url":"https://www.byfood.com/food-experiences?query=abc/"`,
		},

		// --- negatives ---
		{
			name:       "invalid json",
			body:       `{"url":`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid JSON body",
		},
		{
			name:       "missing url",
			body:       `{"operation":"all"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "`url` is required",
		},
		{
			name:       "missing operation",
			body:       `{"url":"https://byfood.com/x"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "`operation` is required",
		},
		{
			name:       "bad operation",
			body:       `{"url":"https://byfood.com/x","operation":"nope"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "`operation` must be one of: canonical, redirection, all",
		},
		{
			name:       "invalid url string",
			body:       `{"url":"not a url","operation":"all"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid URL (must include scheme and host)",
		},
		{
			name:       "url missing scheme",
			body:       `{"url":"byfood.com/abc","operation":"all"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid URL (must include scheme and host)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/process-url", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status got %d want %d; body=%s", rr.Code, tc.wantStatus, rr.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				if !strings.Contains(rr.Body.String(), tc.wantSubstr) {
					t.Fatalf("body missing substring:\nwant: %s\ngot:  %s", tc.wantSubstr, rr.Body.String())
				}
				return
			}

			var er urlErrResp
			if err := json.Unmarshal(rr.Body.Bytes(), &er); err != nil {
				t.Fatalf("invalid json: %v body=%s", err, rr.Body.String())
			}
			if er.Error != tc.wantErr {
				t.Fatalf("error got %q want %q", er.Error, tc.wantErr)
			}
		})
	}
}
