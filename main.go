package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"iter"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

type config struct {
	in, out, padding, ratio, target string
}

type padding [4]int

type ratio struct {
	portrait, album float64
}

type blur struct {
	blur, gamma, saturation float64
}

const (
	name = "blurpad"

	instagramRatio = "4:5 1:1.91"

	instagramTarget = "instagram"
)

var targets = map[string]ratio{
	instagramTarget: {4. / 5., 1. / 1.91},
	/* add more or read from */
}

func main() {
	cfg := parseConfig()

	src, err := imaging.Open(cfg.in, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatal(err)
	}

	var pad padding
	switch cfg.padding {
	case "target":
		target, ok := targets[cfg.target]
		if !ok {
			log.Fatal("unsupported target")
		}
		pad = paddingFromRatio(src, target)
	case "ratio":
		ratio, err := parseRatio(cfg.ratio)
		if err != nil {
			log.Fatal(err)
		}
		pad = paddingFromRatio(src, ratio)
	default:
		_pad, err := parsePadding(cfg.padding)
		if err != nil {
			log.Fatal(err)
		}
		pad = _pad
	}

	dst := process(src, pad)

	if err := imaging.Save(dst, cfg.out); err != nil {
		log.Fatal(err)
	}
}

func parseConfig() config {
	var in, out, pad, rat, tgt string
	flag.StringVar(&in, "i", "input.png", "input")
	flag.StringVar(&out, "o", "", "output")
	flag.StringVar(&pad, "p", "target", "padding (css|ratio|target)")
	flag.StringVar(&rat, "r", instagramRatio, "threshold ratio")
	flag.StringVar(&tgt, "t", instagramTarget, "target")
	flag.Parse()

	if out == "" {
		ext := path.Ext(in)
		base := strings.TrimSuffix(in, ext)
		out = fmt.Sprintf("%s(%s)%s", base, name, ext)
	}

	return config{in, out, pad, rat, tgt}
}

func paddingFromRatio(img image.Image, rat ratio) padding {
	srcWidth, srcHeight := size(img)

	var dx, dy int

	if srcWidth > srcHeight {
		dstHeight := int(rat.album * float64(srcWidth))

		if dstHeight <= srcHeight {
			return padding{}
		}

		dy = dstHeight - srcHeight
	} else {
		dstWidth := int(rat.portrait * float64(srcHeight))

		if dstWidth <= srcWidth {
			return padding{}
		}

		dx = dstWidth - srcWidth
	}

	dx, dy = dx/2, dy/2
	return padding{dy, dx, dy, dx}
}

func parseRatio(str string) (rat ratio, err error) {
	err = errors.New("invalid ratio")

	var (
		portNum, portDenom float64
		albNum, albDenom   float64
	)

	rats := strings.Fields(str)

	parseSingleRatio := func(i int) (float64, float64, bool) {
		rat := strings.Split(rats[i], ":")
		if len(rat) != 2 {
			return 0, 0, false
		}
		num, parseErr := strconv.ParseFloat(rat[0], 64)
		if parseErr != nil {
			return 0, 0, false
		}
		denom, parseErr := strconv.ParseFloat(rat[1], 64)
		if parseErr != nil {
			return 0, 0, false
		}
		return num, denom, true
	}

	switch len(rats) {
	case 1:
		num, denom, ok := parseSingleRatio(0)
		if !ok {
			return
		}
		portNum, portDenom, albNum, albDenom =
			num, denom, num, denom
	case 2:
		_portNum, _portDenom, ok := parseSingleRatio(0)
		if !ok {
			return
		}
		_albNum, _albDenom, ok := parseSingleRatio(1)
		if !ok {
			return
		}
		portNum, portDenom, albNum, albDenom =
			_portNum, _portDenom, _albNum, _albDenom
	default:
		return
	}

	if portNum > portDenom {
		portNum, portDenom = portDenom, portNum
	}
	if albNum > albDenom {
		albNum, albDenom = albDenom, albNum
	}

	return ratio{portNum / portDenom, albNum / albDenom}, nil
}

func parsePadding(str string) (pad padding, err error) {
	err = errors.New("invalid padding")

	pads := strings.Fields(str)
	if len(pads) < 1 || len(pads) > 4 {
		return
	}
	for i, idx := range zipZap(4, len(pads)) {
		padVal, parseErr := strconv.ParseUint(pads[idx], 10, 64)
		if parseErr != nil {
			return
		}
		pad[i] = int(padVal)
	}
	return pad, nil
}

func process(img image.Image, pad padding) image.Image {
	var dst image.Image

	w, h := size(img)

	back := createBackground(img, w+pad[1]+pad[3], h+pad[0]+pad[2])
	dst = imaging.Paste(back, img, image.Pt(pad[3], pad[0]))

	return dst
}

func createBackground(img image.Image, w, h int) image.Image {
	var defaultBlur = blur{6., 0.8, -60.}

	var back image.Image

	width, height := size(img)

	dx, dy := w-width, h-height
	if w*dy < h*dx {
		back = imaging.Resize(img, w, 0, imaging.Lanczos)
	} else {
		back = imaging.Resize(img, 0, h, imaging.Lanczos)
	}

	back = imaging.CropAnchor(back, w, h, imaging.Center)
	back = imaging.AdjustSaturation(back, defaultBlur.saturation)
	back = imaging.AdjustGamma(back, defaultBlur.gamma)
	return imaging.Blur(back, defaultBlur.blur)
}

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
