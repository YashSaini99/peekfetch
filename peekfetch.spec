%define debug_package %{nil}
Name:           peekfetch
Version:        1.0.0
Release:        1%{?dist}
Summary:        Interactive system information TUI with live updates

License:        MIT
URL:            https://github.com/YashSaini99/peekfetch
Source0:        %{url}/archive/v%{version}/%{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.21
BuildRequires:  git

%description
peekfetch is an interactive TUI (Text User Interface) for displaying detailed
system information. It features keyboard navigation, expandable sections,
live CPU and memory updates, and a modern calm color scheme. Built with Go,
Bubble Tea, and Lipgloss.

%prep
%autosetup -n %{name}-%{version}

%build
export CGO_ENABLED=0
export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"

go build -o %{name} ./cmd/peekfetch

%install
# Install binary
install -Dm755 %{name} %{buildroot}%{_bindir}/%{name}

# Install documentation
install -Dm644 README.md %{buildroot}%{_docdir}/%{name}/README.md
install -Dm644 QUICKSTART.md %{buildroot}%{_docdir}/%{name}/QUICKSTART.md
install -Dm644 FEATURES.md %{buildroot}%{_docdir}/%{name}/FEATURES.md

# Install license if you create one
# install -Dm644 LICENSE %{buildroot}%{_licensedir}/%{name}/LICENSE

%check
go test -v ./...

%files
%{_bindir}/%{name}
%{_docdir}/%{name}/README.md
%{_docdir}/%{name}/QUICKSTART.md
%{_docdir}/%{name}/FEATURES.md
# %license LICENSE
%doc README.md

%changelog
* Wed Jan 15 2025 Yash Saini <ysyashsaini3@gmail.com> - 1.0.0-1
- Initial package
