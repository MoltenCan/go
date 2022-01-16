FROM --platform=$BUILDPLATFORM golang:alpine as builder
COPY ./ src/
WORKDIR /go/src
RUN go mod init starter && go mod tidy

# macos
RUN GOOS=linux GOARCH=arm64 go build -o /go/bin/linux/arm64/starter .
# # rpi3+
RUN GOOS=linux GOARCH=arm GOARM=7 go build -o /go/bin/linux/arm/v7/starter .
# # rpi
RUN GOOS=linux GOARCH=arm GOARM=6 go build -o /go/bin/linux/arm/v6/starter .
# # linux
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/linux/amd64/starter .

FROM --platform=$TARGETPLATFORM alpine
ARG TARGETPLATFORM
COPY --from=builder /go/bin/${TARGETPLATFORM}/starter /bin/starter

ENTRYPOINT [ "/bin/starter" ]