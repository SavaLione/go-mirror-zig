# Go Mirror Zig
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Go Report Card](https://goreportcard.com/badge/github.com/savalione/go-mirror-zig)](https://goreportcard.com/report/github.com/savalione/go-mirror-zig) ![Go Version](https://img.shields.io/badge/go-1.26+-blue.svg) ![GitHub last commit](https://img.shields.io/github/last-commit/savalione/go-mirror-zig) ![GitHub issues](https://img.shields.io/github/issues/savalione/go-mirror-zig)

A self-hostable solution written in Go for creating a community mirror for the Zig programming language ([Zig Programming Language](https://ziglang.org/)).
This application is designed for communities, companies, or individuals looking to provide faster local access to Zig toolchains, reducing latency and bandwidth usage on the official servers.

It is lightweight and distributed as a single binary.

## Features
* Artifact caching: Local storage of upstream content uses the official Zig directory structure.
* Integrated security: ACME (Let's Encrypt) support and automatic HTTP to HTTPS redirection.
* Standalone binary: Single, dependency-free binary with no external runtime requirements.
* CLI configuration: Parameter control via commandline flags for ports, paths, and upstream settings.
* Automated maintenance: Periodic removal of stale development builds synchronized with the upstream index.
* Customizable index page: Serve a custom landing page or static directory at the root, with option to completely disable the default index.

## Getting started
### From a precompiled Binary
This is the recommended method for most users.

Navigate to the [latest release page](https://github.com/savalione/go-mirror-zig/releases/latest).

Download the archive for your operating system and architecture (e.g., `go-mirror-zig-v1.0.0-linux-amd64.tar.gz`).
Extract the archive.
```sh
tar -xvzf go-mirror-zig-v1.0.0-linux-amd64.tar.gz
```

Run the application. You can verify it's working by checking the version or help output.
```sh
./go-mirror-zig --version
```

### From source
Ensure you have a recent version of Go installed.

Clone the repository:
```sh
git clone https://github.com/savalione/go-mirror-zig.git
```

Navigate to the project directory:
```sh
cd go-mirror-zig
```

Build the project:
```sh
go build -o go-mirror-zig ./cmd/main.go
```

## Examples
### Running behind a reverse proxy (e.g., nginx)
If you already have nginx or Apache on ports 80/443, run the mirror on a high port (like 8080) and let the proxy handle TLS.
See the [deployment](#deployment) section for the nginx configuration.
```sh
./go-mirror-zig -http-port 8080 -cache-dir="/zig-mirror"
```

### Standalone with automatic TLS (ACME)
For example you have a server without caching proxy.
In that case you can set up the caching mirror with automatic ACME support.

First of all you need to open ports `80` (HTTP) and `443` (HTTPS).

After that you need to decide:
* Where to store mirror cache (for example: `/zig-mirror`)
* The location for storing the obtained TLS certificates, which must be a secure directory (e.g., `/secure-location`).

Then you can run the application with the following flags:
```sh
go-mirror-zig -acme -acme-accept-tos -acme-cache /secure-location -acme-email someone@example.com -acme-host example.com -cache-dir /zig-mirror -redirect-to-https
```

## Configuration
The application is configured using command-line flags.
Run `./go-mirror-zig -help` to see all available options.

|Flag                    |Description                                                                                   |Default value        |
|:-----------------------|:---------------------------------------------------------------------------------------------|:--------------------|
|`-acme`                 |Obtain TLS certificates using the ACME challenge.                                             |                     |
|`-acme-accept-tos`      |Accept the ACME provider's Terms of Service.                                                  |                     |
|`-acme-cache string`    |Directory for storing obtained certificates.                                                  |                     |
|`-acme-directory string`|ACME directory URL.                                                                           |`https://acme-v02.api.letsencrypt.org/directory`|
|`-acme-email string`    |Email address for ACME registration and recovery notices.                                     |                     |
|`-acme-host string`     |The hostname (domain name) for which to obtain the ACME certificate.                          |                     |
|`-cache-dir string`     |Path to the directory where downloaded content will be cached.                                |`./`                 |
|`-enable-tls`           |Enable the TLS (HTTPS) server. Requires `-tls-cert-file` and `-tls-key-file`.                 |                     |
|`-http-port int`        |The port for the plain HTTP listener.                                                         |`80`                 |
|`-listen-address string`|The IP address to listen on. If empty, listens on all available interfaces.                   |                     |
|`-redirect-to-https`    |Enable automatic redirection of HTTP requests to HTTPS. Requires `-enable-tls` or `-acme`.    |                     |
|`-tls-cert-file string` |Path to the TLS certificate file.                                                             |                     |
|`-tls-key-file string`  |Path to the TLS private key file.                                                             |                     |
|`-tls-port int`         |The port for the secure TLS (HTTPS) listener.                                                 |`443`                |
|`-upstream-url string`  |The URL of the upstream server to mirror/proxy.                                               |`https://ziglang.org`|
|`-version`              |Print version information and exit.                                                           |                     |
|`-show-possible-size`   |Print estimation stats of all cacheable upstream artifacts (size, release counts) and exit.   |                     |
|`-show-index-page bool` |Whether to serve a custom index page at the root (/). Set to false to disable.                |`true`               |
|`-index-page string`    |Path to a directory containing static files for the index. If empty, the default page is used.|built-in index page  |
|`-clear-builds-interval`|Interval in seconds to clean up cached dev builds. Set to 0 to disable.                       |`7200`               |

## Deployment
### Using systemd and nginx as a reverse proxy
For example, you have the following setup:
* A headless (CLI access only) Ubuntu server.
* nginx as a caching (and ACME challenge) proxy.

In this setup the ports `80` (HTTP) and `443` (HTTPS) are already occupied by the caching proxy.

First, you need to decide where the cache will be stored.
For this example, we'll assume you want to store it in the `/zig-mirror` directory, which you have already created.

Create a service file:
```sh
sudo nano /etc/systemd/system/go-mirror-zig.service
```

Add the following configuration, adjusting paths and flags as needed:
```ini
[Unit]
Description=Go Mirror Zig Service
After=network.target

[Service]
User=zig-mirror
Group=zig-mirror
Type=simple
WorkingDirectory=/opt/zig-mirror
ExecStart=/go-mirror-zig -http-port=8888 -cache-dir=/zig-mirror
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Then you need to create a nginx configuration for your mirror:
```
# zig.example.com
server {
    listen 80;
    server_name zig.example.com;

    location ~/.well-known {
        allow all;
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Protocol $scheme;
        proxy_set_header X-Forwarded-Host $http_host;
    }
}
```

Enable and start go-mirror-zig, reload nginx:
```sh
sudo systemctl daemon-reload
sudo systemctl enable go-mirror-zig.service
sudo systemctl start go-mirror-zig.service
sudo systemctl reload nginx
```

In the end you may issue a TLS certificate using Certbot/acme.sh or manually set up certificates.

### Using systemd and ACME challenge
Here is an example service file for running the application with systemd and ACME challenge.

Create the service file:
```sh
sudo nano /etc/systemd/system/go-mirror-zig.service
```

Add the following configuration, adjusting paths and flags as needed:
```ini
[Unit]
Description=Go Mirror Zig Service
After=network.target

[Service]
User=zig-mirror
Group=zig-mirror
Type=simple
WorkingDirectory=/opt/zig-mirror
ExecStart=/go-mirror-zig -cache-dir=/zig-mirror -acme -acme-accept-tos -acme-host=zig.example.com -acme-email=someone@example.com -acme-cache=/var/lib/zig-mirror/acme -redirect-to-https
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Reload, enable, and start the service:
```sh
sudo systemctl daemon-reload
sudo systemctl enable go-mirror-zig.service
sudo systemctl start go-mirror-zig.service
```

### Standalone without TLS and systemd
You need to decide where the cache will be stored.
In this example, the cache will be stored in `/zig-mirror`.

Run the application with the following flags:
```sh
go-mirror-zig -cache-dir="/zig-mirror"
```

## Contributing
Contributions are welcome!
We value a healthy and collaborative community.

Please read our [Contributing Guidelines](CONTRIBUTING.md) to get started.
All participants are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Licenses and acknowledgements
This project is licensed under [The GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
See the [LICENSE](LICENSE) file for the full license text.

Copyright (C) 2025-2026 Savelii Pototskii (savalione.com)

### Third-party libraries and assets
This project incorporates code from several third-party libraries and assets.
We are grateful to their developers and maintainers.
* [new.css](https://github.com/xz/new.css) - MIT License
* [Official Zig Project Logo](https://codeberg.org/ziglang/logo) - CC BY-SA 4.0
* [The Inter font family](https://github.com/rsms/inter) - OFL-1.1 License
