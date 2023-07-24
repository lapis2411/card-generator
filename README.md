# card-generator

cd ./cmd
go build main.go -o genetare

put your font file to "cmd/sample/fonts/"

generate -cards ./sample/testCard.csv -style ./sample/testStyle.csv -font ./fonts/hogehoge.ttf -output ./generate
If you want to output the image in a composite form
generate -cards ./sample/testCard.csv -style ./sample/testStyle.csv -font ./fonts/hogehoge.ttf -output ./generate_with_merge　-merge