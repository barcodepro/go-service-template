package app

import (
	"context"
	"encoding/json"
	"github.com/barcodepro/go-service-template/service/internal/log"
	"github.com/barcodepro/go-service-template/service/store"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
	"strings"
	"time"
)

type server struct {
	httpserver *http.Server
	router     *httprouter.Router
	store      *store.Store
	config     *Config
}

type ctxKey int8

const (
	ctxKeyRequestID ctxKey = iota
)

func newServer(config *Config, store *store.Store) *server {
	s := &server{
		router: httprouter.New(),
		store:  store,
		config: config,
	}

	corsPolicy := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(config.AllowedOrigins, ","),
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodGet},
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})

	s.httpserver = &http.Server{
		Addr:              config.ListenAddress.String(),
		Handler:           corsPolicy.Handler(s),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	s.configureRouter()
	return s
}

// ServeHTTP ...
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	m := newMiddlewareStack()
	m.Use(s.setupContext)
	m.Use(s.logRequest)

	// non-protected endpoints
	s.router.Handler(http.MethodGet, "/info", m.Wrap(s.handleInfo()))
	s.router.Handler(http.MethodGet, "/healthz", m.Wrap(s.handleHealthz()))
	s.router.Handler(http.MethodGet, "/readyz", m.Wrap(s.handleReadyz()))

	// authenticate requests and authorize requests
	//m.Use(s.authenticateRequest)
	//m.Use(s.authorizeRequest)

	// example endpoint
	//s.router.Handler(http.MethodGet, "/accounts/:account_id", m.Wrap(s.handleAccountsFind()))
}

func (s *server) setupContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add request id
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), ctxKeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received := time.Now()
		method := r.Method
		requestURI := r.RequestURI
		requestID := r.Context().Value(ctxKeyRequestID)
		agent := r.UserAgent()
		referer := r.Referer()
		proto := r.Proto
		remoteAddr := r.RemoteAddr
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		var logger = log.New()
		logger.Info().
			Time("req_received", received).
			Str("req_method", method).
			Str("req_uri", requestURI).
			Str("agent", agent).
			Str("referer", referer).
			Str("proto", proto).
			Str("remote_addr", remoteAddr).
			Int("status", rw.code).
			Str("status_text", http.StatusText(rw.code)).
			Dur("duration", time.Since(received)).
			Interface("request_id", requestID).
			Msg("")
	})
}

/*
 * HTTP handlers
 */

func (s *server) handleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "i am simple REST API service")
	})
}

func (s *server) handleHealthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"status": "healthy"})
	})
}

func (s *server) handleReadyz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"status": "ready"})
	})
}

/*
 * Respond/Error wrappers
 */
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, e error) {
	type errorsWrapper struct {
		Errors interface{} `json:"errors"`
	}
	w.WriteHeader(code)
	id := r.Context().Value(ctxKeyRequestID)
	log.Errorf("request failed: %s /* request_id:%s */", e, id)

	err := json.NewEncoder(w).Encode(errorsWrapper{Errors: map[string]string{"title": e.Error()}})
	if err != nil {
		log.Errorln("json encode error: ", err)
		_, err2 := w.Write([]byte(e.Error()))
		if err2 != nil {
			log.Errorln("response write error: ", err2)
		}
	}
}

// respond is the wrapper for successful responses
func (s *server) respond(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	type dataWrapper struct {
		Data interface{} `json:"data"`
	}
	w.WriteHeader(code)
	if data != nil {
		ierr := json.NewEncoder(w).Encode(dataWrapper{Data: data})
		if ierr != nil {
			log.Errorf("json encode error: %s", ierr)
		}
	}
}
