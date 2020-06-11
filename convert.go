package main

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
)

func convert(src io.Reader, q int) (*bytes.Buffer, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}

	return buf, nil
}
