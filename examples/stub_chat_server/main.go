// Stub Chat + AIControl implementation. Proves consumers can implement
// the kotg-schema service interfaces.
//
// Run: go run ./examples/stub_chat_server
package main

import (
	"context"
	"log"
	"net"

	kotgv1 "github.com/vellankikoti/kotg-schema/gen/go/kotg/v1"
	"google.golang.org/grpc"
)

type stubChat struct {
	kotgv1.UnimplementedChatServer
}

func (s *stubChat) Send(stream kotgv1.Chat_SendServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil { return err }
		// Echo: stream a TextDelta + Done back.
		if err := stream.Send(&kotgv1.AssistantEvent{
			AnchorId: "1",
			Event:    &kotgv1.AssistantEvent_TextDelta{TextDelta: &kotgv1.TextDelta{Text: "echo: " + msg.GetText()}},
		}); err != nil { return err }
		if err := stream.Send(&kotgv1.AssistantEvent{
			AnchorId: "2",
			Event:    &kotgv1.AssistantEvent_Done{Done: &kotgv1.Done{FinishReason: "stop"}},
		}); err != nil { return err }
	}
}

type stubAIControl struct {
	kotgv1.UnimplementedAIControlServer
}

func (s *stubAIControl) Capabilities(_ context.Context, _ *kotgv1.Empty) (*kotgv1.AICapabilities, error) {
	return &kotgv1.AICapabilities{
		SchemaVersion: "1.0.0",
		AiVersion:     "0.0.1-stub",
		Providers:     []string{"stub"},
		Models:        []string{"stub-1"},
		SupportsUndo:  true,
		SupportsPlans: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil { log.Fatal(err) }
	srv := grpc.NewServer()
	kotgv1.RegisterChatServer(srv, &stubChat{})
	kotgv1.RegisterAIControlServer(srv, &stubAIControl{})
	log.Println("stub_chat_server listening on 127.0.0.1:50052")
	if err := srv.Serve(lis); err != nil { log.Fatal(err) }
}
