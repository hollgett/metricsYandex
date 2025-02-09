package logger

import "net/http"

type Logger interface {
	LogAny(logMess, key string, v interface{})
	LogMess(logMess string)
	LogErr(logMess string, err error)
	RequestMiddleware(next http.Handler) http.Handler
	ResponseMiddleware(next http.Handler) http.Handler
	Flush()
}
