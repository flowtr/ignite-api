FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o main .

VOLUME [ "/var/lib/firecracker", "/run/containerd/containerd.sock" ]

CMD ["/app/main"]
