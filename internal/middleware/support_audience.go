package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SupportAudienceAuthConfig struct {
	HeaderName   string
	TokenPrefix  string
	RequiredRole string
	SkipPaths    []string
}

func NewSupportAudienceAuthMiddleware(config *SupportAudienceAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range config.SkipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}
			headerName := config.HeaderName
			if headerName == "" {
				headerName = "Authorization"
			}
			token := r.Header.Get(headerName)
			if token == "" {
				http.Error(w, "support_audience: missing authentication token", http.StatusUnauthorized)
				return
			}
			prefix := config.TokenPrefix
			if prefix == "" {
				prefix = "Bearer "
			}
			if !strings.HasPrefix(token, prefix) {
				http.Error(w, "support_audience: invalid token format", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type SupportAudienceRateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
	KeyFunc     func(*http.Request) string
}

func NewSupportAudienceRateLimiter(config *SupportAudienceRateLimitConfig) func(http.Handler) http.Handler {
	counters := make(map[string]int)
	resetTimes := make(map[string]time.Time)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			if config.KeyFunc != nil {
				key = config.KeyFunc(r)
			}
			now := time.Now()
			if resetTime, ok := resetTimes[key]; ok && now.After(resetTime) {
				counters[key] = 0
				resetTimes[key] = now.Add(config.Window)
			}
			if _, ok := resetTimes[key]; !ok {
				resetTimes[key] = now.Add(config.Window)
			}
			counters[key]++
			if counters[key] > config.MaxRequests {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(config.Window.Seconds())))
				http.Error(w, "support_audience: rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type SupportAudienceLogEntry struct {
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	RemoteAddr string
	UserAgent  string
}

type SupportAudienceResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *SupportAudienceResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewSupportAudienceLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &SupportAudienceResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)
			_ = SupportAudienceLogEntry{
				Method:     r.Method,
				Path:       r.URL.Path,
				StatusCode: wrapped.statusCode,
				Duration:   time.Since(start),
				RemoteAddr: r.RemoteAddr,
				UserAgent:  r.UserAgent(),
			}
		})
	}
}

type SupportAudienceCORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

func NewSupportAudienceCORS(config *SupportAudienceCORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := false
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
				}
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NewSupportAudienceRecovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "support_audience: internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func NewSupportAudienceTimeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()
			select {
			case <-done:
				return
			case <-time.After(timeout):
				http.Error(w, "support_audience: request timed out", http.StatusGatewayTimeout)
			}
		})
	}
}

type SupportAudienceRequestIDConfig struct {
	HeaderName string
	Generator  func() string
}

func NewSupportAudienceRequestID(config *SupportAudienceRequestIDConfig) func(http.Handler) http.Handler {
	headerName := config.HeaderName
	if headerName == "" {
		headerName = "X-Request-ID"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(headerName)
			if requestID == "" && config.Generator != nil {
				requestID = config.Generator()
			}
			if requestID == "" {
				requestID = fmt.Sprintf("support_audience-%d", time.Now().UnixNano())
			}
			w.Header().Set(headerName, requestID)
			next.ServeHTTP(w, r)
		})
	}
}
