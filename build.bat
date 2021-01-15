go build -ldflags "-s -w -H=windowsgui" .
upx --brute shimeji.exe
