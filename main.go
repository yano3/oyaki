package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"go.elastic.co/apm/module/apmhttp"
)

var client http.Client
var orgSrvURL string
var quality = 90
var version = ""

func main() {
	var ver bool

	flag.BoolVar(&ver, "version", false, "show version")
	flag.Parse()

	if ver {
		fmt.Printf("oyaki %s\n", getVersion())
		return
	}

	orgScheme := os.Getenv("OYAKI_ORIGIN_SCHEME")
	orgHost := os.Getenv("OYAKI_ORIGIN_HOST")
	if orgScheme == "" {
		orgScheme = "https"
	}
	orgSrvURL = orgScheme + "://" + orgHost

	if q := os.Getenv("OYAKI_QUALITY"); q != "" {
		quality, _ = strconv.Atoi(q)
	}

	log.Printf("starting oyaki %s\n", getVersion())
	http.HandleFunc("/", proxy)
	http.ListenAndServe(":8080", apmhttp.Wrap(http.DefaultServeMux))
}

func proxy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	if path == "/" {
		fmt.Fprintln(w, "Oyaki lives!")
		return
	}

	orgURL, err := url.Parse(orgSrvURL + path)
	if err != nil {
		http.Error(w, "Invalid origin URL", http.StatusBadRequest)
		log.Printf("Invalid origin URL. %v\n", err)
		return
	}

	req, err := http.NewRequest("GET", orgURL.String(), nil)
	if err != nil {
		http.Error(w, "Request Failed", http.StatusInternalServerError)
		log.Printf("Request Failed. %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "oyaki")

	xff := r.Header.Get("X-Forwarded-For")
	if len(xff) > 1 {
		req.Header.Set("X-Forwarded-For", xff)
	}

	orgRes, err := client.Do(req)
	if err != nil {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		log.Printf("Get origin failed. %v\n", err)
		return
	}
	defer orgRes.Body.Close()
	if orgRes.StatusCode != http.StatusOK {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		log.Printf("Get origin failed. %v\n", orgRes.Status)
		return
	}

	ct := orgRes.Header.Get("Content-Type")
	cl := orgRes.Header.Get("Content-Length")

	if ct != "image/jpeg" {
		w.Header().Set("Content-Type", ct)
		if cl != "" {
			w.Header().Set("Content-Length", cl)
		}

		_, err := io.Copy(w, orgRes.Body)
		if err != nil {
			// ignore already close client.
			if !errors.Is(err, syscall.EPIPE) {
				http.Error(w, "Read origin body failed", http.StatusInternalServerError)
				log.Printf("Read origin body failed. %v\n", err)
			}
		}
		return
	}

	buf, err := convert(orgRes.Body, quality)
	if err != nil {
		http.Error(w, "Image convert failed", http.StatusInternalServerError)
		log.Printf("Image convert failed. %v\n", err)
		return
	}
	defer buf.Reset()

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	if _, err := io.Copy(w, buf); err != nil {
		// ignore already close client.
		if !errors.Is(err, syscall.EPIPE) {
			http.Error(w, "Write responce failed", http.StatusInternalServerError)
			log.Printf("Write responce  failed. %v\n", err)
		}
	}
}

func getVersion() string {
	if version != "" {
		return version
	}

	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}
	return i.Main.Version
}
