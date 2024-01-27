package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lapis2411/card-generator/resources"
	"image/png"
	"log"
	"os"

	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/common"
	"github.com/lapis2411/card-generator/domain"
	"github.com/lapis2411/card-generator/driver"
	"github.com/lapis2411/card-generator/driver/decoder"
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
	if err := createDirIfNotExist(outPath); err != nil {
		log.Fatal(err)
	}

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
	dc := decoder.NewCsvDecoder()
	ca := adapter.NewCardAdapter(dc)

	nCard := common.NewSize(600, 800)
	id := driver.NewImageDriver(nCard, resources.APFont)
	t := driver.NewA4Template()
	ia := adapter.NewImageAdapter(id, t)

	// ed := dRxport.NewPngExporter()
	// eia := adapter.NewExporter(outPath, ed)
	// ue := usecase.NewExport(ca, ia, eia)

	// if *merge {
	// 	err = ue.ExportCardImagesForPrint(style, card)
	// } else {
	// 	err = ue.ExportCardImages(style, card)
	// }

	ug := usecase.NewGenerate(ca, ia)
	var di []domain.Image
	if *merge {
		di, err = ug.GeneratePrintImages(style, card)
	} else {
		di, err = ug.GenerateCardImages(style, card)
	}
	fmt.Println(len(di))
	for _, d := range di {
		path := fmt.Sprintf("%s/%s.png", outPath, d.Name())
		if err := exportImage(path, d); err != nil {
			log.Fatal(err)
		}
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

func exportImage(path string, image domain.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file(%s): %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b := bufio.NewWriter(f)
	if err := png.Encode(b, image.Image()); err != nil {
		return fmt.Errorf("failed to encode card image: %w", err)
	}
	if err := b.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return nil
}
