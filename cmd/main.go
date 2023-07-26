package main

import (
	"flag"
	"io/ioutil"
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
	if *merge {
		outPath = outPath + "/temporary"
	}

	// open file
	style, err := ioutil.ReadFile(*stylePath)
	if err != nil {
		log.Fatal(err)
	}

	card, err := ioutil.ReadFile(*cardsPath)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := generator.MakeCards(style, card)
	if err != nil {
		log.Fatal(err)
	}
	g := generator.BuildGenerator(outPath, *fontPath)
	if err = g.Generate(cards); err != nil {
		log.Fatal(err)
	}
	if *merge {
		generator.MergeCards(outPath, *outputPath)
		if err := os.RemoveAll(outPath); err != nil {
			log.Fatal(err)
		}
	}
}
