package internal

import (
	"errors"
	"image"
	"strconv"
	"strings"
)

type padding [4]int

func ParsePadding(str string) (pad padding, err error) {
	err = errors.New("Invalid padding.")

	pads := strings.Fields(str)
	if len(pads) < 1 || len(pads) > 4 {
		return
	}

	for i, idx := range zipZap(4, len(pads)) {
		padVal, pErr := strconv.ParseUint(pads[idx], 10, 64)
		if pErr != nil {
			return
		}
		pad[i] = int(padVal)
	}

	return pad, nil
}

func PaddingFromThresholdRatio(img image.Image, rat ratio) padding {
	srcWidth, srcHeight := size(img)

	var dx, dy int

	if srcWidth > srcHeight {
		dstHeight := int(rat[1] * float64(srcWidth))

		if dstHeight <= srcHeight {
			return padding{}
		}

		dy = dstHeight - srcHeight
	} else {
		dstWidth := int(rat[0] * float64(srcHeight))

		if dstWidth <= srcWidth {
			return padding{}
		}

		dx = dstWidth - srcWidth
	}

	dx, dy = dx/2, dy/2

	return padding{dy, dx, dy, dx}
}
