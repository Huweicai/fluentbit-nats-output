package main

import (
	"C"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/sirupsen/logrus"
)

const (
	PluginName = "nats-output"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	var formatter logrus.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)
}

var logs = logrus.WithField("plugin", PluginName)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, PluginName, "Sending data to nats.")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	logs = logs.WithField("stage", "FLBPluginInit")

	config, err := NewNATSConfig(output.FLBPluginConfigKey(plugin, "ID"),
		output.FLBPluginConfigKey(plugin, "URL"),
		output.FLBPluginConfigKey(plugin, "Subject"),
		output.FLBPluginConfigKey(plugin, "TimeoutSeconds"))
	if err != nil {
		logs.WithError(err).Errorf("invalid config")
		return output.FLB_ERROR
	}
	logs.WithField("config", config).Info("config loaded")

	client, err := NewNATSClient(config)
	if err != nil {
		logs.WithError(err).Error("nats client init failed")
		return output.FLB_ERROR
	}

	// Set the context to point to any Go variable
	output.FLBPluginSetContext(plugin, client)

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	logs.Info("flush called without context")
	return output.FLB_ERROR
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	logs = logs.WithField("stage", "FLBPluginFlushCtx")

	// Type assert context back into the original type for the Go variable
	client := output.FLBPluginGetContext(ctx).(*NATSClient)
	dec := output.NewDecoder(data, int(length))

	logs.Debug("data flushed")

	for {
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		var timestamp time.Time
		switch t := ts.(type) {
		case output.FLBTime:
			timestamp = ts.(output.FLBTime).Time
		case uint64:
			timestamp = time.Unix(int64(t), 0)
		default:
			fmt.Println("time provided invalid, defaulting to now.")
			timestamp = time.Now()
		}

		flattened, err := Flatten(record, "", DotStyle)
		if err != nil {
			logs.WithError(err).WithField("record", record).Warn("flattern failed")
			continue
		}

		flattened["__time"] = timestamp.Format(time.RFC3339)

		body, err := json.Marshal(flattened)
		if err != nil {
			logs.WithError(err).WithField("record", record).Warn("invalid unmarshalble record")
			continue
		}

		if err := client.Publish(body); err != nil {
			logs.WithError(err).WithField("body", body).Warn("publish data failed")
			continue
		}

	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	log.Print("[multiinstance] Exit called for unknown instance")
	return output.FLB_OK
}

//export FLBPluginExitCtx
func FLBPluginExitCtx(ctx unsafe.Pointer) int {
	return output.FLB_OK
}

func main() {
}
