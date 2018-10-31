GOPIPELINE := github.com/ffrankies/gopipeline
EXAMPLE := $(GOPIPELINE)/example

.PHONY: local_example
local_example:
	go install $(EXAMPLE)
	example master

.PHONY: dist_example
dist_example:
	go install $(EXAMPLE)
	example -config=GoPipeline.config2.yaml master

.PHONY: install
install:
	go install $(GOPIPELINE)
