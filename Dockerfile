FROM golang:alpine as builder

COPY . /home/

WORKDIR /home/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-image-replacer .

FROM alpine:latest as prod

WORKDIR /root/

COPY --from=0 /home/k8s-image-replacer /root/

ENTRYPOINT ["./k8s-image-replacer"]

EXPOSE 443
