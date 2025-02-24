[![Tests](https://github.com/petlack/dotdev/actions/workflows/tests.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/tests.yml)
[![Compile Binaries](https://github.com/petlack/dotdev/actions/workflows/compile.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/compile.yml)

[![Build Arch Linux Package](https://github.com/petlack/dotdev/actions/workflows/archlinux.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/archlinux.yml)
[![Build Alpine Package](https://github.com/petlack/dotdev/actions/workflows/alpine.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/alpine.yml)
[![Build Debian Package](https://github.com/petlack/dotdev/actions/workflows/debian.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/debian.yml)
[![Build RPM Package](https://github.com/petlack/dotdev/actions/workflows/rpm.yml/badge.svg)](https://github.com/petlack/dotdev/actions/workflows/rpm.yml)

# dotdev
🌐 A lightweight Web server for static HTML with live reload for instant updates during development.
It uses **inotify** for file watching and **WebSocket** for auto reloads.
Written in Go solely with standard library.

![Screen recording](screencast.gif)

## Usage
To run dotdev, provide the HTML file you wish to serve as the first argument.
You can optionally specify the host and port:
```bash
dotdev <file-to-watch> [--host <host>] [--port <port>]
```

### Example
Create an HTML file and serve it:
```bash
echo "<html><body>Hello World</body></html>" > index.html
dotdev index.html
```
The server will output a message similar to:
```
Serving index.html on http://localhost:4774
```
Now, whenever you update `index.html`, connected browsers will automatically reload.

## Command-Line Options
* `--host <HOST>`: Specify the host (default to `HOST` environment variable or `127.0.0.1`).
* `--port <PORT>`: Specify the port (defaults to `PORT` environment variable or `4774`).
* `--help`, `-h`: Print help information.
* `--version`, `-h`: Print version.

## Installation

### Alpine
Head to the [Releases](https://github.com/petlack/dotdev/releases) section and download the latest apk package and public key.
```bash
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev-0.0.1.20250222.01-r1.apk
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev-0.0.1.20250222.01-r1.apk.rsa.pub
cp dotdev-0.0.1.20250224.02-r1.apk.rsa.pub /etc/apk/keys/
apk add dotdev-0.0.1.20250224.02-r1.apk
```

### Arch Linux
**Install from AUR**
```bash
yay -S dotdev-git
```

**Build package from source**
```bash
git clone https://github.com/petlack/dotdev && cd dotdev || return
tar -czf archlinux/pkgbuild-src/dotdev-0.0.1.20250224.02.tar.gz \
    *.go go.mod version.txt
makepkg --dir archlinux/pkgbuild-src --noconfirm
sudo pacman -U ./archlinux/pkgbuild-src/dotdev-0.0.1.20250224.02-1-x86_64.pkg.tar.zst
```

**Install from release**

Head to the [Releases](https://github.com/petlack/dotdev/releases) section and download the latest Arch Linux package.
```bash
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev-0.0.1.20250222.01-1-x86_64.pkg.tar.zst
sudo pacman -U ./dotdev-0.0.1.20250224.02-1-x86_64.pkg.tar.zst
```

### Fedora/openSUSE
Head to the [Releases](https://github.com/petlack/dotdev/releases) section and download the latest **rpm** package.
```bash
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev_0.0.1.20250222.01.fc41.x86_64.rpm
sudo dnf install -y dotdev_0.0.1.20250224.02-1.fc41.x86_64.rpm
```

### Ubuntu/Debian
Head to the [Releases](https://github.com/petlack/dotdev/releases) section and download the latest **deb** package.
```bash
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev_0.0.1.20250222.01-1_amd64.deb
sudo dpkg -i dotdev_0.0.1.20250224.02-1_amd64.deb
```

### Other
**Build from source**

Make sure you have [Go 1.23 installed](https://go.dev/doc/install)
```bash
git clone https://github.com/petlack/dotdev && cd dotdev || return
go build -o dotdev .
install -m 755 dotdev /usr/local/bin/dotdev
```

**Install release binary**

Head to the [Releases](https://github.com/petlack/dotdev/releases) section and download the latest binary for your architecture.
Example:
```bash
wget https://github.com/petlack/dotdev/releases/download/v0.0.1.20250224.02/dotdev-linux-amd64
install -m 755 dotdev-linux-amd64 /usr/local/bin/dotdev
```
