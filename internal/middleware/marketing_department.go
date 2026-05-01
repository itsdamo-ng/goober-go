package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MarketingDepartmentAuthConfig struct {
	HeaderName   string
	TokenPrefix  string
	RequiredRole string
	SkipPaths    []string
}

func NewMarketingDepartmentAuthMiddleware(config *MarketingDepartmentAuthConfig) func(http.Handler) http.Handler {
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
				http.Error(w, "marketing_department: missing authentication token", http.StatusUnauthorized)
				return
			}
			prefix := config.TokenPrefix
			if prefix == "" {
				prefix = "Bearer "
			}
			if !strings.HasPrefix(token, prefix) {
				http.Error(w, "marketing_department: invalid token format", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type MarketingDepartmentRateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
	KeyFunc     func(*http.Request) string
}

func NewMarketingDepartmentRateLimiter(config *MarketingDepartmentRateLimitConfig) func(http.Handler) http.Handler {
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
				http.Error(w, "marketing_department: rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type MarketingDepartmentLogEntry struct {
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	RemoteAddr string
	UserAgent  string
}

type MarketingDepartmentResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *MarketingDepartmentResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewMarketingDepartmentLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &MarketingDepartmentResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)
			_ = MarketingDepartmentLogEntry{
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

type MarketingDepartmentCORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

func NewMarketingDepartmentCORS(config *MarketingDepartmentCORSConfig) func(http.Handler) http.Handler {
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

func NewMarketingDepartmentRecovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "marketing_department: internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func NewMarketingDepartmentTimeout(timeout time.Duration) func(http.Handler) http.Handler {
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
				http.Error(w, "marketing_department: request timed out", http.StatusGatewayTimeout)
			}
		})
	}
}

type MarketingDepartmentRequestIDConfig struct {
	HeaderName string
	Generator  func() string
}

func NewMarketingDepartmentRequestID(config *MarketingDepartmentRequestIDConfig) func(http.Handler) http.Handler {
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
				requestID = fmt.Sprintf("marketing_department-%d", time.Now().UnixNano())
			}
			w.Header().Set(headerName, requestID)
			next.ServeHTTP(w, r)
		})
	}
}
