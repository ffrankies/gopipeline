GOPIPELINE := github.com/ffrankies/gopipeline
EXAMPLE := $(GOPIPELINE)/example

.PHONY: example
example:
	go install $(EXAMPLE)
	example master

.PHONY: install
install:
	go install $(GOPIPELINE)
