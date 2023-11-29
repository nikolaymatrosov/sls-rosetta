package main

import (
	"io"

	"github.com/davidbyttow/govips/v2/vips"
)

const (
	width  = 100
	height = 100
)

func Thumbnail(input io.Reader, output io.WriteCloser) error {
	image, err := vips.NewImageFromReader(input)
	if err != nil {
		return err
	}
	if err = image.Thumbnail(width, height, vips.InterestingCentre); err != nil {
		return err
	}
	data, _, err := image.ExportNative()
	if err != nil {
		return err
	}
	defer output.Close()
	if _, err = output.Write(data); err != nil {
		return err
	}
	return nil
}
