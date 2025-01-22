package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.gw.Write(p)
}

func (cw *compressWriter) Close() error {
	return cw.gw.Close()
}

type compressReader struct {
	io.Reader
	gr *gzip.Reader
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.gr.Read(p)
}

func (cr *compressReader) Close() error {
	return cr.gr.Close()
}

func CompressMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aliasW := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			aliasRespWriter := &compressWriter{
				ResponseWriter: w,
				gw:             gzip.NewWriter(w),
			}
			defer aliasRespWriter.Close()
			w.Header().Set("Content-Encoding", "gzip")
			aliasW = aliasRespWriter
		}
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gReader, err := gzip.NewReader(r.Body)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "decompress", err)
				return
			}
			aliasReqReader := &compressReader{
				Reader: r.Body,
				gr:     gReader,
			}
			defer aliasReqReader.Close()
			r.Body = aliasReqReader
		}
		next.ServeHTTP(aliasW, r)
	})

}

func ContentTypeMiddleware(expectedType ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contT := r.Header.Get("Content-Type")
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
