# Excubitor-Backend

## About

## Building

### Install dependencies

To build the backend application, the following dependencies need to be installed on your system:

- Golang 1.19 / 1.20
- PAM development library
- Yarn (for building the frontend components)

You may install them the following ways, depending on your linux distribution:

#### Debian(-based)

Although this is only tested on Debian 12, it is very likely that it also works on Debian derivatives such as Ubuntu or Linux Mint.

```bash
sudo apt install golang-1.19 libpam0g-dev npm
npm install --global yarn
```

#### Fedora

```bash
sudo dnf install golang pam-devel npm
npm install --global yarn
```

### Build process

To install all go dependencies, run

```bash
make install-deps
```

Then to build the application, run

```bash
make build
```

The built executable can be found in the `bin` directory.

## Packaging

This application also provides a script to package it into .deb format. This can only be done in Debian(-based) systems.

To package the application, you may run

```bash
make package/deb
```

The package can be found in the `package/deb` folder.

## Dependencies

### Build-only dependencies

| Dependency                                                        | Creator                                                    | License                                                                        |
|-------------------------------------------------------------------|------------------------------------------------------------|--------------------------------------------------------------------------------|
| [Golang 1.19 / 1.20](https://github.com/golang/go)                | Robert Griesemer, Rob Pike, Ken Thompson and contributors  | [BSD-3-Clause](https://github.com/golang/go/blob/master/LICENSE)               |
| [PAM Development Library](https://github.com/linux-pam/linux-pam) | Dmitry V. Levin, Tomáš Mráz and contributors               | [BSD-Style or GPL](https://github.com/linux-pam/linux-pam/blob/master/COPYING) |

### Golang dependencies

| Dependency                                        | Creator                                                            | License                                                              |
|---------------------------------------------------|--------------------------------------------------------------------|----------------------------------------------------------------------|
| [Gobwas WS](https://github.com/gobwas/ws)         | Sergey Kamardin and contributors                                   | [MIT](https://github.com/gobwas/ws/blob/master/LICENSE)              |
| [jwt-go](https://github.com/golang-jwt/jwt)       | Luis Gabriel Gomez, Michael Fridman, Alistair Hey and contributors | [MIT](https://github.com/golang-jwt/jwt/blob/main/LICENSE)           |
| [uuid](https://github.com/google/uuid)            | Google Inc. and contributors                                       | [BSD-3-Clause](https://github.com/google/uuid)                       |
| [koanf](https://github.com/knadh/koanf)           | Kailash Nadh and contributors                                      | [MIT](https://github.com/knadh/koanf/blob/master/LICENSE)            |
| [go-sqlite3](https://github.com/mattn/go-sqlite3) | Yasuhiro Matsumoto and contributors                                | [MIT](https://github.com/mattn/go-sqlite3/blob/master/LICENSE)       |
| [pam](https://github.com/msteinert/pam)           | Mike Steinert and contributors                                     | [BSD-2-Clause](https://github.com/msteinert/pam/blob/master/LICENSE) |
| [cors](https://github.com/rs/cors)                | Olivier Poitrey and contributors                                   | [MIT](https://github.com/rs/cors/blob/master/LICENSE)                |
| [pflag](https://github.com/spf13/pflag)           | Steve Francia and contributors                                     | [BSD-3-Clause](https://github.com/spf13/pflag/blob/master/LICENSE)   |
| [testify](https://github.com/stretchr/testify)    | Stretchr Inc. and contributors                                     | [MIT](https://github.com/stretchr/testify/blob/master/LICENSE)       |