# Ignite Rest API

This is a rest api for [Ignite](https://ignite.readthedocs.io/en/stable) that is written in Go along with [GoFiber](https://docs.gofiber.io/).

## Running in Development

You'll need to have root access on your machine, because ignite requires root privileges at the moment.

```bash
sudo go run ./main.go
```

This command will start the the rest api on port `8008`. As of now, there is no authentication system in place - so make sure you have a firewall on your machine.

## Running with Docker

At the moment, there is no published docker image. However, you can easily build and run one with docker compose.
The docker image utilizes the host machine's containerd and firecracker paths.

```bash
docker-compose build
docker-compose up -d
```
