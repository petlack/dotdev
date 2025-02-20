FROM fedora:41 AS builder

RUN dnf install -y git rpm-build sudo golang

RUN useradd -m builder && \
    echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

WORKDIR /home/builder

RUN mkdir -p /home/builder/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS} && \
    chown -R builder:builder /home/builder

USER builder

COPY *.go go.mod version.txt /home/builder/rpmbuild/SOURCES/
COPY .pkgs/rpm/SPECS/*.spec /home/builder/rpmbuild/SPECS/

WORKDIR /home/builder/rpmbuild/SOURCES

RUN tar -czf dotdev-$(cat version.txt).tar.gz \
    --transform "s,^,dotdev-$(cat version.txt)/," \
    *.go go.mod version.txt

WORKDIR /home/builder/rpmbuild

RUN rpmbuild -ba SPECS/*.spec

FROM scratch
WORKDIR /pkg
COPY --from=builder /home/builder/rpmbuild/RPMS/x86_64/*.rpm .
ENTRYPOINT ["dotdev"]
