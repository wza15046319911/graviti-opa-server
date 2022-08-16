SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with verison infos
versionDir = "golang-gin-demo/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitBranch = $(shell git rev-parse --abbrev-ref HEAD)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} \
            -X ${versionDir}.buildDate=${buildDate} \
		    -X ${versionDir}.gitCommit=${gitCommit} \
			-X ${versionDir}.gitBranch=${gitBranch} \
			-X ${versionDir}.gitTreeState=${gitTreeState}"


all:
	go build -v  -ldflags ${ldflags} -i -o ./build/golang-gin-demo
clean:
	rm -f ./build/golang-gin-demo
	find . -name "[._]*.s[a-w][a-z]" | xargs rm -f
gotool:
	gofmt -w .
	go vet . | grep -v vendor;true

help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"

.PHONY: clean gotool help
