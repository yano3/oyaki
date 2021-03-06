package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

func convert(src io.Reader, q int) (*bytes.Buffer, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}

	return buf, nil
}
