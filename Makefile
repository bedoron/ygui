
.PHONY:=build
build:
	go build -o ygui

.PHONY:=run
run: build
	cat demo.yaml | ./ygui

.PHONY:=run-long
run-long:
	cat long_demo.yaml | ./ygui