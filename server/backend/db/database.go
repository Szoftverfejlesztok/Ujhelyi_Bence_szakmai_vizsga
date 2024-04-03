package db

import (
	"database/sql"

	"backend/types"
	"backend/vars"

	_ "github.com/go-sql-driver/mysql"
)

// getDB returns a database handler
func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", vars.ConnectionString)
	//db.SetMaxOpenConns(50)
	//db.SetMaxIdleConns(1000)
	return db, err
}

// AddRecord adds a log record about state2 of a lamp
// We don't give the ID due the database will create it
func AddRecord(eventLog types.Lamp) (types.Lamp, error) {
	db, err := getDB()
	if err != nil {
		return types.Lamp{}, err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO event_logs (lamp, date, state) VALUES (?, NOW(), ?)",
		eventLog.Lamp,
		eventLog.State)
	if err != nil {
		return types.Lamp{}, err
	}

	var res types.Lamp
	var lamp, date string
	var state bool

	events := db.QueryRow("SELECT lamp, date, state FROM event_logs WHERE lamp=? ORDER BY id DESC", eventLog.Lamp)
	if err = events.Scan(&lamp, &date, &state); err != nil {
		return types.Lamp{}, err
	}

	res.Lamp = lamp
	res.Date = date
	res.State = state

	return res, nil
}

// GetLastByLamp return a record with the provided lamp's name
func GetLastByLamp(recordLamp string) (types.Lamp, error) {
	db, err := getDB()
	if err != nil {
		return types.Lamp{}, err
	}
	defer db.Close()

	events, err := db.Query("SELECT lamp, date, state FROM event_logs WHERE lamp=? ORDER BY date DESC LIMIT 1", recordLamp)
	if err != nil {
		return types.Lamp{}, err
	}

	var res types.Lamp

	if events.Next() {
		var lamp, date string
		var state bool
		err = events.Scan(&lamp, &date, &state)
		if err != nil {
			return types.Lamp{}, err
		}
		res.Lamp = lamp
		res.Date = date
		res.State = state
	}

	return res, nil
}

func GetDistinctLamp() ([]types.Lamp, error) {
	db, err := getDB()
	if err != nil {
		return []types.Lamp{}, err
	}
	defer db.Close()

	lampArray, err := db.Query("SELECT lamp, state FROM ( SELECT id, lamp, date, state, ROW_NUMBER() OVER (PARTITION BY lamp ORDER BY date DESC) AS rn FROM event_logs ) AS subquery WHERE rn = 1;")
	if err != nil {
		return []types.Lamp{}, err
	}
	defer lampArray.Close()

	var res []types.Lamp
	for lampArray.Next() {
		var tmp types.Lamp
		var lamp string
		var state bool
		err = lampArray.Scan(&lamp, &state)
		if err != nil {
			return []types.Lamp{}, err
		}
		tmp.Lamp = lamp
		tmp.State = state
		res = append(res, tmp)
	}

	return res, nil
}

func GetStates() (string, error) {
	db, err := getDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	lampArray, err := db.Query("SELECT state FROM ( SELECT id, lamp, date, state, ROW_NUMBER() OVER (PARTITION BY lamp ORDER BY date DESC) AS rn FROM event_logs ) AS subquery WHERE rn = 1;")
	if err != nil {
		return "", err
	}
	defer lampArray.Close()

	var res []bool
	for lampArray.Next() {
		var tmp bool
		var state bool
		err = lampArray.Scan(&state)
		if err != nil {
			return "", err
		}
		tmp = state
		res = append(res, tmp)
	}

	var states string
	for _, state := range res {
		if state {
			states += "1"
		} else {
			states += "0"
		}
	}

	return states, nil
}

// IsLampExist return tru if lamp exits in the database
func IsLampExist(lampName string) (bool, error) {
	db, err := getDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var lamp string
	events := db.QueryRow("SELECT lamp FROM event_logs WHERE lamp=? ORDER BY id DESC LIMIT 1", lampName)
	if err = events.Scan(&lamp); err != nil {
		return false, err
	}

	if lamp != "" {
		return true, nil
	}

	return false, nil
}

func HealthCheck() error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
