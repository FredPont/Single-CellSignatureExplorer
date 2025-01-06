date="20240613"

echo "Start compilation for windows 64bits"

#CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_ENABLED=1 go build -o filename .
#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o file.exe .
GOOS=windows CGO_ENABLED=1 go build -o Winx64_scMerger_$date.exe .


echo "Start compilation for Mac 64bits"
GOOS=darwin GOARCH=amd64 go build -o Mac_x86_64_scMerger_$date.bin .
GOOS=darwin GOARCH=arm64 go build -o Mac_arm_64_scMerger_$date.bin .

echo "Start compilation for linux 64bits"
go build -o Linux_x86_64_scMerger_$date.bin .