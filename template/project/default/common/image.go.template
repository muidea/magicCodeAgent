package common

import (
	cd "github.com/muidea/magicCommon/def"
	fc "github.com/muidea/magicFile/common"
)

const (
	FilterImage = "/image/query/"
	QueryImage  = "/image/query/:id"
	UpdateImage = "/image/update/:id"
	DeleteImage = "/image/destroy/:id"
)

const ImageModule = "/module/image"

type ImageResult struct {
	cd.Result
	Image *fc.FileView `json:"image"`
}

type ImageListResult struct {
	cd.Result
	Total int64          `json:"total"`
	Image []*fc.FileView `json:"image"`
}

type ImageStatisticResult struct {
	ImageListResult
}
