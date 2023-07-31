package main

import (
	"flag"
	"log"
	"os"

	generator "github.com/lapis2411/card-generator"
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

	cards, err := generator.MakeCards(style, card)
	if err != nil {
		log.Fatal(err)
	}
	g := generator.NormalSizeCardGenerator(*fontPath)
	cimgs, err := g.Generate(cards)
	if err != nil {
		log.Fatal(err)
	}
	if *merge {
		l := generator.NewA4Layout()
		cvs, err := l.Arrange(cimgs)
		if err != nil {
			log.Fatal(err)
		}
		if err := cvs.ExportImages(outPath); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := cimgs.ExportImages(outPath); err != nil {
			log.Fatal(err)
		}
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
