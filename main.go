package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

	err = convert(orgRes.Body, w)
	if err != nil {
		http.Error(w, "Image onvert failed", http.StatusInternalServerError)
		return
	}

	return
}

func convert(src io.Reader, w io.Writer) error {
	img, format, err := image.Decode(src)
	if err != nil {
		return err
	}

	if format == "png" {
		bg := image.NewUniform(image.White)
		dstImg := image.NewRGBA(img.Bounds())

		draw.Draw(dstImg, dstImg.Bounds(), bg, image.ZP, draw.Src)
		draw.Draw(dstImg, dstImg.Bounds(), img, image.ZP, draw.Over)

		if err := jpeg.Encode(w, dstImg, &jpeg.Options{Quality: quality}); err != nil {
			return err
		}
	}

	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	return nil
}
