Kafka to Watson
===============

Simple program to fetch data from one Kafka topic, and send it to IBM Watson, pretending we are an IOT Gateway.

The kafka message MUST be of the form:

    {
        "DeviceType": "...",
        "DeviceId": "...",
        "Event": "...",
        "Message": { ... }
    }

it describes the device to impersonate.

# Usage

    kafka-to-watson <broker1,...,brokerN> <topic1,...,topicN> <watson-org> <watson-gw-type> <watson-gw-id> <watson-gw-token>

# Example

1. Start the kafka consumer:

        ./kafka-to-watson 127.0.0.1:9092 watson tpktfp CurlGatewayType CurlGateway 'XXX...YYY'

2. Push a message to the `watson` kafka topic:

        kafkacat -P -t 'watson' -b localhost:9092 <<< '{"DeviceType": "Thermometer", "DeviceId": "th032", "Event": "newState", "Message": {"temp":40}}'

# TODO

- [x] HTTP Watson client
- [ ] MQTT Watson client
