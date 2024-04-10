package misc

import (
	"backend/db"
	"backend/types"
	"backend/vars"
	"encoding/hex"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"
)

var xorKey = func() []byte {
	keyString := vars.GetKey()
	if keyString == "" {
		slog.Error("Environment variable not set", slog.String("env", keyString))
		os.Exit(1)
	}

	key, err := hex.DecodeString(keyString)
	if err != nil {
		slog.Error("Can't decode XOR key", slog.Any("error", err))
	}

	return key
}()

// XorData takes an input string and XOR it with a key
func XorData(input string) (output string) {
	if len(input) != len(xorKey) {
		return ""
	}

	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ xorKey[i%len(xorKey)])
	}

	return output
}

// SetupDevices read DEVICES system environment and assign devices
func SetupDevices() error {
	deviceList := os.Getenv("DEVICES")
	devices := strings.Fields(deviceList)

	slog.Info("Adding devices to the database", slog.Any("devices", devices))
	for _, device := range devices {
		d := types.Device{
			Device: device,
			State:  false,
		}
		if _, dbErr := db.AddRecord(d); dbErr != nil {
			slog.Error("Error adding device to the database", slog.Any("error", dbErr),
				slog.String("device", device))
			return dbErr
		}
	}
	return nil
}

func Seed() error {
	devices, err := db.GetDistinctDevice()
	if err != nil {
		slog.Error("Error getting distinct devices", slog.Any("error", err))
	}
	states := make(map[string]bool)

	slog.Info("Seeding started")
	for _, device := range devices {
		states[device.Device] = false
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	records := 100
	for i := 0; i < records; i++ {

		// Generate random like change
		index := random.Intn(len(devices))
		device := devices[index].Device
		states[device] = !states[device]
		randomTime := random.Intn(360) + 30

		// Insert into database
		record := types.Device{
			Device: device,
			State:  states[device],
		}
		if _, err := db.AddSeededRecord(record, randomTime); err != nil {
			slog.Error("Error seeding database with record", slog.String("key", device), slog.Bool("value", states[device]))
		}
	}
	slog.Info("Seeding completed", slog.Int("records_added", records))

	return nil
}
