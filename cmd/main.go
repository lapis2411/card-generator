package main

import (
	"flag"
	"log"
	"os"

	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/common"
	"github.com/lapis2411/card-generator/driver"
	"github.com/lapis2411/card-generator/driver/template"
	"github.com/lapis2411/card-generator/usecase"
)

func main() {
	stylePath := flag.String("style", "", "Path to the style csv file")
	cardsPath := flag.String("cards", "", "Path to the cards csv file")
	outputPath := flag.String("output", "", "Path to the output folder")
	fontPath := flag.String("font", "", "Path to the font file")
	merge := flag.Bool("merge", false, "Merge cards")

	flag.Parse()
	if *stylePath == "" {
		log.Fatal("style path is not specified")
	}
	if *cardsPath == "" {
		log.Fatal("cards path is not specified")
	}
	if *fontPath == "" {
		log.Fatal("font path is not specified")
	}
	outPath := *outputPath
	if outPath == "" {
		outPath = "./out"
	}
	createDirIfNotExist(outPath)

	// open file
	style, err := os.ReadFile(*stylePath)
	if err != nil {
		log.Fatal(err)
	}

	card, err := os.ReadFile(*cardsPath)
	if err != nil {
		log.Fatal(err)
	}

	//
	dc := driver.NewCsvDecoder()
	ca := adapter.NewCardAdapter(dc)

	nCard := common.NewSize(600, 800)
	id := driver.NewImageDriver(nCard, *fontPath)
	t := template.NewA4Template()
	ia := adapter.NewImageAdapter(id, t)

	ue := usecase.NewExport(ca, ia)

	if *merge {
		err = ue.ExportCardImagesForPrint(style, card)
	} else {
		err = ue.ExportCardImages(style, card)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}
