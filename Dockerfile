FROM golang:alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go build -o main .

VOLUME [ "/var/lib/firecracker", "/run/containerd/containerd.sock" ]

CMD ["/app/main"]
