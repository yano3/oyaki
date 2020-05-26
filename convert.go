package main

import (
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
)

func convert(src io.Reader, w io.Writer, q int) error {
	img, format, err := image.Decode(src)
	if err != nil {
		return err
	}

	if format == "png" {
		bg := image.NewUniform(image.White)
		dstImg := image.NewRGBA(img.Bounds())

		draw.Draw(dstImg, dstImg.Bounds(), bg, image.ZP, draw.Src)
		draw.Draw(dstImg, dstImg.Bounds(), img, image.ZP, draw.Over)

		if err := jpeg.Encode(w, dstImg, &jpeg.Options{Quality: q}); err != nil {
			return err
		}
	}

	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	return nil
}
