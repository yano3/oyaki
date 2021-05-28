package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	orgRes, err := client.Do(req)
	if err != nil {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		log.Printf("Get origin failed. %v\n", err)
		return
	}
	if orgRes.StatusCode != http.StatusOK {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		log.Printf("Get origin failed. %v\n", orgRes.Status)
		return
	}
	defer orgRes.Body.Close()

	ct := orgRes.Header.Get("Content-Type")
	cl := orgRes.Header.Get("Content-Length")

	if ct != "image/jpeg" {
		body, err := ioutil.ReadAll(orgRes.Body)
		if err != nil {
			http.Error(w, "Read origin body failed", http.StatusInternalServerError)
			log.Printf("Read origin body failed. %v\n", err)
			return
		}

		w.Header().Set("Content-Type", ct)
		if cl != "" {
			w.Header().Set("Content-Length", cl)
		}

		w.Write(body)
		return
	}

	buf, err := convert(orgRes.Body, quality)
	if err != nil {
		http.Error(w, "Image convert failed", http.StatusInternalServerError)
		log.Printf("Image convert failed. %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	w.Write(buf.Bytes())
	buf.Reset()

	return
}
