srcfiles := $(wildcard src/*.txt)
ifeq ($(OS),Windows_NT)
	bin_file := urem.exe
else
	bin_file := urem
endif

.PHONY : clean, fmt, gen, test

build : $(bin_file)

$(bin_file) : gen fmt $(srcfiles)
	go build
fmt :
	go fmt ./...
gen :
	go generate ./...
test :
	go test ./...
clean :
	-rm $(bin_file)
