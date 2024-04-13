package main

import (
	"flag"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TimEngleSF/url-shortener-go/internal/db"
)

type application struct {
	Postgres *db.Postgres
	logger   *slog.Logger
}

type Link struct {
	ShortUrl    string
	RedirectUrl string
	Suffix      string
}

func main() {
	addr := flag.String("addr", "8080", "HTTP networking address")
	dbHost := flag.String("dbhost", "localhost", "PSQL database host")
	dbName := flag.String("dbname", "srscore", "PSQL database name")
	dbPort := flag.String("dbport", "5432", "PSQL database port")
	dbUser := flag.String("dbuser", "user", "PSQL database user")
	dbPass := flag.String("dbpass", "", "PSQL database password")
	dbSSLFlag := flag.Bool("dbssl", false, "PSQL database ssl mode")
	// useEnvFile := flag.Bool("useEnvFile", false, "Use a .env file")

	flag.Parse()

	/* INIT LOGGER */
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	var port string
	port = *addr

	/* INIT POSTGRES STRUCT */
	Postgres := db.Postgres{}
	dsn := db.PGDataSource{
		Host:   *dbHost,
		Port:   *dbPort,
		DbName: *dbName,
		User:   *dbUser,
		Pass:   *dbPass,
		SSL:    db.ConvSSL(*dbSSLFlag),
	}

	Postgres.Dsn = &dsn

	/* OPEN DB */
	err := Postgres.OpenDb()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		Postgres: &Postgres,
		logger:   logger,
	}
	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(port, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("server running")
}

func ConvPort(port string) string {
	hasPrefix := strings.HasPrefix(port, ":")
	if !hasPrefix {
		return ":" + port
	}
	return port
}

func isValidUrl(input string) bool {
	u, err := url.Parse(input)
	return err == nil && u.Scheme != "" && u.Host != ""
}
