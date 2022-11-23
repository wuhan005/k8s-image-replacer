FROM golang:alpine as builder
COPY . /home/
WORKDIR /home/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-image-replacer .

FROM alpine:latest as prod
RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \

WORKDIR /root/

COPY --from=0 /home/k8s-image-replacer /root/
ADD config.yaml /root/

WORKDIR /root/
ENTRYPOINT ["./k8s-image-replacer"]
EXPOSE 443
