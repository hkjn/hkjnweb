// +build !appengine

// main.go holds bits needed when not serving hkjnweb on appengine.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"hkjn.me/hkjnweb"
)

// registerStatic registers handler for static files.
func registerStatic(dir string) {
	if dir == "" {
		dir = "static"
	}
	path := fmt.Sprintf("/%s/", dir)
	h := http.StripPrefix(
		path,
		http.FileServer(http.Dir(dir)))
	http.Handle(path, h)
}

func main() {
	flag.Parse()
	hkjnweb.Register(os.Getenv("PROD") != "")
	registerStatic(os.Getenv("STATIC_DIR"))

	addr := os.Getenv("BIND_ADDR")
	if addr == "" {
		addr = ":12345"
	}
	if os.Getenv("SERVE_HTTP") != "" {
		log.Printf("webserver serving HTTP on %s..\n", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatalf("FATAL: http server exited: %v\n", err)
		}
	} else {
		log.Println("Since SERVE_HTTP isn't set, we should serve https")
		certFile := os.Getenv("HTTPS_CERT_FILE")
		keyFile := os.Getenv("HTTPS_KEY_FILE")
		if certFile == "" || keyFile == "" {
			log.Fatalln("FATAL: No HTTPS_CERT_FILE or HTTPS_KEY_FILE specified.")
		}
		log.Printf("webserver serving HTTPS on %s..\n", addr)
		err := http.ListenAndServeTLS(addr, certFile, keyFile, nil)
		if err != nil {
			log.Fatalf("FATAL: https server exited: %v\n", err)
		}
	}
}
