package main

import (
	"bytes"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"net/http"
	"os"
	"strings"
)

var WATSONURITMPL = "https://%s.messaging.internetofthings.ibmcloud.com/api/v0002/device/types/%s/devices/%s/events/%s"

type WatsonConnection interface {
	Init() error
	Close() error
	Publish(deviceType string, deviceId string, event string, payload []byte) error
}

type WatsonHTTP struct {
	org     string
	gwType  string
	gwId    string
	gwToken string

	auth string
}

type WatsonMQTT struct {
	// XXX
}

func (w *WatsonHTTP) Init() error {
	w.auth = fmt.Sprintf("g/%s/%s/%s", w.org, w.gwType, w.gwId)
	return nil
}

func (w *WatsonHTTP) Close() error {
	return nil
}

func (w *WatsonHTTP) Publish(deviceType string, deviceId string, event string, payload []byte) error {
	uri := fmt.Sprintf(WATSONURITMPL, w.org, deviceType, deviceId, event)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(w.auth, w.gwToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (w *WatsonMQTT) Init(org string, gwType string, gwId string, gwToken string) error {
	// XXX
	return nil
}

func (w *WatsonMQTT) Close() error {
	// XXX
	return nil
}

func (w *WatsonMQTT) Publish(deviceType string, deviceId string, event string, payload []byte) error {
	// XXX
	return nil
}

func usage() {
	println("Usage: kafka-to-watson <broker1,...,brokerN> <topic1,...,topicN> <watson-org> <watson-gw-type> <watson-gw-id> <watson-gw-token>")
}

func main() {
	args := os.Args[1:]
	if len(args) != 6 {
		usage()
		os.Exit(1)
	}

	brokers := strings.Split(args[0], ",")
	topics := strings.Split(args[1], ",")
	watsonOrg := args[2]
	watsonGwType := args[3]
	watsonGwId := args[4]
	watsonGwToken := args[5]

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	consumer, err := cluster.NewConsumer(brokers, "kafka-to-watson", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	log.Printf("Kafka Consumer created")

	go func() {
		for err := range consumer.Errors() {
			log.Printf("[KAFKA] Errors: %s", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("[KAFKA] Rebalanced: %+v", ntf)
		}
	}()

	var c WatsonConnection = &WatsonHTTP{
		org:     watsonOrg,
		gwType:  watsonGwType,
		gwId:    watsonGwId,
		gwToken: watsonGwToken,
	}
	if err := c.Init(); err != nil {
		log.Panic("[WATSON] Unable to create HTTP connection")
	}
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				log.Printf("New messages received Topic: '%s', Key: '%s', Message: '%s'",
					msg.Topic, msg.Key, msg.Value)

				device := strings.Split(string(msg.Key), "/")
				if len(device) == 3 {
					deviceType, deviceId, event := device[0], device[1], device[2]
					if err := c.Publish(deviceType, deviceId, event, msg.Value); err != nil {
						log.Printf("[WATSON] Error while publishing data to Watson, message skipped…")
					}
				} else {
					log.Printf("[KAFKA] Key must be of the form 'deviceType/deviceId/event', message skipped…")
				}
				consumer.MarkOffset(msg, "")
			}
		}
	}
	print("This is the end")
}
