// BUG(sparky) When using GZIP browsers are much more difficult about the Content-Type headers. If improperly set the file will often download instead of loading properly.

// Provides a way to wrap handlers with on the fly GZIP compression
//
// See https:// groups.google.com/forum/?fromgroups=#!topic/golang-nuts/eVnTcMwNVjM
package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Wrap the ResponseWriter with another writer.
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Map the write operation to the Writer not the ResponseWriter.
func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Returns the same handler, but the normal ResponseWriter is replaced with a GzipResponseWriter where the write operation is passed through a gzip Writer from the compress/gzip package
func GzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if supported
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		// Set encoding
		w.Header().Set("Content-Encoding", "gzip")

		// Build new writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Call Handler
		h.ServeHTTP(GzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
