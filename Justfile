set positional-arguments := true

IMAGE_NAME := "dotdev"
BIN_NAME := "dotdev"

compile:
    go build -o build/{{BIN_NAME}} .

compile-static:
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o build/{{BIN_NAME}} .

dev *args:
    find . -name "*.go" | entr -cr sh -c 'just compile run ${@}' . "${@}"

docker-build:
    docker build -t {{IMAGE_NAME}} .

docker-run *args:
    docker run --rm \
        -e PORT=8080 \
        -e HOST=0.0.0.0 \
        -p 14774:8080 \
        -v .:/app \
        -w /app \
        --name {{IMAGE_NAME}} \
        {{IMAGE_NAME}} "${@}"

install: compile-static
    sudo install -Dm755 build/{{BIN_NAME}} /usr/local/bin/{{BIN_NAME}}

run *args='':
    ./build/{{BIN_NAME}} "${@}"

test:
    go clean -testcache && go test -v .
