package main

import (
	"context"
	"os"
	"time"

	"github.com/Unfield/Statuz/logger"
	"github.com/Unfield/Statuz/monitors"
	"github.com/Unfield/Statuz/scheduler"
)

func main() {
	os.Setenv("APP_ENV", "production")
	os.Setenv("LOGGER_MODE", "cute")

	log := logger.NewLogger()

	log.Info("🚀 Starting Statuz...")

	monitor := monitors.HTTPMonitor{
		ID:                  "test-monitor",
		URL:                 "https://www.speedtest.net/",
		HBInterval:          10 * time.Second,
		Retries:             3,
		RetryInterval:       2 * time.Second,
		ReqTimeout:          5 * time.Second,
		MaxRedirects:        0,
		AcceptedStatusCodes: []int{200},
		IPFamily:            "auto",
		HTTPMethod:          "GET",
	}

	monitor2 := monitors.HTTPMonitor{
		ID:                  "test-monitor2",
		URL:                 "https://www.google.com/",
		HBInterval:          60 * time.Second,
		Retries:             3,
		RetryInterval:       2 * time.Second,
		ReqTimeout:          5 * time.Second,
		MaxRedirects:        0,
		AcceptedStatusCodes: []int{200},
		IPFamily:            "auto",
		HTTPMethod:          "GET",
	}

	monitor3 := monitors.HTTPMonitor{
		ID:                  "test-monitor3",
		URL:                 "https://www.google.com/",
		HBInterval:          10 * time.Second,
		Retries:             3,
		RetryInterval:       2 * time.Second,
		ReqTimeout:          5 * time.Second,
		MaxRedirects:        0,
		AcceptedStatusCodes: []int{200},
		IPFamily:            "auto",
		HTTPMethod:          "GET",
		DegradedThreshold:   10 * time.Millisecond,
	}

	monitor4 := monitors.HTTPMonitor{
		ID:                  "test-monitor4",
		URL:                 "https://www.thisdomaincannotbevalid438753984tu38ut3u09rt309.com/",
		HBInterval:          60 * time.Second,
		Retries:             3,
		RetryInterval:       2 * time.Second,
		ReqTimeout:          5 * time.Second,
		MaxRedirects:        0,
		AcceptedStatusCodes: []int{200},
		IPFamily:            "auto",
		HTTPMethod:          "GET",
	}

	monitorsList := []monitors.Monitor{&monitor, &monitor2, &monitor3, &monitor4}

	s := scheduler.NewScheduler(context.Background(), monitorsList)
	go s.Run()

	for {
		select {
		case result, ok := <-s.ResultChannel:
			if !ok {
				log.Warn("Result channel closed")
				return
			}

			log.Infof("%s [%s] %s → %s (success=%t) (duration=%d)",
				monitors.GetStatusIcon(result.Status),
				result.MonitorID,
				result.Type,
				result.Status,
				result.Success,
				result.Duration.Milliseconds(),
			)

		case <-s.Context().Done():
			log.Warn("Scheduler context canceled, stopping listener.")
			return
		}
	}
}
