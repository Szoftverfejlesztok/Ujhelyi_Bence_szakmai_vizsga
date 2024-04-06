package misc

import (
	"backend/db"
	"backend/types"
	"encoding/hex"
	"log/slog"
	"os"
	"strings"
)

var xorKey = func() []byte {
	keyString := os.Getenv("XOR_KEY")
	if keyString == "" {
		slog.Error("Environment variable not set", slog.String("env", "XOR_KEY"))
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

// SetupDevices read ROOMS system environment and assign one device to it
// This function os called only once before everything else starts
func SetupDevices() error {
	rooms := os.Getenv("ROOMS")
	devices := strings.Fields(rooms)

	slog.Info("Adding devices to the database", slog.Any("devices", devices))
	for _, device := range devices {
		l := types.Device{
			Device: device,
			State:  false,
		}
		if _, dbErr := db.AddRecord(l); dbErr != nil {
			slog.Error("Error adding device to the database", slog.Any("error", dbErr),
				slog.String("device", device))
			return dbErr
		}
	}
	return nil
}
