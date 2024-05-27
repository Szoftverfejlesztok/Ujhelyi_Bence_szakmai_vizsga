package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"backend/db"
	"backend/misc"
	"backend/router"
	"backend/vars"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Seed the running database")
	flag.Parse()

	// Wait for MySQL to start
	for i := 0; i < vars.GetMaxTry(); i++ {
		if err := db.HealthCheck(); err == nil {
			slog.Info("Database reached", slog.Int("attempt", i+1))
			break
		} else {
			slog.Info("Tried to reaching database", slog.Int("attempt", i+1))
		}

		if i == vars.GetMaxTry()-1 {
			slog.Error("Could not connect to the database, exiting...")
			os.Exit(1)
		}

		time.Sleep(5 * time.Second)
	}

	// Seed database with random data
	if *seedFlag {
		slog.Info("Running in seed mode")
		if err := misc.Seed(); err != nil {
			slog.Error("Error during seeding database", slog.Any("error", err))
		}
		os.Exit(0)
	}

	// Setup devices
	devices, err := db.GetDistinctDevice()
	if err != nil {
		slog.Error("Error getting devices", slog.Any("error", err))
		os.Exit(1)
	}

	if len(devices) == 0 {
		slog.Info("Starting device setup")
		if err := misc.SetupDevices(); err != nil {
			slog.Error("Error during device setup", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// HTTP
	go func() {
		var port = vars.GetPort()
		r := chi.NewRouter()
		c := cors.New(cors.Options{
			AllowedOrigins:  []string{"*"},
			AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
			AllowedMethods:  []string{"GET", "POST"},
		})
		// Backend
		r.Use(c.Handler)
		r.Post("/api/addRecord", router.AddRecordHandler)
		r.Get("/api/getLastByDevice/{device}", router.GetLastByDeviceHandler)
		r.Get("/api/getDevices", router.GetDevices)
		r.Get("/api/getDevicesUptime", router.GetDevicesUptime)
		r.Get("/api/hc", router.HealthCheckHandler)
		// Frontend
		router.FileServer(r, "/", http.Dir("./frontend"))
		router.FileServer(r, "/statistics", http.Dir("./frontend/statistics"))

		slog.Info("Smarthome API is running", slog.String("port", port))
		if err := http.ListenAndServe(":"+port, r); err != nil {
			slog.Error("Could not serve HTTP API", slog.String("port", port))
		}
	}()

	// WS
	go func() {
		var port = vars.GetWSPort()

		slog.Info("Smarthome WS is running", slog.String("port", port), slog.String("path", "/smart-home"))

		http.HandleFunc("/smart-home", router.HandleClient)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			slog.Error("Error serving WebSocket port", slog.String("port", port))
		}
	}()

	select {}
}
