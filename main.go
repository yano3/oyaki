package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var client http.Client
var orgSrvURL string
var quality = 90

func main() {
	orgScheme := os.Getenv("OYAKI_ORIGIN_SCHEME")
	orgHost := os.Getenv("OYAKI_ORIGIN_HOST")
	if orgScheme == "" {
		orgScheme = "https"
	}
	orgSrvURL = orgScheme + "://" + orgHost

	if q := os.Getenv("OYAKI_QUALITY"); q != "" {
		quality, _ = strconv.Atoi(q)
	}

	http.HandleFunc("/", proxy)
	http.ListenAndServe(":8080", nil)
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
		return
	}

	orgRes, err := client.Get(orgURL.String())
	if err != nil {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		return
	}
	defer orgRes.Body.Close()

	buf, err := convert(orgRes.Body, quality)
	if err != nil {
		http.Error(w, "Image convert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	w.Write(buf.Bytes())
	buf.Reset()

	return
}
