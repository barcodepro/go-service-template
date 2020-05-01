package app

import (
	"fmt"
	"github.com/barcodepro/go-service-template/service/store"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_TestServiceHandlers(t *testing.T) {
	st, _ := store.TestStore(t, store.TestStoreConfig(t))
	s := newServer(&Config{}, st)
	assert.NotNil(t, s)

	var testcases = []struct {
		path string
		want int
	}{
		{path: "/info", want: http.StatusOK},
		{path: "/healthz", want: http.StatusOK},
		{path: "/readyz", want: http.StatusOK},
	}

	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.want, rec.Code)
		})
	}
}

func TestServer_error(t *testing.T) {
	st, _ := store.TestStore(t, store.TestStoreConfig(t))
	s := newServer(&Config{}, st)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.error(rec, req, http.StatusBadRequest, fmt.Errorf("test error"))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServer_respond(t *testing.T) {
	st, _ := store.TestStore(t, store.TestStoreConfig(t))
	s := newServer(&Config{}, st)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.respond(rec, req, http.StatusOK, "test")
	assert.Equal(t, http.StatusOK, rec.Code)
}
