
build-arm:
	GOOS=linux GOARCH=arm64 go build -o server-arm64 .

build-win:
	GOOS=windows GOARCH=amd64 go build -o server.exe .

s:
	git push -u origin main