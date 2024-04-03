package vars

import (
	"fmt"
	"os"
	"strconv"
)

var (
	ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		"smarthome")
)

func GetPort() string {
	var port = os.Getenv("SH_PORT")
	if port == "" {
		port = "8088"
	}
	return port
}

func GetWSPort() string {
	var port = os.Getenv("WS_PORT")
	if port == "" {
		port = "8089"
	}
	return port
}

func GetMaxTry() int {
	var mt, err = strconv.Atoi(os.Getenv("MAX_TRY"))
	if err != nil || mt == 0 {
		mt = 25
	}
	return mt
}
