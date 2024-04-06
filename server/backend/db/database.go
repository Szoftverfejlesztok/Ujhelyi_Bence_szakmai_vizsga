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

// AddRecord adds a log record about state2 of a device
// We don't give the ID due the database will create it
func AddRecord(eventLog types.Device) (types.Device, error) {
	db, err := getDB()
	if err != nil {
		return types.Device{}, err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO event_logs (device, date, state) VALUES (?, NOW(), ?)",
		eventLog.Device,
		eventLog.State)
	if err != nil {
		return types.Device{}, err
	}

	var res types.Device
	var device, date string
	var state bool

	events := db.QueryRow("SELECT device, date, state FROM event_logs WHERE device=? ORDER BY id DESC", eventLog.Device)
	if err = events.Scan(&device, &date, &state); err != nil {
		return types.Device{}, err
	}

	res.Device = device
	res.Date = date
	res.State = state

	return res, nil
}

// GetLastByDevice return a record with the provided device's name
func GetLastByDevice(recordDevice string) (types.Device, error) {
	db, err := getDB()
	if err != nil {
		return types.Device{}, err
	}
	defer db.Close()

	events, err := db.Query("SELECT device, date, state FROM event_logs WHERE device=? ORDER BY date DESC LIMIT 1", recordDevice)
	if err != nil {
		return types.Device{}, err
	}

	var res types.Device

	if events.Next() {
		var device, date string
		var state bool
		err = events.Scan(&device, &date, &state)
		if err != nil {
			return types.Device{}, err
		}
		res.Device = device
		res.Date = date
		res.State = state
	}

	return res, nil
}

func GetDistinctDevice() ([]types.Device, error) {
	db, err := getDB()
	if err != nil {
		return []types.Device{}, err
	}
	defer db.Close()

	deviceArray, err := db.Query("SELECT device, state FROM ( SELECT id, device, date, state, ROW_NUMBER() OVER (PARTITION BY device ORDER BY date DESC) AS rn FROM event_logs ) AS subquery WHERE rn = 1;")
	if err != nil {
		return []types.Device{}, err
	}
	defer deviceArray.Close()

	var res []types.Device
	for deviceArray.Next() {
		var tmp types.Device
		var device string
		var state bool
		err = deviceArray.Scan(&device, &state)
		if err != nil {
			return []types.Device{}, err
		}
		tmp.Device = device
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

	deviceArray, err := db.Query("SELECT state FROM ( SELECT id, device, date, state, ROW_NUMBER() OVER (PARTITION BY device ORDER BY date DESC) AS rn FROM event_logs ) AS subquery WHERE rn = 1;")
	if err != nil {
		return "", err
	}
	defer deviceArray.Close()

	var res []bool
	for deviceArray.Next() {
		var tmp bool
		var state bool
		err = deviceArray.Scan(&state)
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

// IsDeviceExist return tru if device exits in the database
func IsDeviceExist(deviceName string) (bool, error) {
	db, err := getDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var device string
	events := db.QueryRow("SELECT device FROM event_logs WHERE device=? ORDER BY id DESC LIMIT 1", deviceName)
	if err = events.Scan(&device); err != nil {
		return false, err
	}

	if device != "" {
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
