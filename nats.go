package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NATSClient struct {
	conf *NATSConfig
	conn *nats.Conn
}

const (
	natsMaxReconnect  = 1000000
	natsReconnectWait = 5 * time.Second
)

func NewNATSClient(c *NATSConfig) (*NATSClient, error) {
	options := nats.GetDefaultOptions()
	options.Url = c.URL
	options.RetryOnFailedConnect = true
	options.ReconnectedCB = func(conn *nats.Conn) {
		logs.Warn("NATS reconnected successfully")
	}

	options.MaxReconnect = natsMaxReconnect
	options.ReconnectWait = natsReconnectWait
	options.Timeout = time.Duration(c.TimeoutSeconds) * time.Second
	options.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		logs.Errorf("NATS disconnected, err:%v", err)
	}

	conn, err := tryConnection(options)
	if err != nil {
		return nil, err
	}

	logs.Infof("nats client connected: %s", c.URL)

	return &NATSClient{
		conf: c,
		conn: conn,
	}, nil
}

func tryConnection(options nats.Options) (*nats.Conn, error) {
	conn, err := options.Connect()
	if err == nil && conn.IsConnected() {
		logrus.Infof("nats connected with normal config: %s", options.Url)
		return conn, nil
	}
	conn.Close()

	if err := nats.Secure(&tls.Config{InsecureSkipVerify: true})(&options); err != nil {
		return nil, err
	}

	conn, err = options.Connect()
	if err == nil && conn.IsConnected() {
		logrus.Infof("nats connected with insecure config: %s", options.Url)
		return conn, nil
	}
	conn.Close()

	return nil, fmt.Errorf("nats connecting with address '%s' finally failed", options.Url)
}

func (c *NATSClient) Publish(data []byte) error {
	return c.conn.Publish(c.conf.Subject, data)
}

type NATSConfig struct {
	ID             string
	URL            string
	Subject        string
	TimeoutSeconds int
}

func NewNATSConfig(id, urlStr, subject, timeoutSecondsStr string) (*NATSConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if urlStr == "" {
		return nil, errors.New("nats url is required")
	}

	if subject == "" {
		return nil, errors.New("nats subject is required")
	}

	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid url '%s' :%w", urlStr, err)
	}

	// default 5 seconds timeout
	timeoutSeconds := 5
	if timeoutSecondsStr != "" {
		timeoutSeconds, err = strconv.Atoi(timeoutSecondsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid TimeoutSeconds '%s' :%w", timeoutSecondsStr, err)
		}
	}

	return &NATSConfig{
		ID:             id,
		URL:            urlStr,
		Subject:        subject,
		TimeoutSeconds: timeoutSeconds,
	}, nil
}
