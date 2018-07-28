package zip

import (
	"compress/gzip"
	"net/http"
)

func Middleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Content-Encoding", "gzip")
		next(&zipWriter{w, http.Header{}, 0}, r)
	}
}

type zipWriter struct {
	http.ResponseWriter
	header http.Header
	status int
}

func (w *zipWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	var writer, err = gzip.NewWriterLevel(w.ResponseWriter, gzip.BestCompression)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

func (w *zipWriter) Header() http.Header {
	return w.header
}

func (w *zipWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	var header = w.ResponseWriter.Header()
	for k, v := range w.header {
		for _, val := range v {
			header.Add(k, val)
		}
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
