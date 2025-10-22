package models

import "gopkg.in/gographics/imagick.v3/imagick"

type IMagickService struct {
	MagicWand *imagick.MagickWand
}
