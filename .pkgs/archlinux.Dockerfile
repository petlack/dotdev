FROM archlinux:base-devel AS builder

RUN pacman -Syu --noconfirm
RUN pacman -S --noconfirm git

RUN useradd -m builduser && \
    passwd -d builduser && \
    printf 'builduser ALL=(ALL) ALL\n' | tee -a /etc/sudoers

USER builduser

COPY --chown=builduser:builduser .pkgs/archlinux/pkgbuild-src/PKGBUILD /home/builduser/
COPY --chown=builduser:builduser . /home/builduser

WORKDIR /home/builduser

RUN tar -czf dotdev-$(cat version.txt).tar.gz \
    *.go go.mod version.txt

RUN mkdir package
RUN makepkg -s --noconfirm
RUN mv *.pkg.tar.zst /home/builduser/package/

FROM scratch
WORKDIR /pkg
COPY --from=builder /home/builduser/package .
ENTRYPOINT ["dotdev"]
