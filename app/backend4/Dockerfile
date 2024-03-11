FROM golang:1.22 as builder

WORKDIR /src

COPY . /src

RUN CGO_ENABLED=0 go build -v -o /app

FROM scratch

COPY --from=0 /app /app

EXPOSE 8080

CMD ["/app"]
