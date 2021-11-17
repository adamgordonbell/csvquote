VERSION --referenced-save-only 0.5

FROM  golang:1.17-alpine3.14

RUN apk add --update --no-cache g++

WORKDIR /workdir

deps:
    COPY cmd cmd
    COPY go.mod .

build:
    FROM +deps
    RUN go build -o csvquote cmd/cvsquote/main.go
    SAVE ARTIFACT csvquote AS LOCAL csvquote

for-darwin-amd64:
   FROM +deps
   RUN GOOS=darwin GOARCH=amd64  go build -o csvquote cmd/cvsquote/main.go
   SAVE ARTIFACT csvquote AS LOCAL csvquote 

for-darwin-arm64:
   FROM +deps
   RUN GOOS=darwin GOARCH=arm64 go build -o csvquote cmd/cvsquote/main.go
   SAVE ARTIFACT csvquote AS LOCAL csvquote 

for-windows-amd64:
    FROM +deps
    RUN GOOS=windows GOARCH=amd64 go build -o csvquote.exe cmd/cvsquote/main.go
    SAVE ARTIFACT csvquote.exe AS LOCAL csvquote.exe


#    build:
#     FROM golang:1.15-alpine3.13
#     WORKDIR /example
#     ARG GOOS=linux
#     ARG GOARCH=amd64
#     ARG GOARM
#     COPY main.go ./
#     RUN go build -o main main.go
#     SAVE ARTIFACT ./main

# build-amd64:
#     FROM --platform=linux/amd64 alpine:3.13
#     COPY +build/main ./example/main
#     ENTRYPOINT ["/example/main"]
#     SAVE IMAGE --push org/myimage:latest

# build-arm-v7:
#     FROM --platform=linux/arm/v7 alpine:3.13
#     COPY \
#         --platform=linux/amd64 \
#         --build-arg GOARCH=arm \
#         --build-arg GOARM=v7 \
#         +build/main ./example/main