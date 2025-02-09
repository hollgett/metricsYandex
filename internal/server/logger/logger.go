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

type logger struct {
	l *zap.Logger
}

func New() (Logger, error) {
	lvl, err := zap.ParseAtomicLevel("info")
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	lg, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &logger{
		l: lg,
	}, nil
}

func (log *logger) LogAny(logMess, key string, v interface{}) {
	log.l.Info(logMess, zap.Any(key, v))
}

func (log *logger) LogMess(logMess string) {
	log.l.Info(logMess)
}

func (log *logger) LogErr(logMess string, err error) {
	log.l.Info(logMess, zap.Error(err))
}

func (log *logger) RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.l.Info("Request start", zap.String("uri", r.RequestURI), zap.String("method", r.Method))

		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.l.Info("Request complete", zap.String("uri", r.RequestURI), zap.String("method", r.Method), zap.Duration("duration", duration))

	})
}

func (log *logger) ResponseMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		responseData := responseData{}

		lr := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData,
		}

		next.ServeHTTP(&lr, r)
		log.l.Info("Response data", zap.Int("status", responseData.status), zap.Int("size", responseData.size))
	}

	return http.HandlerFunc(fn)
}

func (log *logger) Flush() {
	log.LogMess("log flush")
	log.l.Sync()
}
