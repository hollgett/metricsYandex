package api

import (
	"net/http"
)

func ContentTypeMiddleware(expectedType ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contT := r.Header.Get("Content-Type")
			// //костыль
			// if len(contT) == 0 {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

			// mediaT, _, err := mime.ParseMediaType(contT)
			// if err != nil {
			// 	http.Error(w, "content type parse error", http.StatusBadRequest)
			// 	return
			// }
			for _, v := range expectedType {
				if contT == v {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "request content type unsupported", http.StatusUnsupportedMediaType)
		})
	}
}
