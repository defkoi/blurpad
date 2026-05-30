package internal

import (
	"image"

	"github.com/disintegration/imaging"
)

func Process(img image.Image, pad padding) image.Image {
	var dst image.Image

	w, h := size(img)

	back := createBackground(img, w+pad[1]+pad[3], h+pad[0]+pad[2])
	dst = imaging.Paste(back, img, image.Pt(pad[3], pad[0]))

	return dst
}
