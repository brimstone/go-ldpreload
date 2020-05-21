backdoor.so: Makefile main.go
	go build -x -v -buildmode=c-shared -o backdoor.so -ldflags '-s -w'
	ldd backdoor.so

.PHONY: watch
watch:
	find *.go Makefile | entr make
