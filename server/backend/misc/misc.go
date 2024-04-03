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

// SetupLamps read ROOMS system environment and assign one lamp to it
// This function os called only once before everything else starts
func SetupLamps() error {
	rooms := os.Getenv("ROOMS")
	lamps := strings.Fields(rooms)

	slog.Info("Adding lamps to the database", slog.Any("lamps", lamps))
	for _, lamp := range lamps {
		l := types.Lamp{
			Lamp:  lamp,
			State: false,
		}
		if _, dbErr := db.AddRecord(l); dbErr != nil {
			slog.Error("Error adding lamp to the database", slog.Any("error", dbErr),
				slog.String("lamp", lamp))
			return dbErr
		}
	}
	return nil
}
