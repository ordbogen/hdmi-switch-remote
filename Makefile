.PHONY: types

.SUFFIXES:

GULP = $(shell npm bin)/gulp

all: node_modules bower_components public/hello.js public/index.html public/vendor.js public/hello.css public/angular-material.min.css public/font-awesome.min.css public/fonts

run: go_hello_web all
	./go_hello_web

run\:dev: go_hello_web_bare all
	./go_hello_web_bare

go_hello_web_bare: $(shell find . -name "*.go")
	go build -o go_hello_web_bare

go_hello_web: go_hello_web_bare
	rm -f go_hello_web
	cp go_hello_web_bare go_hello_web
	rice append --exec go_hello_web

node_modules:
	npm -q update

bower_components:
	bower -q update

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

types:
	java -jar bower_components/closure-compiler/compiler.jar \
		--warning_level=VERBOSE \
		--externs support/externs/angular-1.2.js \
		--externs support/externs/externs.js \
		--angular_pass \
		tmp/hello.js

public.rice-box.go: all
	rice embed-go

hello_embedded: public.rice-box.go
	go build

setup:
	# For rice
	sudo apt-get install zip
	npm install -q -g bower gulp

hdmi-remote-darwin-amd64:
	rm -rf hdmi-remote-darwin-amd64
	rice embed-go
	go-darwin-amd64 build -o hdmi-remote-darwin-amd64
