BINARY_NAME = todoer

all: build run

build:
	go build -ldflags="-s -w" -o ./dist/${BINARY_NAME}

run:
	@./dist/${BINARY_NAME}

gif:
	vhs < ./media/demo.tape