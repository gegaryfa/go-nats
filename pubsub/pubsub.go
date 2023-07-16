package main

import (
	"fmt"
	//"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

type Event struct {
	Message string
}

func main() {

	url := nats.DefaultURL

	nc, _ := nats.Connect(url)

	defer func(nc *nats.Conn) {
		_ = nc.Drain()
	}(nc)

	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer func(nc *nats.Conn) {
		_ = ec.Drain()
	}(nc)

	// this message will be lost because no one has subscribed to this subject yet
	_ = nc.Publish("greet.foo", []byte("hello"))

	// The wildcard in this synchronous subscriber means that we subscribe to all the subjects that start with <greet.>
	synvsub, _ := nc.SubscribeSync("greet.*")

	// This is an asynchronous subscriber that subscribes to only one subject
	// and needs an encoded connection
	_, _ = ec.Subscribe(
		"greet.async",
		func(event *Event) {
			fmt.Printf("message data on async subsriber: %s", event.Message)
		},
	)

	msg, _ := synvsub.NextMsg(10 * time.Millisecond)
	fmt.Println("subscribed after a publish...")
	fmt.Printf("msg is nil? %v\n", msg == nil)

	_ = nc.Publish("greet.joe", []byte("hello"))
	_ = nc.Publish("greet.pam", []byte("hello"))

	msg, _ = synvsub.NextMsg(10 * time.Millisecond)
	if msg != nil {
		fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
	}

	msg, _ = synvsub.NextMsg(10 * time.Millisecond)
	if msg != nil {
		fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
	}

	_ = nc.Publish("greet.bar", []byte("hello"))

	msg, _ = synvsub.NextMsg(10 * time.Millisecond)
	if msg != nil {
		fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
	}

	event := Event{
		Message: "Async message",
	}
	// Publish using JSON encoder. this one will be picked up by the async subscriber
	_ = ec.Publish("greet.async", event)

	// we sleep to make sure that the async subscriber has enough time to react to the received message
	time.Sleep(1 * time.Second)
}
