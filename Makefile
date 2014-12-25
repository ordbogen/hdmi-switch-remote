.PHONY: types hdmi-switch-remote-darwin-amd64 all

.SUFFIXES:

GULP = $(shell npm bin)/gulp

NAME = hdmi-switch-remote

build: $(NAME)

all: node_modules public/hello.js public/index.html public/vendor.js public/hello.css public/angular-material.min.css public/font-awesome.min.css public/fonts
	@echo "TRAVIS_GO_VERSION: $(TRAVIS_GO_VERSION)"
	go get -d -v ./...
	go get -v github.com/GeertJohan/go.rice/rice

run\:dev: $(NAME)-bare all
	./$(NAME)-bare

$(NAME)-bare: $(shell find . -name "*.go") all
	rm -rf public.rice-box.go
	go build -o $(NAME)-bare

$(NAME): all
	rm -rf $(NAME)
	rm -rf public.rice-box.go
	go build -o $(NAME)
	rice append --exec $(NAME)

node_modules:
	npm update > /dev/null

public/hello.js: app/hello.ls
	$(GULP)

public/index.html: app/index.jade
	$(GULP)

public/vendor.js: bower.json
	$(GULP)

public/hello.css: app/hello.styl
	$(GULP)

public/angular-material.min.css:
	$(GULP)

public/font-awesome.min.css:
	$(GULP)

public/fonts:
	$(GULP)

public.rice-box.go: all
	rice embed-go

hello_embedded: public.rice-box.go
	go build

$(NAME)-darwin-amd64: all
	rm -rf $(NAME)-darwin-amd64
	rice embed-go
	GOOS=darwin GOARCH=amd64 go build -o $(NAME)-darwin-amd64
