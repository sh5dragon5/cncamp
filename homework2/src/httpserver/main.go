package main

import (
	"flag"
	"fmt"
	"github.com/felixge/httpsnoop"
	"github.com/golang/glog"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	flag.Set("v", "2")
	glog.V(2).Info("Starting http server...")

	log.Printf("Starting http server...")

	mux := &http.ServeMux{}

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/header", headerHandler)

	var handler http.Handler = mux
	handler = logRequestHandler(handler)

	err := http.ListenAndServe(":5055", handler)
	if err != nil {
		log.Fatal(err)
	}

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering healthz handler")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Response Status Code - 200"))
}

func headerHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("entering header handler")
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	version := os.Getenv("VERSION");
	w.Header().Add("Version", version)

	for k, v := range w.Header() {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

type RequestInfo struct {
	uri string
	ip string
	code int
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		info := &RequestInfo {
			uri: r.URL.String(),
		}

		m := httpsnoop.CaptureMetrics(h, w, r)
		info.code = m.Code
		info.ip = getRemoteAddress(r)

		logRequestInfo(info)
	}
	return http.HandlerFunc(fn)
}