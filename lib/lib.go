package lib

import (
	"fmt"
	"image"
	"path"
	"strings"

	"github.com/disintegration/imaging"
)

func OpenDoSave(
	in, out string,
	do func(image.Image) (image.Image, error),
) error {
	if out == "" {
		out = defaultOut(in)
	}

	src, err := imaging.Open(in, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	dst, err := do(src)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	if err := imaging.Save(dst, out); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func defaultOut(in string) string {
	ext := path.Ext(in)
	base := strings.TrimSuffix(in, ext)
	return fmt.Sprintf("%s(blurpad)%s", base, ext)
}
