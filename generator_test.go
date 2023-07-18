package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTemplateCard(t *testing.T) {
	type args struct {
		border color.RGBA
		size   Size
		width  int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				border: color.RGBA{R: 0, G: 103, B: 192, A: 255},
				size:   Size{Width: 200, Height: 300},
				width:  15,
			},
			want:    "./generatortest/border.png",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := templateCard(tt.args.border, tt.args.size, tt.args.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("templateCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want, err := decodeImage(tt.want)
			if err != nil {
				t.Errorf("want image decode error = %v", err)
				return
			}
			if diff := cmp.Diff(got, want); diff != "" {
				exportImage("./generatortest/border_got.png", got)
				t.Error("templateCard() mismatch and got png exported at ./generatortest/border_got.png")
			}
		})
	}
}

var style1 = &Style{
	name:     "style",
	fontsize: 60,
	position: image.Point{X: 100, Y: 100},
	rgba:     color.RGBA{255, 0, 0, 255},
}
var style2 = &Style{
	name:     "style2",
	fontsize: 30,
	position: image.Point{X: 450, Y: 700},
	rgba:     color.RGBA{0, 150, 200, 255},
}

func TestGenerate(t *testing.T) {

	type args struct {
		cards Cards
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cards: Cards{
					"test_card1": &Card{
						styledTexts: []StyledText{
							{
								text:  "style1",
								style: style1,
							}, {
								text:  "style2",
								style: style2,
							},
						},
						// styles: map[string]struct{}{}, for duplicate check, dont need in generation
					},
					"test_card2": &Card{
						styledTexts: []StyledText{
							{
								text:  "style2",
								style: style2,
							},
						},
						// styles: map[string]struct{}{}, for duplicate check, dont need in generation
					},
				},
			},
			want: map[string]string{
				"test_card1": "./generatortest/test_card1.png",
				"test_card2": "./generatortest/test_card2.png",
			},
			wantErr: false,
		},
	}

	outputdir := "../temporary/"

	files, _ := ioutil.ReadDir(outputdir)
	if files != nil {
		for _, file := range files {
			t.Errorf("temporary directory is not empty: %s", file.Name())
		}
		return
	}
	defer func() {
		files, err := ioutil.ReadDir(outputdir)
		if err != nil {
			t.Error(err)
		}

		// すべてのファイルを削除します。
		for _, file := range files {
			err = os.Remove(outputdir + "/" + file.Name())
			if err != nil {
				t.Errorf("can't delete generated files for test: %v", err)
			}
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Generate(tt.args.cards)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			im := make(map[string]*image.RGBA, len(tt.want))
			for k, p := range tt.want {
				im[k], err = decodeImage(p)
				if err != nil {
					t.Errorf("want image decode error = %v", err)
					return
				}
			}
			for k, v := range im {
				got, err := decodeImage("../temporary/" + k + ".png")
				if err != nil {
					t.Errorf("got image decode error = %v", err)
					return
				}
				if diff := cmp.Diff(got, v); diff != "" {
					exportImage("./generatortest/"+k+"_got.png", got)
					t.Error("Generate() mismatch and got png exported at ./generatortest/border_got.png")
				}
			}
		})
	}
}

func decodeImage(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file(%s): %v", path, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return nil, fmt.Errorf("failed to convert image to RGBA")
	}
	return rgba, nil
}

// generate card image for test
func wantGeneration() {
	img, _ := templateCard(color.RGBA{
		R: 0,
		G: 103,
		B: 192,
		A: 255,
	}, Size{
		Width:  200,
		Height: 300,
	}, 15)
	exportImage("./generatortest/border.png", img)

	cards := Cards{
		"test_card1": &Card{
			styledTexts: []StyledText{
				{
					text:  "style1",
					style: style1,
				}, {
					text:  "style2",
					style: style2,
				},
			},
			// styles: map[string]struct{}{}, for duplicate check, dont need in generation
		},
		"test_card2": &Card{
			styledTexts: []StyledText{
				{
					text:  "style2",
					style: style2,
				},
			},
			// styles: map[string]struct{}{}, for duplicate check, dont need in generation
		},
	}

	Generate(cards)

}
