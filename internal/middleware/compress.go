// Package middleware provides HTTP middleware for compressing HTTP responses and decompressing HTTP request bodies using gzip encoding.
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by
// the `compressWriter`. This header map is a reference
// to the original http.ResponseWriter header map.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// This method compresses the data using gzip before writing it to
// the underlying ResponseWriter. It returns the number of bytes
// written and any write error encountered.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip.Writer, ensuring that all data is flushed to the underlying http.ResponseWriter.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads compressed data from the underlying gzip.Reader.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes the compressReader by first closing the underlying io.ReadCloser and then closing the gzip.Reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// CompressorMiddleware is a middleware function that compresses HTTP responses and decompresses request bodies based on the "Accept-Encoding" and "Content-Encoding" headers.
func CompressorMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}

			h.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
