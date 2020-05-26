package main

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
)

func convert(src io.Reader, q int) (*bytes.Buffer, error) {
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	if format == "png" {
		bg := image.NewUniform(image.White)
		dstImg := image.NewRGBA(img.Bounds())

		draw.Draw(dstImg, dstImg.Bounds(), bg, image.ZP, draw.Src)
		draw.Draw(dstImg, dstImg.Bounds(), img, image.ZP, draw.Over)

		if err := jpeg.Encode(buf, dstImg, &jpeg.Options{Quality: q}); err != nil {
			return nil, err
		}
	}

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}

	return buf, nil
}
