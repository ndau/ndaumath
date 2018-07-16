ANDROID_HOME ?= $(HOME)/Library/Android/sdk/

.PHONY: default all test build clean

default: all

all: clean test build ios android

test:
	go test

build: ios android

ios: Keyaddr.framework

Keyaddr.framework: api.go
	gomobile bind -target ios -v

keyaddr.zip: Keyaddr.framework
	rm keyaddr.zip
	zip -ur keyaddr Keyaddr.framework/

android: keyaddr.aar

keyaddr.aar: api.go
	gomobile bind -target android -v

clean:
	rm -rf Keyaddr.framework
	rm keyaddr.aar
	rm keyaddr-sources.jar
