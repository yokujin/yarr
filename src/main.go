package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkanaev/yarr/src/platform"
	"github.com/nkanaev/yarr/src/server"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/utils/cli"
	"github.com/nkanaev/yarr/src/utils/logging"
	"github.com/rs/zerolog/log"
)

var Version string = "0.0"
var GitHash string = "unknown"

func parseAuthfile(authfile io.Reader) (username, password string, err error) {
	scanner := bufio.NewScanner(authfile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("wrong syntax (expected `username:password`)")
		}
		username = parts[0]
		password = parts[1]
		break
	}
	return username, password, nil
}

func main() {
	err := platform.FixConsoleIfNeeded()
	if err != nil {
		log.Fatal().Err(err).Msg("fix console")
	}

	cli.Parse()

	if cli.Ver {
		fmt.Printf("v%s (%s)\n", Version, GitHash)
		return
	}

	logging.Setup()

	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Msg("get config dir")
	}

	if cli.DB == "" {
		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatal().Err(err).Msg("create app config dir")
		}
		cli.DB = filepath.Join(storagePath, "storage.db")
	}

	log.Printf("using db file %s", cli.DB)

	var username, password string
	if cli.AuthFile != "" {
		f, err := os.Open(cli.AuthFile)
		if err != nil {
			log.Fatal().Err(err).Msg("open auth file")
		}
		defer f.Close()
		username, password, err = parseAuthfile(f)
		if err != nil {
			log.Fatal().Err(err).Msg("parse auth file")
		}
	} else if cli.Auth != "" {
		username, password, err = parseAuthfile(strings.NewReader(cli.Auth))
		if err != nil {
			log.Fatal().Err(err).Msg("parse auth literal")
		}
	}

	if (cli.CertFile != "" || cli.KeyFile != "") && (cli.CertFile == "" || cli.KeyFile == "") {
		log.Fatal().Msg("Both cert & key files are required")
	}

	store, err := storage.New(cli.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("initialise database")
	}

	srv := server.NewServer(store, cli.Addr)

	if cli.BasePath != "" {
		srv.BasePath = "/" + strings.Trim(cli.BasePath, "/")
	}

	if cli.CertFile != "" && cli.KeyFile != "" {
		srv.CertFile = cli.CertFile
		srv.KeyFile = cli.KeyFile
	}

	if username != "" && password != "" {
		srv.Username = username
		srv.Password = password
	}

	log.Printf("starting server at %s", srv.GetAddr())
	if cli.Open {
		err := platform.Open(srv.GetAddr())
		if err != nil {
			log.Fatal().Err(err).Msg("open platform")
		}
	}
	platform.Start(srv)
}
