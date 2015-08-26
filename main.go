// Copyright 2015 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

// garsue is path to original garsue icon file. Here it is thought to be converted with go-bindata.
// To install go-bindata, just 'go get github.com/jteeuwen/go-bindata'.
const garsue = "data/garsue_trans.png"

var (
	defaultColor = color.NRGBA{255, 0, 0, 255}
	colorCode    = flag.String("c", "FF0000", "color code")
)

func main() {
	// parse given
	flag.Parse()
	given, err := hexToNRGBA(*colorCode)
	if err != nil {
		log.Fatal(err)
	}

	data, err := Asset(garsue)
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewBuffer(data)

	src, err := png.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	b := src.Bounds()
	rect := image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
	buf := image.NewNRGBA(rect)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := src.At(x, y)
			if _, _, _, a := c.RGBA(); a == 0 {
				buf.Set(x, y, given)
				continue
			}
			buf.Set(x, y, c)
		}
	}

	file, err := os.Create("./garsue.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	png.Encode(file, buf)
}

func hexToNRGBA(code string) (color.NRGBA, error) {
	if len(code) != 6 {
		return defaultColor, fmt.Errorf("code must be in 6 hex sytle.")
	}
	code = strings.ToUpper(code)
	r, err := strconv.ParseInt(code[0:2], 16, 16)
	if err != nil {
		return defaultColor, err
	}
	g, err := strconv.ParseInt(code[2:4], 16, 16)
	if err != nil {
		return defaultColor, err
	}
	b, err := strconv.ParseInt(code[4:6], 16, 16)
	if err != nil {
		return defaultColor, err
	}
	c := color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	log.Println(c)
	return c, nil
}
