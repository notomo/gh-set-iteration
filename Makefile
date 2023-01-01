GH_NAME:=set-iteration

build:
	go build -o gh-${GH_NAME} main.go

install: build
	gh extension remove ${GH_NAME} || echo
	gh extension install .

start: install
	gh ${GH_NAME} -project-url=https://github.com/users/notomo/projects/1 -content-url=https://github.com/notomo/todo/issues/702 -field=Iteration -log=/dev/stdout -offset-days=-7

test:
	go test -v ./...
