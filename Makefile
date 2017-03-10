TARGET?=restfulmq
OS?=linux
ARCH?=amd64

XTARGET=${TARGET}-${OS}-${ARCH}

SRCS= $(shell find . -name '*.go')

build: ${TARGET}

xbuild: ${XTARGET}

${TARGET}: ${SRCS}
	go build -o $@

${XTARGET}: ${SRCS}
	${MAKE} build GOOS=${OS} GOARCH=${ARCH} TARGET=$@

clean:
	go clean
	rm -f *~ ${TARGET}-*-*
