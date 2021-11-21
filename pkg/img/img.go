package img

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/igor-koniukhov/img_generator/pkg/colors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strconv"
)

const (
	imgColorDefault         = "E5E5E5"
	msgColorDefault         = "AAAAAA"
	imgWDefault             = 300
	imgHDefault             = 300
	fontSizeDefault         = 0
	fontFileDefault         = "wqy-zenhei.ttc"
	dpiDefault      float64 = 72
	hintingDefault          = "none"
)

func GenerateFavicon() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	m := image.NewRGBA(image.Rect(0, 0, 16, 16))
	clr := color.RGBA{B: 0, A: 0}
	draw.Draw(m, m.Bounds(), &image.Uniform{C: clr}, image.ZP, draw.Src)
	var img image.Image = m

	if err := jpeg.Encode(buffer, img, nil); err != nil {
		return nil, err
	}
	return buffer, nil
}

func Generate(urlPart []string) (*bytes.Buffer, error) {
	var (
		err      error
		imgColor = imgColorDefault
		msgColor = msgColorDefault
		imgW     = imgWDefault
		imgH     = imgHDefault
		fontSize = fontSizeDefault
	)
	msg := ""
	// TODO pars urlPart in separate method, and fillout structure
	for i, val := range urlPart {
		switch i {
		case 1:
			if val != "" {
				imgW, err = strconv.Atoi(val)
				if err != nil {
					log.Println("Can not parse 'imgW', err: ", err)
					return nil, err
				}
			}
		case 2:
			if val != "" {
				imgH, err = strconv.Atoi(val)
				if err != nil {
					log.Println("Can not parse 'imgH', err:", err)
					return nil, err
				}
			}
		case 3:
			if val != "" {
				imgColor = val
			}
		case 4:
			if val != "" {
				msg = val
			}
		case 5:
			if val != "" {
				msgColor = val
			}
		case 6:
			fontSize, err = strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
		}
	}
	if ((imgW > 0 || imgH > 0) && msg == "") || msg == "" {
		msg = fmt.Sprintf("%d x %d", imgW, imgH)
	}
	if fontSize == 0 {
		fontSize = imgW / 9
		if imgH < imgW {
			fontSize = imgH / 5
		}
	}
	hx := colors.Hex(imgColor)
	rgb, err := hx.ToRGB()
	if err != nil {
		return nil, err
	}
	m := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	imgRgba := color.RGBA{R: rgb.Red, G: rgb.Green, B: rgb.Blue, A: 10}
	draw.Draw(m, m.Bounds(), &image.Uniform{C: imgRgba}, image.ZP, draw.Src)
	addLabel(m, imgW, imgH, msg, fontSize, colors.Hex(msgColor))
	var img image.Image = m
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}
	return buffer, nil

}
func addLabel(img *image.RGBA, imgW, imgH int, msg string, msgFontSize int, msgColor colors.Hex) {
	var (
		fontFile = fontFileDefault
		dpi      = dpiDefault
		hinting  = hintingDefault
	)
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err !=nil{
		log.Println(err)
		return
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	rgb, err := colors.Hex2RGB(msgColor)
	if err !=nil {
		log.Println(err)
		return
	}
	clr := color.Color(color.RGBA{R: rgb.Red, G: rgb.Blue, B: rgb.Green, A: 255})
	d:=&font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(clr),
		Face: truetype.NewFace(fnt, &truetype.Options{
			Size:              float64(msgFontSize),
			DPI:               dpi,
			Hinting:           h,
		}),
	}
	y := imgH/2 + msgFontSize/2 - 12
	d.Dot = fixed.Point26_6{
		X:(fixed.I(imgW)-d.MeasureString(msg))/2,
		Y:fixed.I(y),
	}
	d.DrawString(msg)
}
