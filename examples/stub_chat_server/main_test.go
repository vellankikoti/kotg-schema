package main

import (
	"context"
	"net"
	"testing"

	kotgv1 "github.com/kubilitics/kotg-schema/gen/go/kotg/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStubChatServer_Capabilities(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil { t.Fatal(err) }
	defer lis.Close()

	srv := grpc.NewServer()
	kotgv1.RegisterAIControlServer(srv, &stubAIControl{})
	go srv.Serve(lis)
	defer srv.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { t.Fatal(err) }
	defer conn.Close()
	client := kotgv1.NewAIControlClient(conn)

	got, err := client.Capabilities(context.Background(), &kotgv1.Empty{})
	if err != nil { t.Fatal(err) }
	if got.GetSchemaVersion() != "1.0.0" {
		t.Fatalf("expected schema_version=1.0.0, got %q", got.GetSchemaVersion())
	}
	if !got.GetSupportsUndo() {
		t.Fatal("supports_undo should be true")
	}
}
