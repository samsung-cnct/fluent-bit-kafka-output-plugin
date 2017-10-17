package main

import (
	"context"
	"testing"
	"time"
	"unsafe"

	"github.com/Shopify/sarama/mocks"
	"github.com/fluent/fluent-bit-go/output"
)

func TestFLBPluginRegister(t *testing.T) {
	bk := context.Background()
	var ctx = unsafe.Pointer(&bk)

	result := FLBPluginRegister(ctx)
	if result != 0 {
		t.Error("Failed to register plugin, expected", 0, "but found", result)
	}
}

func TestFLBPluginInitFails(t *testing.T) {
	bk := context.Background()
	var ctx = unsafe.Pointer(&bk)

	timeout = 1 * time.Minute
	result := FLBPluginInit(ctx)
	if result != output.FLB_ERROR {
		t.Error("Expected kafka producer to error out, but found success instead", result)
	}
}

func TestFLBPluginInitSucceeds(t *testing.T) {
	bk := context.Background()
	var ctx = unsafe.Pointer(&bk)
	timeout = 0
	producer = mocks.NewSyncProducer(t, nil)
	result := FLBPluginInit(ctx)
	if result != output.FLB_OK {
		t.Error("Expected kafka to successfully connect, but failed instead", result)
	}
}
