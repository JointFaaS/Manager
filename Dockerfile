FROM golang:1.13 AS build

WORKDIR /go/src/app

COPY . .

RUN make manager

WORKDIR /env-addons

RUN git clone https://github.com/JointFaaS/aliyun-env-addons.git ali
RUN git clone https://github.com/JointFaaS/aws-env-addons.git aws

FROM alpine:3

WORKDIR /root/

RUN mkdir .jfManager

COPY config.yml .jfManager/

COPY --from=build /env-addons/ .jfManager
COPY --from=build /go/src/app/build/ .

CMD ["/root/manager"]