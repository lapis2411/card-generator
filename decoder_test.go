package main

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestCSVStyleDecoder(t *testing.T) {
	type args struct {
		csv []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Styles
		wantErr bool
	}{
		{
			name: "correct style csv",
			args: args{
				csv: []byte(`name,font_size,x,y,color_r,color_g,color_b,color_a
title,30,100,100,0,0,0,255
description,15,200,200,100,150,200,255
cost,15,10,10,255,255,255,100`),
			},
			want: Styles{
				"title": &Style{
					name:     "title",
					fontsize: 30,
					position: image.Point{X: 100, Y: 100},
					rgba:     color.RGBA{0, 0, 0, 255},
				},
				"description": &Style{
					name:     "description",
					fontsize: 15,
					position: image.Point{X: 200, Y: 200},
					rgba:     color.RGBA{100, 150, 200, 255},
				},
				"cost": &Style{
					name:     "cost",
					fontsize: 15,
					position: image.Point{X: 10, Y: 10},
					rgba:     color.RGBA{255, 255, 255, 100},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate style name",
			args: args{
				csv: []byte(`name,font_size,x,y,color_r,color_g,color_b,color_a
title,30,100,100,0,0,0,255
description,15,200,200,100,150,200,255
title,35,100,100,0,0,0,255`),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CSVDecoder.DecodeStyles(tt.args.csv)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVDecoder.DecodeStyles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.want != nil {
					t.Errorf("CSVDecoder.DecodeStyles() = %v, want %v", got, tt.want)
				}
				return
			}
			for s := range got {
				if !reflect.DeepEqual(*(got[s]), *(tt.want[s])) {
					t.Errorf("CSVDecoder.DecodeStyles() = %v, want %v", *got[s], *(tt.want[s]))
				}
			}

		})
	}
}

func TestCSVCardDecoder(t *testing.T) {
	titleStyle := &Style{
		name:     "title",
		fontsize: 30,
		position: image.Point{X: 100, Y: 100},
		rgba:     color.RGBA{0, 0, 0, 255},
	}
	descStyle := &Style{
		name:     "description",
		fontsize: 15,
		position: image.Point{X: 200, Y: 200},
		rgba:     color.RGBA{100, 150, 200, 255},
	}
	costStyle := &Style{
		name:     "cost",
		fontsize: 15,
		position: image.Point{X: 10, Y: 10},
		rgba:     color.RGBA{255, 255, 255, 100},
	}
	useStyles := Styles{
		"title":       titleStyle,
		"description": descStyle,
		"cost":        costStyle,
	}

	type args struct {
		csv    []byte
		styles Styles
	}
	tests := []struct {
		name    string
		args    args
		want    Cards
		wantErr bool
	}{
		{
			name: "correct card csv",
			args: args{
				csv: []byte(`name,style,text
main,title,some title
main,description,test
sub1,title,sub1 title
sub1,description, desc
main,cost,12`),
				styles: useStyles,
			},
			want: Cards{
				"main": &Card{
					styledTexts: []StyledText{
						{
							text:  "some title",
							style: titleStyle,
						}, {
							text:  "test",
							style: descStyle,
						}, {
							text:  "12",
							style: costStyle,
						},
					},
					styles: map[string]struct{}{
						"title":       {},
						"description": {},
						"cost":        {},
					},
				},
				"sub1": &Card{
					styledTexts: []StyledText{
						{
							text:  "sub1 title",
							style: titleStyle,
						}, {
							text:  " desc",
							style: descStyle,
						},
					},
					styles: map[string]struct{}{
						"title":       {},
						"description": {},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate card name",
			args: args{
				csv: []byte(`name,style,text
main,title,some title
main,title,test
sub1,sub1 title, title
sub1,sub1 description, description`),
				styles: useStyles,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "undefined style",
			args: args{
				csv: []byte(`name,style,text
main,undefined,some title
main,title,some title`),
				styles: useStyles,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CSVDecoder.DecodeCards(tt.args.csv, tt.args.styles)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVDecoder.DecodeCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVDecoder.DecodeCards() = %v, want %v", got, tt.want)
			}
		})
	}
}
