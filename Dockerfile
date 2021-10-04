FROM golang:1.15 as build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build


FROM ubuntu:latest as proxy-api-server
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=build /project/bin/api /
CMD ["./api"]


FROM ubuntu:latest as proxy-server
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=build /project/bin/proxy /
CMD ["./proxy"]
