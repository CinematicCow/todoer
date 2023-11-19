BINARY_NAME = todoer
OS_NAME = linux

all: build run

build:
ifeq ($(OS_NAME), windows)
	go build -ldflags="-s -w" -o ./dist/${BINARY_NAME}.exe
else 
	go build -ldflags="-s -w" -o ./dist/${BINARY_NAME}
endif

run:
	@./dist/${BINARY_NAME}

gif:
	vhs < ./media/demo.tape