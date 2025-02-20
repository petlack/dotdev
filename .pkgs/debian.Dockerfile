FROM debian:bookworm-slim AS builder

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN apt-get update && apt-get install -y \
    build-essential \
    dpkg-dev \
    debhelper \
    devscripts \
    dh-make \
    fakeroot \
    lintian \
    software-properties-common \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

RUN echo "deb http://deb.debian.org/debian bookworm-backports main" > /etc/apt/sources.list.d/backports.list \
    && apt-get update \
    && apt-get -t bookworm-backports install -y \
    golang-go

RUN groupadd -g ${GROUP_ID} builder \
    && useradd -m -u ${USER_ID} -g builder builder

USER builder

COPY --chown=builder:builder . /home/builder/src
COPY --chown=builder:builder .pkgs/debian /home/builder/src/debian

WORKDIR /home/builder/src

RUN chmod +x debian/rules
RUN dch -v "$(cat version.txt)-1" "Automated build of $(cat version.txt)-1 with Docker"

RUN dpkg-buildpackage -us -uc -ui -b

FROM scratch
WORKDIR /pkg
COPY --from=builder /home/builder/*.deb .
ENTRYPOINT ["dotdev"]
