package cmdconfig

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Argv struct {
		Reverse bool

		Host                string
		Port                uint
		BasicAuth           string
		AuthorizationHeader string
		Help                bool
		Version             bool
		MarkAsDone          bool
		User                string
		Group               string
		Log                 string

		MaxOpenFiles      uint64
		NProc             uint
		PprofHostPort     string
		ChHostsString     string
		ChHosts           []string
		ChDatabase        string
		ChUser            string
		ChPassword        string
		Config            string
		Dir               string
		MaxSendSize       int64
		MaxFileSize       int64
		RotateIntervalSec int64

		LogToTables bool
	}
)

func init() {
	// actions
	flag.BoolVar(&Argv.Help, `h`, false, `show this help`)
	flag.BoolVar(&Argv.Version, `version`, false, `show version`)
	flag.BoolVar(&Argv.Reverse, `reverse`, false, `start reverse proxy server instead (ch-addr is used as clickhouse host-port)`)

	// common options
	Argv.Host = getStringOrDefault(os.Getenv("KH_HOST"), `0.0.0.0`)
	Argv.Port = getUintOrDefault(os.Getenv("KH_PORT"), 8080)
	Argv.BasicAuth = getStringOrDefault(os.Getenv("BASIC_AUTH"), "")
	Argv.User = getStringOrDefault(os.Getenv("SYSTEM_USER"), `kitten`)
	Argv.Group = getStringOrDefault(os.Getenv("SYSTEM_GROUP"), `kitten`)
	Argv.Log = getStringOrDefault(os.Getenv("LOG_FILE"), "")
	Argv.ChHostsString = getStringOrDefault(os.Getenv("CLICKHOUSE_HOSTS"), "127.0.0.1:8123")
	Argv.NProc = getUintOrDefault(os.Getenv("KH_CORES"), 0)
	Argv.PprofHostPort = getStringOrDefault(os.Getenv("PPROF_HOST"), ``)
	Argv.MaxOpenFiles = getUint64OrDefault(os.Getenv("MAX_OPEN_FILES"), 262144)
	Argv.ChDatabase = getStringOrDefault(os.Getenv("CLICKHOUSE_DATABASE"), `default`)
	Argv.ChUser = getStringOrDefault(os.Getenv("CLICKHOUSE_USER"), ``)
	Argv.ChPassword = getStringOrDefault(os.Getenv("CLICKHOUSE_PASSWORD"), ``)
	// local proxy options
	Argv.Config = getStringOrDefault(os.Getenv("CONFIG_PATH"), ``)
	Argv.Dir = getStringOrDefault(os.Getenv("PERSISTENCE_DIR"), `/tmp/kittenhouse`)
	Argv.MaxSendSize = getInt64OrDefault(os.Getenv("MAX_SEND_SIZE"), 1<<20)
	Argv.MaxFileSize = getInt64OrDefault(os.Getenv("MAX_FILE_SIZE"), 50<<20)
	Argv.RotateIntervalSec = getInt64OrDefault(os.Getenv("ROTATE_INTERVAL_SEC"), 1800)
	Argv.MarkAsDone = false
	Argv.LogToTables = getBoolOrDefault(os.Getenv("LOG_TO_TABLES"))
	flag.Parse()
	Argv.ChHosts = strings.Split(Argv.ChHostsString, `;`)
	if Argv.BasicAuth != "" {
		Argv.AuthorizationHeader = fmt.Sprintf(
			"Basic: %s",
			base64.StdEncoding.EncodeToString([]byte(Argv.BasicAuth)),
		)
	}
}

func getStringOrDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func getUintOrDefault(value string, defaultValue uint) uint {
	if value != "" {
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return defaultValue
		}
		return uint(v)
	}
	return defaultValue
}

func getUint64OrDefault(value string, defaultValue uint64) uint64 {
	if value != "" {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return defaultValue
		}
		return v
	}
	return defaultValue
}

func getInt64OrDefault(value string, defaultValue int64) int64 {
	if value != "" {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return defaultValue
		}
		return v
	}
	return defaultValue
}

func getBoolOrDefault(value string) bool {
	if strings.ToLower(value) == "true" {
		return true
	}
	return false
}
