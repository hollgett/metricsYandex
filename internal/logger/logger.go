package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size

	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

var Log *zap.Logger = zap.NewNop()

func InitLogger() error {
	lvl, err := zap.ParseAtomicLevel("info")
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	lg, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = lg
	return nil
}

func LogInfo(logMess string, v ...zap.Field) {
	Log.Info(logMess, v...)
}

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		LogInfo("Request start", zap.String("uri", r.RequestURI), zap.String("method", r.Method))

		next.ServeHTTP(w, r)
		duration := time.Since(start)
		LogInfo("Request complete", zap.String("uri", r.RequestURI), zap.String("method", r.Method), zap.Duration("duration", duration))

	})
}

func ResponseMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		responseData := responseData{}

		lr := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData,
		}

		next.ServeHTTP(&lr, r)
		LogInfo("Response data", zap.Int("status", responseData.status), zap.Int("size", responseData.size))
	}

	return http.HandlerFunc(fn)
}
