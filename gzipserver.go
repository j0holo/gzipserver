// Copyright 2017 Ruben Klink. All rights reserved.
// See LICENSE for license details.

package gzipserver

import (
	"compress/flate"
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

// Response defines the ContentType and Content for the response on a request.
type Response struct {
	Content []byte
}

// ServeHTTP implements http.Handler to provide an automatic Content-Type and
// gzip compression.
func (res *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compressedResponse(w, r, res.Content)
}

// compressedResponse compresses the response if Accept-Encoding contains gzip or deflate.
func compressedResponse(w http.ResponseWriter,
	r *http.Request,
	response []byte) {
	encoding := r.Header.Get("Accept-Encoding")
	w.Header().Set("Content-Type", http.DetectContentType(response))
	if strings.Contains(encoding, "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		_, err := gzw.Write(response)
		if err != nil {
			log.Fatal(err)
		}

		if err := gzw.Close(); err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(encoding, "deflate") {
		w.Header().Set("Content-Encoding", "deflate")
		flw, err := flate.NewWriter(w, 1)
		if err != nil {
			log.Fatal(err)
		}

		_, err = flw.Write(response)
		if err != nil {
			log.Fatal(err)
		}

		if err = flw.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		w.Header().Set("Content-Encoding", "deflate")
		// w.WriteHeader() will be implicitly called if WriteHeader() is not
		// been called before Write(). 200 OK is assumed.
		if _, err := w.Write(response); err != nil {
			log.Fatal(err)
		}
	}
}
