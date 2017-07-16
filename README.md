# GzipServer

GzipServer is a small wrapper around the http.Handler interface to provide automatic gzip compression when a client allows gzip or deflate compression in Accept-Encoding header field. Gzip compression is the default, deflate compression will only be used if gzip is not included in Accept-Encoding.

## Get the package:

You can get the package via:

`$ go get github.com/j0holo/gzipServer`

Don't foget to set your GOPATH before running go get.

## Public Functions and Structs:

Response is a struct that contains the content that will be served to the client on a request. The Response struct contains a method (ServerHTTP) which implements the http.Hanlder interface.

```
type Response struct {
	Content []byte
}
```

ServeHTTP implements http.Handler to provide an automatic Content-Type and gzip compression. It will be called automatically when a Response is added to http.ServeMux via (*ServeMux) Handle().

```
func (res *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compressedResponse(w, r, res.Content)
}
```

**Note:** (*ServeMux) HandleFunc() can not be used because it doesn't use the gzipserver.Response struct.

## Example:

You can create a self-signed certificate with openssl (on Linux):

`openssl genrsa -out server.key 2048`

`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365`

```
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

	index := &gzipserver.Response{
		Content: []byte("<h1>Welcome to this webpage.</h1>"),
	}

	mux.Handle("/", index)

	tlsConf := &tls.Config{
		PreferServerCipherSuites: true,
	}

	srv := http.Server{
		Addr:           "127.0.0.1:8080",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    20 * time.Second,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConf,
		TLSNextProto:   nil,
		ConnState:      nil,
		ErrorLog:       nil,
	}

	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
```

## Contribute:

If you have any suggestions, issues, or feedback open a issue or send a pull request.

## License:

See [LICENSE](LICENSE).
