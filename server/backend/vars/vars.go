package vars

import (
	"fmt"
	"os"
	"strconv"
)

var (
	ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		getDBUser(),
		getDBPassword(),
		getDBHost(),
		getDBPort(),
		"smarthome")
)

func getDBUser() string {
	var user = os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	return user
}

func getDBPassword() string {
	var pass = os.Getenv("DB_PASS")
	if pass == "" {
		pass = "supersecret"
	}
	return pass
}

func getDBHost() string {
	var host = os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	return host
}

func getDBPort() string {
	var port = os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	return port
}

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

func GetKey() string {
	var key = os.Getenv("XOR_KEY")
	if key == "" {
		key = "0100010001000100"
	}
	return key
}

func GetMaxTry() int {
	var mt, err = strconv.Atoi(os.Getenv("MAX_TRY"))
	if err != nil || mt == 0 {
		mt = 10
	}
	return mt
}
