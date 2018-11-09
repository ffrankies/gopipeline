GOPIPELINE := github.com/ffrankies/gopipeline
EXAMPLE := $(GOPIPELINE)/example

.PHONY: example
example:
	go install $(EXAMPLE)
	example master

.PHONY: local_example
local_example:
	go install $(EXAMPLE)
	example -config=GoPipeline.localconfig.yaml master

.PHONY: dist_example
dist_example:
	go install $(EXAMPLE)
	example -config=GoPipeline.distconfig.yaml master

.PHONY: install
install:
	go install $(GOPIPELINE)
