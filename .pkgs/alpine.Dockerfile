FROM alpine:3.20 AS builder
RUN apk add --no-cache alpine-sdk go git sudo

RUN adduser -G abuild -g "Alpine Package Builder" -s /bin/ash -D builder && \
    echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

WORKDIR /home/builder

RUN mkdir -p /home/builder/.abuild && \
    chown builder:abuild /home/builder/.abuild

USER builder
RUN abuild-keygen -a -i -n

COPY *.go go.mod version.txt ./
COPY .pkgs/alpine/APKBUILD .

RUN tar -czf dotdev-$(cat version.txt).tar.gz \
    --transform "s,^,dotdev-$(cat version.txt)/," \
    *.go go.mod version.txt

RUN abuild checksum && abuild -r

FROM scratch
WORKDIR /pkg
COPY --from=builder /home/builder/packages/home/x86_64/*.apk .
COPY --from=builder /home/builder/.abuild/*.pub .
ENTRYPOINT ["dotdev"]
