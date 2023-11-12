FONT_PATH=cmd/fonts/AP.ttf
run:
	go run cmd/main.go -cards cmd/sample/testCard.csv -style cmd/sample/testStyle.csv -font $(FONT_PATH) -output cmd/out -merge

run_separate:
	go run cmd/main.go -cards cmd/sample/testCard.csv -style cmd/sample/testStyle.csv -font $(FONT_PATH) -output cmd/out

rm_out:
	rm cmd/out/*

