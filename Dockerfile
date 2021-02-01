FROM golang@sha256:205d5cf61216a16da4431dd5796793c650236159fa04e055459940ddc4c6389c as build

WORKDIR /srv/grpc

COPY go.mod .
COPY proto/calculator/service.proto ./proto/
COPY *.go ./server/

ARG VERS="3.11.4"
ARG ARCH="linux-x86_64"
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${VERS}/protoc-${VERS}-${ARCH}.zip \
    --output-document=./protoc-${VERS}-${ARCH}.zip && \
    apt update && apt install -y unzip && \
    unzip -o protoc-${VERS}-${ARCH}.zip -d protoc-${VERS}-${ARCH} && \
    mv protoc-${VERS}-${ARCH}/bin/* /usr/local/bin && \
    mv protoc-${VERS}-${ARCH}/include/* /usr/local/include && \
    go get -u github.com/golang/protobuf/protoc-gen-go

RUN protoc \
    --proto_path=. \
    --go_out=plugins=grpc:. \
    ./proto/service.proto

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /go/bin/server \
    grpc-tutorial/server

FROM scratch

COPY --from=build /go/bin/server /server

ENTRYPOINT ["/server"]