FROM alpine:3.5

ENV CHARTS_SERVICE_PORT 8080
ENV CHARTS_SERVICE_HEALTH_PORT 8082

EXPOSE $CHARTS_SERVICE_PORT
EXPOSE $CHARTS_SERVICE_HEALTH_PORT

COPY charts-server /
ADD packages/ /var/lib/charts/

# https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem
RUN apk add --no-cache --virtual .build-deps \
    curl \
    && curl -sSL -o /sbin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 \
    && chmod +x /sbin/dumb-init \
    && apk del .build-deps

ENTRYPOINT ["/sbin/dumb-init", "--"]

CMD ["/charts-server"]
