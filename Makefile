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

install_sequential:
	go install $(EXAMPLE)/exampleSequential
	go install $(APRIORI)/aprioriSequential

example: install_example
	example master

local_example: install_example
	example -config=GoPipeline.localconfig.yaml master

dist_example: install_example
	example -config=GoPipeline.distconfig.yaml master

dist_example1: install_example
	example -config=GoPipeline.distconfig1.yaml master

dist_example5: install_example
	example -config=GoPipeline.distconfig5.yaml master

dist_example10: install_example
	example -config=GoPipeline.distconfig10.yaml master

dist_example15: install_example
	example -config=GoPipeline.distconfig15.yaml master

dist_example20: install_example
	example -config=GoPipeline.distconfig20.yaml master

sequential_example: install_sequential
	exampleSequential

apriori: install_apriori
	example_apriori master

local_apriori: install_apriori
	example_apriori -config=GoPipeline.localconfig.yaml master

dist_apriori: install_apriori
	example_apriori -config=GoPipeline.distconfig.yaml master

dist_apriori1: install_apriori
	example_apriori -config=GoPipeline.distconfig1.yaml master

dist_apriori5: install_apriori
	example_apriori -config=GoPipeline.distconfig15.yaml master

dist_apriori10: install_apriori
	example_apriori -config=GoPipeline.distconfig10.yaml master

dist_apriori15: install_apriori
	example_apriori -config=GoPipeline.distconfig15.yaml master

dist_apriori20: install_apriori
	example_apriori -config=GoPipeline.distconfig20.yaml master

sequential_apriori: install_sequential
	aprioriSequential
