.PHONY: build clean deploy format gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/sticker sticker/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

format:
	go fmt ./sticker

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
