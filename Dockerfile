FROM golang:1.15-buster as build

WORKDIR /srv/grpc

COPY go.mod .
COPY proto/ ./proto
COPY /server/*.go ./server/

ARG VERS="3.11.4"
ARG ARCH="linux-x86_64"
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${VERS}/protoc-${VERS}-${ARCH}.zip \
    --output-document=./protoc-${VERS}-${ARCH}.zip && \
    apt update && apt install -y unzip && \
    unzip -o protoc-${VERS}-${ARCH}.zip -d protoc-${VERS}-${ARCH} && \
    mv protoc-${VERS}-${ARCH}/bin/* /usr/local/bin && \
    mv protoc-${VERS}-${ARCH}/include/* /usr/local/include && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN protoc \
    -I ./proto \
    --go_out ./proto \
    --go_opt paths=source_relative \
    --go-grpc_out ./proto \
    --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./proto \
    --grpc-gateway_opt paths=source_relative \
    ./proto/calculator/service.proto

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /go/bin/server \
    ./server

FROM scratch

COPY --from=build /go/bin/server /server

ENTRYPOINT ["/server"]