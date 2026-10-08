package services

import (
	"log"
	"time"
)

func StartHealthCheckScheduler() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		log.Println("Health check scheduler started")

		for {
			<-ticker.C

			log.Println("Running scheduled health checks...")
			RunHealthChecks()
		}
	}()
}
