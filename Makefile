GOPIPELINE := github.com/ffrankies/gopipeline
EXAMPLE := $(GOPIPELINE)/example
APRIORI := $(GOPIPELINE)/example_apriori

.PHONY: install install_example install_apriori example apriori local_example local_apriori dist_example dist_apriori

install:
	go install $(GOPIPELINE)

install_example:
	go install $(EXAMPLE)

install_apriori:
	go install $(APRIORI)

example: install_example
	example master

local_example: install_example
	example -config=GoPipeline.localconfig.yaml master

dist_example: install_example
	example -config=GoPipeline.distconfig.yaml master

apriori: install_apriori
	example_apriori master

local_apriori: install_apriori
	example_apriori -config=GoPipeline.localconfig.yaml master

dist_apriori: install_apriori
	example_apriori -config=GoPipeline.distconfig.yaml master
