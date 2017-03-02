// Copyright 2017 Ruben Klink. All rights reserved.
// See LICENSE for license details.

package gzipserver

import (
	"compress/gzip"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// ContentType abstracts away the int type to contain more context.
type ContentType int

// ServerConfig is a struct that is used to configure the *http.Server
// that gets returned by NewGzipServer().
type ServerConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	Handler        *http.ServeMux
	MaxHeaderBytes int
	TLSConfig      *tls.Config
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState      func(net.Conn, http.ConnState)
	ErrorLog       *log.Logger
}

// Response defines the ContentType and Content for the response on a request.
type Response struct {
	Content string
}

// NewGzipServer creates a new server with sane defaults.
func NewGzipServer(sc *ServerConfig) *http.Server {
	return &http.Server{
		Addr:           sc.Addr,
		ReadTimeout:    sc.ReadTimeout,
		WriteTimeout:   sc.WriteTimeout,
		IdleTimeout:    sc.IdleTimeout,
		Handler:        sc.Handler,
		MaxHeaderBytes: sc.MaxHeaderBytes,
		TLSConfig:      sc.TLSConfig,
		TLSNextProto:   sc.TLSNextProto,
		ConnState:      sc.ConnState,
		ErrorLog:       sc.ErrorLog,
	}
}

// ServeHTTP implements http.Handler to provide an automatic Content-Type and
// gzip compression.
func (res *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compressedResponse(w, r, []byte(res.Content))
}

// compressedResponse compresses the response if Accept-Encoding contains gzip.
func compressedResponse(w http.ResponseWriter,
	r *http.Request,
	response []byte) {
	encoding := r.Header.Get("Accept-Encoding")
	if strings.Contains(encoding, "gzip") {
		w.Header().Set("Content-Type", http.DetectContentType(response))
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		_, err := gzw.Write(response)
		if err != nil {
			log.Fatal(err)
		}

		if err := gzw.Close(); err != nil {
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
