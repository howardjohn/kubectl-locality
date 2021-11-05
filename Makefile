.ONESHELL:
SHELL := /bin/bash
MODULE = github.com/howardjohn/kubectl-locality

all: format
.PHONY: check-git
check-git:
	@
	if [[ -n $$(git status --porcelain) ]]; then
		echo "Error: git is not clean"
		git status
		git diff
		exit 1
	fi

.PHONY: gen-check
gen-check: check-git format

.PHONY: format
format:
	@go mod tidy
	@goimports -l -w -local $(MODULE) .

.PHONY: install
install:
	@go install
