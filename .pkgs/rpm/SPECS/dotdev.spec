Name:           dotdev
Version:        0.0.1.20250224.03
Release:        1%{?dist}
Summary:        Lightweight Web server for static HTML with built-in live reload written in Go.

License:        GPL-3.0-or-later
URL:            https://github.com/petlack/dotdev
Source:         %{name}-%{version}.tar.gz

BuildRequires:  go

%global debug_package %{nil}  # Disable automatic debuginfo package generation

%description
%{summary}

%prep
%autosetup -n %{name}-%{version}

%build
export GOOS=linux
go version
go build -a -ldflags="-linkmode=external" -o build/%{name} .

%install
install -Dm755 build/%{name} %{buildroot}%{_bindir}/%{name}

%check
go test ./...

%files
%{_bindir}/%{name}

%changelog
* Fri Jul 12 2024 Peter Laca <peter@laca.me> - %{version}-%{release}
- Initial RPM release
