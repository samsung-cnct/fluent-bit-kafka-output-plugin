# fluent-bit-kafka-output-plugin
Kafka Output Plugin for FluentBit

To deploy as a sidecar in Kubernetes: https://github.com/leahnp/fluentbit-sidecar

Message me with any comments or suggestions!

Info on FluentBit: https://github.com/fluent/fluent-bit

Info on Golang Output Plugins for FluentBit: https://github.com/fluent/fluent-bit-go

### Run locally for development:

Follow the instructions [here](https://kafka.apache.org/quickstart) to start zookeeper and kafka locally on your computer.

Download [Fluent-bit](http://fluentbit.io/download/) locally. [Build Fluent-bit](http://fluentbit.io/documentation/0.12/installation/build_install.html). Test Fluent-bit is installed and built properly by running: `bin/fluent-bit -i random -o stdout` from your Fluent-bit build directory. This command starts Fluent-bit, tells it to use "random" for the input plugin and "stdout" for the output plugin.

Now we want to run Fluent-bit with the kafka output plugin. `bin/fluent-bit -v -e <path to ../fluent-bit-kafka-output-plugin/out_kafka.so> -i random -o out_kafka`. Remember to change the broker list in the plugin to: "localhost:9092" and re-"MAKE". 

Note: make sure you configure the proper topic locally - currently Fluent-bit-output-plugin uses: `logs_default` but the Kafka quickstart only adds `test` topic. 




