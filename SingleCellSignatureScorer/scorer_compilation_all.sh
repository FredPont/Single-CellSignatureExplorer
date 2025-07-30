formatted_date=$(date +'%Y%m%d')
echo $formatted_date


export GOOS=linux
echo 'compilation linux'
go build -o Scorer_Linux64_$formatted_date.bin

echo 'compilation windows'
export GOOS=windows
go build -o Scorer_Win64_$formatted_date.exe

echo 'compilation Mac'
export GOOS=darwin
go build -o Scorer_Mac64_$formatted_date.bin
GOOS=darwin GOARCH=arm64 go build -o Scorer_Mac_arm64_$formatted_date.bin

