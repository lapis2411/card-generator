package main

import generator "github.com/lapis2411/card-generator"

func main() {
	// stylePath := flag.String("style", "", "Path to the style file")
	// cardsPath := flag.String("cards", "", "Path to the cards file")
	// outputPath := flag.String("output", "", "Path to the output folder")
	// fontPath := flag.String("font", "", "Path to the font file")

	// flag.Parse()

	// // open file
	// style, err := ioutil.ReadFile(*stylePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// styles, err := generator.CSVDecoder.DecodeStyles(style)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// card, err := ioutil.ReadFile(*cardsPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cards, err := generator.CSVDecoder.DecodeCards(card, styles)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// g := generator.BuildGenerator(*outputPath, *fontPath)
	// if err = g.Generate(cards); err != nil {
	// 	log.Fatal(err)
	// }

	generator.MergeCards("./temporary/")
}
