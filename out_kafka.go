package main

import (
  "github.com/fluent/fluent-bit-go/output"
  "github.com/ugorji/go/codec"
  "github.com/Shopify/sarama"
  "encoding/json"
  "reflect"
  "unsafe"
  "fmt"
  "io"
  "C"
)

var brokerList []string = []string{"kafka-0.kafka.default.svc.cluster.local:9092"}
var producer sarama.SyncProducer

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
  return output.FLBPluginRegister(ctx, "out_kafka", "out_kafka GO!")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
  var err error
  producer, err = sarama.NewSyncProducer(brokerList, nil)

  if err != nil {
    fmt.Printf("Failed to start Sarama producer: %v\n", err)
    return output.FLB_ERROR
  }

  return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
  var h codec.MsgpackHandle

  var b []byte
  var m interface{}
  var err error
  var enc_data []byte

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

    if format == "json" {
      enc_data, err = encode_as_json(m)
    } else if format == "msgpack" {
      enc_data, err = encode_as_msgpack(m)
    } else if format == "string" {
      // enc_data, err == encode_as_string(m)
    }
    if err != nil {
      fmt.Printf("Failed to encode %s data: %v\n", format, err)
      return output.FLB_ERROR
    }

    producer.SendMessage(&sarama.ProducerMessage {
      Topic: "logs_default",
      Key:   nil,
      Value: sarama.ByteEncoder(enc_data),
    })
    
  }
  return output.FLB_OK
}

  // record is map of interface to interfaces, which the json Marshaler will
  // not encode automagically. we need to iterate over it, and create a new
  // map of strings to interfaces that the json Marshaler can handle
func prepare_data(record interface{}) interface{} {
  // base case, 
  // if val is map, return
  r, ok := record.(map[interface{}] interface{})

  if ok != true {
    return record
  }

  record2 := make(map[string] interface{})
  for k, v := range r {
    key_string := k.(string)
    // convert C-style string to go string, else the JSON encoder will attempt
    // to base64 encode the array
    if val, ok := v.([]byte); ok {
      // if it IS a byte array, make string
      v2 := string(val)
      // add to new record map
      record2[key_string] = v2
    } else {
      // if not, recurse to decode interface &
      // add to new record map
      record2[key_string] = prepare_data(v)
    }
  }

  return record2
}

func encode_as_json(m interface {}) ([]byte, error) {
  slice := reflect.ValueOf(m)
  timestamp := slice.Index(0).Interface().(uint64)
  record := slice.Index(1).Interface()

  type Log struct {
    Time uint64
    Record interface{}
  }

  log := Log {
    Time: timestamp,
    Record: prepare_data(record),
  }

  return json.Marshal(log)
} 

func encode_as_msgpack(m interface {}) ([]byte, error) {
  var (
    mh codec.MsgpackHandle
    w io.Writer
    b []byte
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
