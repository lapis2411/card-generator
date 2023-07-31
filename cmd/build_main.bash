set -e
cd "$(dirname "$0")"
go build -o cmd.exe main.go
chmod a+x cmd.exe