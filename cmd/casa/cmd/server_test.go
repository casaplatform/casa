package cmd

import (
	"testing"

	"github.com/casaplatform/casa"
	"github.com/gomqtt/broker"
	"github.com/gomqtt/packet"
)

func Test_getRand(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRand(); got != tt.want {
				t.Errorf("getRand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_brokerLogger_Log(t *testing.T) {
	type fields struct {
		Logger casa.Logger
	}
	type args struct {
		event   broker.LogEvent
		client  *broker.Client
		packet  packet.Packet
		message *packet.Message
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &brokerLogger{
				Logger: tt.fields.Logger,
			}
			bl.Log(tt.args.event, tt.args.client, tt.args.packet, tt.args.message, tt.args.err)
		})
	}
}
