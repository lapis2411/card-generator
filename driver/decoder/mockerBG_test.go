package decoder

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/lapis2411/card-generator/domain"
)

func TestJSONDecoder(t *testing.T) {
	jsonDecoder := NewJSONDecoder()
	type args struct {
		csv []byte
	}
	titleS := domain.NewStyle("title", image.Point{X: 100, Y: 100}, 30, color.RGBA{0, 0, 0, 255})
	descS := domain.NewStyle("description", image.Point{X: 200, Y: 200}, 15, color.RGBA{100, 150, 200, 255})
	costS := domain.NewStyle("cost", image.Point{X: 10, Y: 10}, 15, color.RGBA{255, 255, 255, 100})
	tests := []struct {
		name    string
		args    args
		want    map[string]*domain.Style
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
			want: map[string]*domain.Style{
				"title":       &titleS,
				"description": &descS,
				"cost":        &costS,
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
			got, err := csvDecoder.DecodeStyles(tt.args.csv)
			if (err != nil) != tt.wantErr {
				t.Errorf("csvDecoder.DecodeStyles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.want != nil {
					t.Errorf("csvDecoder.DecodeStyles() = %v, want %v", got, tt.want)
				}
				return
			}
			for s := range got {
				if !reflect.DeepEqual(*(got[s]), *(tt.want[s])) {
					t.Errorf("csvDecoder.DecodeStyles() = %v, want %v", *got[s], *(tt.want[s]))
				}
			}

		})
	}
}

func TestCSVCardDecoder(t *testing.T) {
	csvDecoder := NewCsvDecoder()
	titleS := domain.NewStyle("title", image.Point{X: 100, Y: 100}, 30, color.RGBA{0, 0, 0, 255})
	descS := domain.NewStyle("description", image.Point{X: 200, Y: 200}, 15, color.RGBA{100, 150, 200, 255})
	costS := domain.NewStyle("cost", image.Point{X: 10, Y: 10}, 15, color.RGBA{255, 255, 255, 100})

	useStyles := map[string]*domain.Style{
		"title":       &titleS,
		"description": &descS,
		"cost":        &costS,
	}
	type args struct {
		csv    []byte
		styles map[string]*domain.Style
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]domain.FormattedText
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
			want: map[string][]domain.FormattedText{
				"main": {
					domain.NewFormattedText("some title", &titleS),
					domain.NewFormattedText("test", &descS),
					domain.NewFormattedText("12", &costS),
				},
				"sub1": {
					domain.NewFormattedText("sub1 title", &titleS),
					domain.NewFormattedText(" desc", &descS),
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
			got, err := csvDecoder.DecodeCards(tt.args.csv, tt.args.styles)
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
