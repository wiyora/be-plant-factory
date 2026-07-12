package websocket

import (
	"encoding/json"

	"github.com/gofiber/contrib/v3/socketio"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, handler SocketHandler) {
	handler.SubscribeEvents()

	app.Get("/ws", socketio.New(func(kws *socketio.Websocket) {
		// Per-connection setup handled by global events via SubscribeEvents
	}))

	app.Get("/health/ws", func(c fiber.Ctx) error {
		payload, _ := json.Marshal(map[string]interface{}{
			"clients": handler.ClientCount(),
		})
		return c.SendString(string(payload))
	})
}
