# go-mirror-zig
This is a simple Go application for hosting a Zig community mirror.

It's **draft version**.
The behavior of this application will change.

If you still want to use this application, don't forget to change the `cmd/templates/index.html` file.

Contributions are welcomed!

## Features
* ACME challenge support (obtaining TLS certificates automatically).
* HTTP and TLS (HTTPS) server support.
* Automatic redirection of HTTP requests to HTTPS.
* Choosing custom upstream.

## Building
1. Clone the repository: `git clone https://github.com/SavaLione/go-mirror-zig.git`
2. Navigate to the project directory: `cd go-mirror-zig`
3. Build the project: `go build -o go-mirror-zig cmd/main.go`

## Usage
```
Usage of go-mirror-zig:
  -acme
        Obtain TLS certificates using ACME challenge.
  -acme-accept-tos
        Accept the ACME provider's Terms of Service.
  -acme-cache string
        Directory for storing obtained certificates.
  -acme-directory string
        ACME directory URL. (default "https://acme-v02.api.letsencrypt.org/directory")
  -acme-email string
        Email address for ACME registration and recovery notices.
  -acme-host string
        The hostname (domain name) for which to obtain the ACME certificate.
  -cache-dir string
        Path to the directory where downloaded content will be cached. (default "./")
  -enable-tls
        Enable the TLS (HTTPS) server. Requires -tls-cert-file and -tls-key-file.
  -http-port int
        The port for the plain HTTP listener. (default 80)
  -listen-address string
        The IP address to listen on. If empty, listens on all available interfaces.
  -redirect-to-https
        Enable automatic redirection of HTTP requests to HTTPS. Requires -enable-tls or -acme.
  -tls-cert-file string
        Path to the TLS certificate file.
  -tls-key-file string
        Path to the TLS private key file.
  -tls-port int
        The port for the secure TLS (HTTPS) listener. (default 443)
  -upstream-url string
        The URL of the upstream server to mirror/proxy. (default "https://ziglang.org")
```

An example of integration the application with Systemd:
1. Edit the service file with preferred editor (for example nano): `nano /etc/systemd/system/go-mirror-zig.service`
2. Add the configuration:
    ```
    [Unit]
    Description=A simple Go application for hosting a Zig community mirror
    After=network.target

    [Service]
    User=zig-mirror
    Group=zig-mirror
    Type=simple
    WorkingDirectory=/zfs-pool-fast/zig-mirror
    ExecStart=/go-mirror-zig -cache-dir=/zfs-pool-fast/zig-mirror
    Restart=always

    [Install]
    WantedBy=multi-user.target
    ```
    * Don't forget to change the `WorkingDirectory` and `-cache-dir` flag
3. Start the service: `systemctl start go-mirror-zig`
4. Enable (make the script start after boot) it: `systemctl start go-mirror-zig`

## Licenses and Acknowledgements
This project is licensed under [The GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
See the [LICENSE](LICENSE) file for the full license text.

Copyright (C) 2025 Savelii Pototskii (savalione.com)
