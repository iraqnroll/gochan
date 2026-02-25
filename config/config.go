package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// TODO: Make this configurable (CLI arg or whatever)
const (
	CONFIG_FILENAME = "config/config.toml"
)

var Config struct {
	Global struct {
		Shortname              string
		Subtitle               string
		RecentPostsNum         int
		AllowedMediaTypes      []string
		FingerprintSalt        string
		TripcodeSalt           string
		AuthenticatedTripcodes map[string]string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		SSLMode  string
	}
	Api struct {
		Enabled          bool
		RecentPostsNum   int
		SessionTokenSize int
	}
	Frontend struct {
		StaticDir string
		Enabled   bool
	}
}

func dbConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Config.Database.Host,
		Config.Database.Port,
		Config.Database.User,
		Config.Database.Password,
		Config.Database.Database,
		Config.Database.SSLMode)
}

// Open will open a SQL connection with the provided
// Postgres database. Callers of Open need to ensure
// the connection is eventually closed via the
// db.Close() method.
func OpenDBConn() (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConnString())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func InitConfig() {
	_, err := toml.DecodeFile(CONFIG_FILENAME, &Config)
	if err != nil {
		fmt.Println("Error decoding configuration file :", err)
		os.Exit(1)
	}
}

// Getter methods for config parameters :

func ApiEnabled() bool {
	return Config.Api.Enabled
}

func FrontendEnabled() bool {
	return Config.Frontend.Enabled
}

func FrontendStaticDir() string {
	return Config.Frontend.StaticDir
}

func FingerprintSalt() string {
	return Config.Global.FingerprintSalt
}

func TripcodeSalt() string {
	return Config.Global.TripcodeSalt
}

func AllowedMediaTypes() []string {
	return Config.Global.AllowedMediaTypes
}

func SessionTokenSize() int {
	return Config.Api.SessionTokenSize
}

func Shortname() string {
	return Config.Global.Shortname
}

func Subtitle() string {
	return Config.Global.Subtitle
}

func NumberOfRecentPosts() int {
	return Config.Api.RecentPostsNum
}

func AuthenticatedTripcode(password string) string {
	return Config.Global.AuthenticatedTripcodes[password]
}
