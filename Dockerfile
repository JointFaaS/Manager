FROM golang:1.13 

WORKDIR /go/src/app

COPY . .

RUN make manager

WORKDIR /root/

RUN mkdir .jfManager

RUN git clone https://github.com/JointFaaS/aliyun-env-addons.git .jfManager/ali
RUN git clone https://github.com/JointFaaS/aws-env-addons.git .jfManager/aws

COPY config.yml .jfManager/

CMD ["/go/src/app/manager"]
