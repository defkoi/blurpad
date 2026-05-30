package internal

import (
	"image"
	"iter"
)

func size(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func zipZap(count, threshold int) iter.Seq2[int, int] {
	if count <= 0 || threshold <= 0 {
		return func(yield func(int, int) bool) {}
	}

	period := 2 * (threshold - 1)
	if period == 0 {
		period = 1
	}

	return func(yield func(int, int) bool) {
		for i := range count {
			pos := i % period
			var value int
			if pos < threshold {
				value = pos
			} else {
				value = 2*(threshold-1) - pos
			}
			if !yield(i, value) {
				return
			}
		}
	}
}
