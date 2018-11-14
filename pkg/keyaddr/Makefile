ANDROID_HOME ?= $(HOME)/Library/Android/sdk/

.PHONY: default all test build clean

default: all

all: test build ios android keyaddr.zip

test:
	go test

build: ios android

ios: Keyaddr.framework

sources: address.go key.go key_conv.go signature.go words.go

Keyaddr.framework: sources
	gomobile bind -target ios -v

keyaddr.zip: Keyaddr.framework
	rm -f keyaddr.zip
	zip -ur keyaddr Keyaddr.framework/

android: keyaddr.aar

keyaddr.aar: sources
	gomobile bind -target android -v

clean:
	rm -rf Keyaddr.framework
	rm -f keyaddr.aar
	rm -f keyaddr-sources.jar
