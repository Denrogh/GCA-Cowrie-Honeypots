package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Denrogh/GCA-Cowrie-Honeypots/pkg/config"
	"github.com/Denrogh/GCA-Cowrie-Honeypots/pkg/kafka"
)

func main() {
	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if cfg.LogPath == "" {
		cfg.LogPath = "/app/test.json"
	}

	// Login and get auth token
	// CHANGE TO LOAD UP LOG FILE, relative path > absolute path for some reason//
	// fmt.Println(cfg.LogPath)
	logfile, err := os.Open(cfg.LogPath)
	if err != nil {
		log.Fatalf("failed to load log file: %v: %v", cfg.LogPath, err)
		return
	}

	fmt.Println("Log file loaded successfully")

	// Create a new instance of the Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer producer.Close()

	fmt.Println("producer returned successfully")
	fmt.Println(cfg.KafkaBrokers, cfg.KafkaTopic)

	// Fetch sessions from the logfile (reading the json file, care for offset. limit, before and after ONLY (use some sort of filter))
	// TODO: methodise
	scanner := bufio.NewScanner(logfile)

	// Index tracks how far in log file, used for polling.
	index := 0

	for scanner.Scan() {
		line := scanner.Text()
		index++

		// Extract the first string from the line and place it in HoneypotName
		segments := strings.Fields(line)

		// Ignores initialization lines with > 1
		if len(segments) > 1 {
			// var session kafka.CowrieSession
			session := kafka.CowrieSession{HoneypotName: segments[0]}
			braceIndex := strings.Index(line, "{")
			if braceIndex >= 0 {
				modifiedLine := line[braceIndex:]

				if err = json.Unmarshal([]byte(modifiedLine), &session); err != nil {
					log.Printf("Error unmarshaling JSON: %v", err)
					continue // Skip this line if unmarshaling fails
				}
			}

			if err = producer.PublishSession(session, cfg.KafkaTopic); err != nil {
				log.Fatalf("Failed to publish session to Kafka:", err)
			}

			fmt.Printf("Session sent successfully")
		}
	}
	logfile.Close()

	fmt.Println("Logs from before kafka-producer start sent, now starting polling if enabled")

	// Check if polling is enabled
	if cfg.EnablePolling {
		// Continuously poll for new sessions based on the configured interval
		pollingInterval, err := time.ParseDuration(cfg.PollingInterval)
		if err != nil {
			log.Fatal("Failed to parse pollingInterval value:", err)
		}

		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()

		// Poll at every interval, this should run until stopped
		for range ticker.C {
			pollfile, err := os.Open("/app/test.json")
			if err != nil {
				fmt.Printf("Failed to load polledlog file, check path in config ", err)
				return
			}
			defer pollfile.Close()
			// File has to be reread at the beginning of every poll
			startTime := time.Now()
			scanner2 := bufio.NewScanner(pollfile)
			pollingIndex := 0

			fmt.Println(index)

			for scanner2.Scan() {
				line := scanner2.Text()
				pollingIndex++

				// The pollingIndex is now bigger, read and publish new sessions
				segments := strings.Fields(line)
				if len(segments) > 1 && (index < pollingIndex) {
					fmt.Println("this is being ran")
					session := kafka.CowrieSession{HoneypotName: segments[0]}
					braceIndex := strings.Index(line, "{")
					if braceIndex >= 0 {
						modifiedLine := line[braceIndex:]

						if err = json.Unmarshal([]byte(modifiedLine), &session); err != nil {
							log.Printf("Error unmarshaling JSON: %v", err)
							continue // Skip this line if unmarshaling fails
						}
					}

					if err = producer.PublishSession(session, cfg.KafkaTopic); err != nil {
						log.Printf("Failed to publish session to Kafka:", err)
					}
					fmt.Println("Polled message published!")
				}
			}

			// Mark where in the logs was reached so next polling loop starts from there
			index = pollingIndex
			fmt.Printf("Session sent successfully")
			pollfile.Close()
			elapsedTime := time.Since(startTime)
			fmt.Printf("Polling loop took %s to run\n", elapsedTime)
		}

		log.Println("Polling has come to an end")
	}
}
