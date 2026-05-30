package internal

import (
	"image"

	"github.com/disintegration/imaging"
)

type effect struct {
	blur, gamma, saturation float64
}

var defaultEffect = effect{6., 0.8, -60.}

func createBackground(img image.Image, w, h int) image.Image {
	var back image.Image

	width, height := size(img)

	dx, dy := w-width, h-height
	if w*dy < h*dx {
		back = imaging.Resize(img, w, 0, imaging.Lanczos)
	} else {
		back = imaging.Resize(img, 0, h, imaging.Lanczos)
	}

	back = imaging.CropAnchor(back, w, h, imaging.Center)
	back = imaging.AdjustSaturation(back, defaultEffect.saturation)
	back = imaging.AdjustGamma(back, defaultEffect.gamma)
	return imaging.Blur(back, defaultEffect.blur)
}
