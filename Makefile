.PHONY: build render all dev

MARKDOWN=$(wildcard blog/*.md)
TARGETS= \
	$(patsubst blog/%,public/blog/%,$(wildcard blog/*.html)) \
	$(patsubst blog/%,public/blog/%,$(MARKDOWN:.md=.html)) \
	$(patsubst html/pages/%,public/%,$(wildcard html/pages/*.html))

all: build render

build: renderer

dev:
	bash sync-dev.sh

render: $(TARGETS)

renderer: *.go
	go build -o renderer .

public/blog/index.html: blog/index.html $(MARKDOWN) renderer
	mkdir -p public/blog html/pages/blog
	./renderer -src $< > html/pages/$<
	./renderer -src html/pages/$< > $@

public/blog/%.html: blog/%.md renderer
	mkdir -p public/blog html/pages
	./renderer -src $< > html/pages/$(<:.md=.html)
	./renderer -src html/pages/$(<:.md=.html) > $@

public/%.html: html/pages/%.html renderer
	mkdir -p public
	./renderer -src $< > $@
