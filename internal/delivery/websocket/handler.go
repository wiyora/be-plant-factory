package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/contrib/v3/socketio"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	mqttInfra "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/mqtt"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type SocketHandler interface {
	SubscribeEvents()
	OnMQTTMessage(topic string, qos byte, callback mqttLib.MessageHandler) error
	PublishMQTT(topic string, qos byte, retained bool, payload interface{}) error
	BroadcastToClients(event string, data []byte)
	ClientCount() int
	EmitToClient(uuid, event string, data []byte) error
}

type socketHandler struct {
	mqttClient *mqttInfra.Client
	log        zerolog.Logger
	mu         sync.RWMutex
	clients    map[string]*socketio.Websocket
}

func NewSocketHandler(i do.Injector) (SocketHandler, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	mqttClient := do.MustInvoke[*mqttInfra.Client](i)

	log := logger.WithLayer(rawLog, logger.LayerSocketIO)
	log.Info().Msg("initializing socket.io handler")

	return &socketHandler{
		mqttClient: mqttClient,
		log:        log,
		clients:    make(map[string]*socketio.Websocket),
	}, nil
}

func (h *socketHandler) SubscribeEvents() {
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		id := ep.Kws.GetUUID()
		h.mu.Lock()
		h.clients[id] = ep.Kws
		h.mu.Unlock()

		h.log.Info().Str("uuid", id).Msg("client connected")
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		id := ep.Kws.GetUUID()
		h.mu.Lock()
		delete(h.clients, id)
		h.mu.Unlock()

		h.log.Info().Str("uuid", id).Msg("client disconnected")
	})

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		h.log.Debug().Str("uuid", ep.Kws.GetUUID()).Msg("message received")
	})
}

func (h *socketHandler) OnMQTTMessage(topic string, qos byte, callback mqttLib.MessageHandler) error {
	return h.mqttClient.Subscribe(topic, qos, callback)
}

func (h *socketHandler) PublishMQTT(topic string, qos byte, retained bool, payload interface{}) error {
	return h.mqttClient.Publish(topic, qos, retained, payload)
}

func (h *socketHandler) BroadcastToClients(event string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, kws := range h.clients {
		if kws.IsAlive() {
			kws.EmitEvent(event, data)
		}
	}
}

func (h *socketHandler) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *socketHandler) EmitToClient(uuid, event string, data []byte) error {
	h.mu.RLock()
	kws, ok := h.clients[uuid]
	h.mu.RUnlock()

	if !ok || !kws.IsAlive() {
		return fmt.Errorf("client %s not found or disconnected", uuid)
	}

	kws.EmitEvent(event, data)
	return nil
}

type MQTTMessage struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}
