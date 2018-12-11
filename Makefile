build: main.go info/*.go pages/*.go
	@GOARCH=amd64 GOOS=linux go build -o build/tldr -ldflags "-s -w"
	@upx build/tldr > /dev/null
	@GOARCH=amd64 GOOS=windows go build -o build/tldr.exe -ldflags "-s -w"
	@upx build/tldr.exe > /dev/null

nightly: build
	@mv build/cLC build/tldr_nightly
	@mv build/cLC.exe build/tldr_nightly.exe