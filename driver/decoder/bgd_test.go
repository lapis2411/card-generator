package decoder

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/lapis2411/card-generator/domain"
)

func TestBGFDecodeCard(t *testing.T) {
	bgdDecoder := NewBGDDecoder()
	st1 := domain.NewStyle("style1", image.Point{X: 100, Y: 100}, 30, color.RGBA{0, 0, 0, 255})
	st2 := domain.NewStyle("style2", image.Point{X: 200, Y: 200}, 15, color.RGBA{100, 150, 200, 255})
	st3 := domain.NewStyle("style3", image.Point{X: 105, Y: 10}, 15, color.RGBA{255, 255, 255, 100})
	st4 := domain.NewStyle("style4", image.Point{X: 110, Y: 10}, 15, color.RGBA{255, 255, 255, 100})
	st5 := domain.NewStyle("style5", image.Point{X: 115, Y: 10}, 15, color.RGBA{255, 255, 255, 100})
	st6 := domain.NewStyle("style6", image.Point{X: 120, Y: 10}, 15, color.RGBA{255, 255, 255, 100})

	useStyles := map[string]*domain.Style{
		"1": &st1,
		"2": &st2,
		"3": &st3,
		"4": &st4,
		"5": &st5,
		"6": &st6,
	}
	type args struct {
		json   []byte
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
				json: []byte(`[
{"id":1,"name":"name1",
"data":[
    {"styleID":1,"text":"name1"},
	{"styleID":2,"text":"lu1"},
	{"styleID":4,"text":"this is text"},
	{"styleID":5,"text":"ld1"}]},
{"id":2,"name":"name2",
"data":[
	{"styleID":1,"text":"name2"},
	{"styleID":3,"text":"ru2"},
	{"styleID":4,"text":"これはテキスト"},
	{"styleID":6,"text":"rd1"}]},
{"id":3,"name":"name3",
"data":[
	{"styleID":1,"text":"name3"},
	{"styleID":2,"text":"lu3"},
	{"styleID":3,"text":"ru3"},
	{"styleID":5,"text":"ld3"},
	{"styleID":6,"text":"rd3"}]}]`),
				styles: useStyles,
			},
			want: map[string][]domain.FormattedText{
				"1": {
					domain.NewFormattedText("name1", &st1),
					domain.NewFormattedText("lu1", &st2),
					domain.NewFormattedText("this is text", &st4),
					domain.NewFormattedText("ld1", &st5),
				},
				"2": {
					domain.NewFormattedText("name2", &st1),
					domain.NewFormattedText("ru2", &st3),
					domain.NewFormattedText("これはテキスト", &st4),
					domain.NewFormattedText("rd1", &st6),
				},
				"3": {
					domain.NewFormattedText("name3", &st1),
					domain.NewFormattedText("lu3", &st2),
					domain.NewFormattedText("ru3", &st3),
					domain.NewFormattedText("ld3", &st5),
					domain.NewFormattedText("rd3", &st6),
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate card id",
			args: args{
				json: []byte(`[
{"id":0,"name":"name1",
"data":[
    {"styleID":1,"text":"name1"},
	{"styleID":2,"text":"lu1"},
	{"styleID":4,"text":"this is text"},
	{"styleID":5,"text":"ld1"}]},
{"id":0,"name":"name2",
"data":[
	{"styleID":1,"text":"name2"},
	{"styleID":3,"text":"ru2"},
	{"styleID":4,"text":"これはテキスト"},
	{"styleID":6,"text":"rd1"}]}]`),
				styles: useStyles,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "undefined style",
			args: args{
				json: []byte(`[
{"id":0,"name":"name1",
"data":[
    {"styleID":100,"text":"name1"},
	{"styleID":2,"text":"lu1"},
	{"styleID":4,"text":"this is text"},
	{"styleID":5,"text":"ld1"}]}]`),
				styles: useStyles,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bgdDecoder.DecodeCards(tt.args.json, tt.args.styles)
			if (err != nil) != tt.wantErr {
				t.Errorf("BGDDecoder.DecodeCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("len(BGDDecoder.DecodeCards()) = %v, want %v", len(got), len(tt.want))
			}
			for k := range got {
				if !reflect.DeepEqual(got[k], tt.want[k]) {
					t.Errorf("BGDDecoder.DecodeCards()[%s] = %v, want %v", k, got[k], tt.want[k])
				}
			}
		})
	}
}
