FILENAME := thoughtsense.out
DOCKER := docker
OUTPUT := output
DEPLOY := deploy

GIT_TAG := $(shell git describe --tags --abbrev=0)
DOCKER_TAG := jdxj/thoughtsense:$(GIT_TAG)

.PHONY: clean
clean:
	@rm -rf $(OUTPUT)
	@rm -rf $(DOCKER)/$(FILENAME)
