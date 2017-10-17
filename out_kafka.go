// Copyright © 2017 Samsung CNCT
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Adapted from https://github.com/fluent/fluent-bit/blob/master/GOLANG_OUTPUT_PLUGIN.md

package main

import (
	"C"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"
	"unsafe"

	"github.com/Shopify/sarama"
	"github.com/fluent/fluent-bit-go/output"
	"github.com/ugorji/go/codec"
)

var brokerList = []string{"kafka-0.kafka.default.svc.cluster.local:9092"}
var producer sarama.SyncProducer
var timeout = 0 * time.Minute

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "out_kafka", "out_kafka GO!")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	var err error

	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	// If Kafka is not running on init, wait to connect
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		if producer == nil {
			producer, err = sarama.NewSyncProducer(brokerList, nil)
		}
		if err == nil {
			return output.FLB_OK
		}
		log.Printf("Cannot connect to Kafka: (%s) retrying...", err)
		time.Sleep(time.Second * 30)
	}
	log.Printf("Kafka failed to respond after %s", timeout)
	return output.FLB_ERROR
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	var h codec.MsgpackHandle

	var b []byte
	var m interface{}
	var err error
	var encData []byte

	b = C.GoBytes(data, length)
	dec := codec.NewDecoderBytes(b, &h)

	// Iterate the original MessagePack array
	for {
		// decode the msgpack data
		err = dec.Decode(&m)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Failed to decode msgpack data: %v\n", err)
			return output.FLB_ERROR
		}

		// select format until config files are available for fluentbit
		format := "json"

		switch format {
		case "json":
			encData, err = encodeAsJson(m)
		case "msgpack":
			encData, err = encodeAsMsgpack(m)
		case "string":
			// encData, err == encode_as_string(m)
		}

		if err != nil {
			fmt.Printf("Failed to encode %s data: %v\n", format, err)
			return output.FLB_ERROR
		}

		producer.SendMessage(&sarama.ProducerMessage{
			Topic: "logs_default",
			Key:   nil,
			Value: sarama.ByteEncoder(encData),
		})

	}
	return output.FLB_OK
}

// record is map of interface to interfaces, which the json Marshaler will
// not encode automagically. we need to iterate over it, and create a new
// map of strings to interfaces that the json Marshaler can handle
func prepareData(record interface{}) interface{} {
	// base case,
	// if val is map, return
	r, ok := record.(map[interface{}]interface{})

	if ok != true {
		return record
	}

	record2 := make(map[string]interface{})
	for k, v := range r {
		keyString := k.(string)
		// convert C-style string to go string, else the JSON encoder will attempt
		// to base64 encode the array
		if val, ok := v.([]byte); ok {
			// if it IS a byte array, make string
			v2 := string(val)
			// add to new record map
			record2[keyString] = v2
		} else {
			// if not, recurse to decode interface &
			// add to new record map
			record2[keyString] = prepareData(v)
		}
	}

	return record2
}

func encodeAsJson(m interface{}) ([]byte, error) {
	slice := reflect.ValueOf(m)
	timestamp := slice.Index(0).Interface().(uint64)
	record := slice.Index(1).Interface()

	type Log struct {
		Time   uint64
		Record interface{}
	}

	log := Log{
		Time:   timestamp,
		Record: prepareData(record),
	}

	return json.Marshal(log)
}

func encodeAsMsgpack(m interface{}) ([]byte, error) {
	var (
		mh codec.MsgpackHandle
		w  io.Writer
		b  []byte
	)

	enc := codec.NewEncoder(w, &mh)
	enc = codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(&m)
	return b, err
}

func FLBPluginExit() int {
	return 0
}

func main() {
}
