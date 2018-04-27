package main

import (
	"github.com/nfnt/resize"
	"flag"
	"fmt"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"bytes"
	"image"
	"strings"
	"bufio"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	"image/color"
	"reflect"
)

var ascii = []string{
	"@", "c","*", "+", ":", "."," ",
}
var (
	path  string
	width int
)

func write(img image.Image) {
	writer := bufio.NewWriter(os.Stdout)
	gray:=image.NewGray(img.Bounds())
	for y := 0; y < gray.Rect.Max.Y; y++ {
		for x := 0; x < gray.Rect.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			yv := reflect.ValueOf(grayColor).FieldByName("Y").Uint()
			pos := int(yv * 6 / 255)
			writer.Write([]byte(ascii[pos]))
		}
		writer.WriteByte('\n')
	}
	writer.Flush()
}

func scaleImage(img image.Image, w int) image.Image {
	r := img.Bounds()
	h := (r.Max.Y * w * 10) / (r.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img
}

func getImageFromNet(url string) string {
	file, err := os.Create("in.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(file, bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	return "in.png"
}

func loadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func init() {
	flag.StringVar(&path, "p", "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", "The source path of image.")
	flag.IntVar(&width, "w", 80, "The width of output.")
	flag.Parse()
}

func main() {
	newPath := path
	if strings.Contains(path, "http") {
		newPath = getImageFromNet(path)
	}
	img, err := loadImage(newPath)
	if err != nil {
		log.Fatal(err)
	}
	scaledImg:= scaleImage(img, width)
	write(scaledImg)

}
