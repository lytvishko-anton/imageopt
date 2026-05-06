package main

import (
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	amqp "github.com/rabbitmq/amqp091-go"
)

// This matches the inner data sent by Symfony
type ImageTask struct {
	FilePath string `json:"filePath"`
}

// Symfony Messenger wraps the object in a 'message' key
type SymfonyEnvelope struct {
	Message ImageTask `json:"message"`
}

func main() {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, _ := ch.QueueDeclare("messages", true, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	log.Println(" [*] WORKER READY. Awaiting images...")

	for d := range msgs {
		log.Printf("Raw Data Received: %s", string(d.Body))

		var data map[string]interface{}
		if err := json.Unmarshal(d.Body, &data); err != nil {
			log.Printf("Unmarshal Error: %v", err)
			continue
		}

		// Grab 'filePath' directly from the top level
		path, ok := data["filePath"].(string)
		if !ok || path == "" {
			log.Printf("Error: 'filePath' is empty or missing in JSON")
			continue
		}

		log.Printf("SUCCESS! Found Path: %s", path)

		// Verify it exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("CRITICAL: File physically missing at %s", path)
			continue
		}

		// Convert and delete
		err := convertToWebp(path)
		if err != nil {
			log.Printf("Conversion Failed: %v", err)
			continue
		}

		// Delete the ORIGINAL (jpg/png) now that we have the webp
		err = os.Remove(path)
		if err != nil {
			log.Printf("Could not delete original: %v", err)
		} else {
			log.Printf("SUCCESS: Original deleted. Only WebP remains for download.")
		}
	}
}

func convertToWebp(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".webp"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return webp.Encode(outputFile, img, &webp.Options{Lossless: false, Quality: 75})
}
