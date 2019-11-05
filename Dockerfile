# 基于alpine添加了bash和时区更改的基础镜像
FROM registry.cn-shanghai.aliyuncs.com/welfare-sign/alpine:0.0.1

WORKDIR /app

COPY ./bin/app /app/app
COPY ./docs /app/docs/
COPY ./config /app/config/
COPY ./public /app/public/

EXPOSE 8080

ENTRYPOINT ["./app"]