package web

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/winnerx0/jille/internal/utils"
)

func SseHandler(b *utils.Broker) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		client := make(chan utils.Event)
		b.Add <- client

		reader, writer := io.Pipe()

		go func() {
			defer func() {
				b.Remove <- client
				writer.Close()
			}()

			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case event, ok := <-client:
					if !ok {
						return
					}
					data, _ := json.Marshal(event)
					if _, err := fmt.Fprintf(writer, "data: %s\n\n", data); err != nil {
						return
					}
				case <-ticker.C:
					// Heartbeat to detect client disconnect
					if _, err := fmt.Fprintf(writer, ": keep-alive\n\n"); err != nil {
						return
					}
				}
			}
		}()

		return c.SendStream(reader)
	}
}
