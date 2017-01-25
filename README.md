# fluent-bit-kafka-output-plugin
Kafka Output Plugin for FluentBit

Status (01-25-17): In development, this plugin is functioning and will send json from fluentbit to kafka. 

To deploy as a sidecar in Kubernetes: https://github.com/leahnp/fluentbit-sidecar

Message me with any comments or suggestions!

Info on FluentBit: https://github.com/fluent/fluent-bit

Info on Golang Output Plugins for FluentBit: https://github.com/fluent/fluent-bit-go

TODO: 
- Fluentbit still needs Kubernetes metadata filter, waiting for Golang plugin filter capabilities. 
- Add string and test msgpack encoding
- Add ability to connect via brokers or zookeeper
- Add variable to conf file for topic
