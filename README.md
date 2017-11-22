Kafka to Watson
===============

Simple program to fetch data from one Kafka topic, and send it to IBM Watson, pretending we are an IOT Gateway.

The kafka message key MUST be of the form:

    deviceType/deviceId/event

it describes the device to impersonate.

# Usage

    kafka-to-watson <broker1,...,brokerN> <topic1,...,topicN> <watson-org> <watson-gw-type> <watson-gw-id> <watson-gw-token>

# Example

1. Start the kafka consumer:

        ./kafka-to-watson 127.0.0.1:9092 watson tpktfp CurlGatewayType CurlGateway 'XXX...YYY'

2. Push a message to the `watson` kafka topic:

        kafkacat -P -t 'watson' -b localhost:9092 -K '#' <<< 'TestDevice/curldevice/cpu#{"temp":40}'

# TODO

- [x] HTTP Watson client
- [ ] MQTT Watson client