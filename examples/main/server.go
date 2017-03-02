// Copyright 2017 Ruben Klink. All rights reserved.
// See LICENSE for license details.

package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/j0holo/gzipserver"
)

func main() {
	mux := http.NewServeMux()

	newHandler := &gzipserver.Response{
		Content: "<h1>Welcome to this webpage.</h1>",
	}

	m.Handle("/", newHandler)

	tlsConf := &tls.Config{
		PreferServerCipherSuites: true,
	}

	srv := gzipserver.NewGzipServer(&gzipserver.ServerConfig{
		Addr:           "127.0.0.1:8080",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    20 * time.Second,
		Handler:        m,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConf,
		TLSNextProto:   nil,
		ConnState:      nil,
		ErrorLog:       nil,
	})

	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
