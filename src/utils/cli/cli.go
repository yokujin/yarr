package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var Addr, DB, AuthFile, Auth, CertFile, KeyFile, BasePath, LogFile string
var Ver, Open bool

var OptList = make([]string, 0)

func opt(envVar, defaultValue string) string {
	OptList = append(OptList, envVar)
	value := os.Getenv(envVar)
	if value != "" {
		return value
	}
	return defaultValue
}

func init() {
	flag.CommandLine.SetOutput(os.Stdout)

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(out, "\nThe environmental variables, if present, will be used to provide\nthe default values for the params above:")
		fmt.Fprintln(out, " ", strings.Join(OptList, ", "))
	}

	flag.StringVar(&Addr, "addr", opt("YARR_ADDR", "127.0.0.1:7070"), "address to run server on")
	flag.StringVar(&BasePath, "base", opt("YARR_BASE", ""), "base path of the service url")
	flag.StringVar(&AuthFile, "auth-file", opt("YARR_AUTHFILE", ""), "`path` to a file containing username:password. Takes precedence over --auth (or YARR_AUTH)")
	flag.StringVar(&Auth, "auth", opt("YARR_AUTH", ""), "string with username and password in the format `username:password`")
	flag.StringVar(&CertFile, "cert-file", opt("YARR_CERTFILE", ""), "`path` to cert file for https")
	flag.StringVar(&KeyFile, "key-file", opt("YARR_KEYFILE", ""), "`path` to key file for https")
	flag.StringVar(&DB, "db", opt("YARR_DB", ""), "storage file `path`")
	flag.StringVar(&LogFile, "log-file", opt("YARR_LOGFILE", ""), "`path` to log file to use instead of stdout")
	flag.BoolVar(&Ver, "version", false, "print application version")
	flag.BoolVar(&Open, "open", false, "open the server in browser")

}

func Parse() {
	flag.Parse()
}
