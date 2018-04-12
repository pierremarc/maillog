package main

import (
	"fmt"
	"log"
	"path"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func init() {
	imagick.Initialize()
}

// type thumbnail struct {
// 	orig   string
// 	size   uint
// 	result ResultString
// }

// type thumbnailer struct {
// 	notice chan thumbnail
// }

const (
	ThumbnailSmall  = 400
	ThumbnailMedium = 800
	ThumbnailLarge  = 1600
)

const divf = 1000
const divf2 = divf * divf

func thumSize(w, h, s uint) (uint, uint) {
	r := s * divf / w * divf
	// log.Printf("thumSize(%v, %v, %v) r=%v -> %v, %v", w, h, s, r, w*r/divf2, h*r/divf2)
	return w * r / divf2, h * r / divf2
}

func sizeName(p string, s uint) string {
	b := path.Base(p)
	d := path.Dir(p)
	return path.Join(d, fmt.Sprintf("thumbnails/%d", s), b)
}

func process(p string, s uint) ResultString {
	var err error
	dest := sizeName(p, s)
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImage(p)
	if err != nil {
		log.Printf("Error Reading Thumbnail(%s) %v", p, err)
		return ErrString(err)
	}
	cols, rows := thumSize(mw.GetImageWidth(), mw.GetImageHeight(), s)
	err = mw.ResizeImage(cols, rows, imagick.FILTER_LANCZOS2_SHARP, 1)
	if err != nil {
		log.Printf("Error Resizing Thumbnail(%s) %v", p, err)
		return ErrString(err)
	}

	err = mw.SetImageCompressionQuality(90)
	if err != nil {
		log.Printf("Error Compression Thumbnail(%s) %v", p, err)
		return ErrString(err)
	}

	ensureDir(path.Dir(dest))
	err = mw.WriteImage(dest)
	if err != nil {
		log.Printf("Error Writing Thumbnail(%s) %v", p, err)
		return ErrString(err)
	}
	log.Printf("Processed Thumbnail(%s) %s", p, dest)

	return OkString(dest)
}

func GetThumbnail(p string, s uint) string {
	sn := sizeName(p, s)
	if pathExists(sn) {
		return sn
	}
	return process(p, s).FoldString(p, IdString)
}

func StopThumbnailer() {
	imagick.Terminate()
}
