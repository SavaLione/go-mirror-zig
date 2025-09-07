package flags

import (
	"flag"
	"strconv"
)

type Flags struct {
	CacheDir    string
	UpstreamURL string
	Port        int
	IP          string
}

func NewFlags() (f Flags) {
	flag.StringVar(&f.CacheDir, "cache-dir", "./", "Directory to store cache")
	flag.StringVar(&f.UpstreamURL, "upstream-url", "https://ziglang.org", "Zig upstream mirror")
	flag.IntVar(&f.Port, "port", 8080, "Port to listen on")
	flag.StringVar(&f.IP, "ip", "", "IP to listen on")

	flag.Parse()

	return f
}

func (f Flags) Address() string {
	return f.IP + ":" + strconv.Itoa(f.Port)
}
