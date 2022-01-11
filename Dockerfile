FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/ProblemMicro
COPY . .
ENV GOPROXY https://proxy.golang.org,direct
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/problemservice

FROM scratch
COPY --from=builder /go/src/ProblemMicro/build/problemservice /usr/bin/problemservice
COPY --from=builder /go/src/ProblemMicro/problserv.* /home/certificates/
ENTRYPOINT [ "/usr/bin/problemservice" ]