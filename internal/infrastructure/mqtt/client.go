package mqtt

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Client struct {
	client mqtt.Client
	log    zerolog.Logger
}

func New(i do.Injector) (*Client, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	cfg := do.MustInvoke[*config.Config](i).MQTT

	log := logger.WithLayer(rawLog, logger.LayerMQTT)
	log.Info().Msg("connecting to mqtt broker")

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Address()).
		SetClientID(cfg.ClientIdentifier()).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			log.Info().Msg("connected to mqtt broker")
		}).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			log.Error().Err(err).Msg("mqtt connection lost")
		})

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("failed to connect to mqtt broker")
		return nil, fmt.Errorf("connect mqtt: %w", token.Error())
	}

	log.Info().Msg("successfully connected to mqtt broker")

	return &Client{
		client: client,
		log:    log,
	}, nil
}

func (c *Client) Client() mqtt.Client {
	return c.client
}

func (c *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, qos, callback)
	token.Wait()
	return token.Error()
}

func (c *Client) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) error {
	token := c.client.SubscribeMultiple(filters, callback)
	token.Wait()
	return token.Error()
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) HealthCheckWithContext(ctx context.Context) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("mqtt client is not connected")
	}

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	if c.client != nil && c.client.IsConnected() {
		c.log.Info().Msg("disconnecting from mqtt broker")
		c.client.Disconnect(250)
	}

	return nil
}
