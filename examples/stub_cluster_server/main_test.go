package main

import (
	"context"
	"net"
	"testing"

	kotgv1 "github.com/kubilitics/kotg-schema/gen/go/kotg/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStubClusterServer_GetCluster(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil { t.Fatal(err) }
	defer lis.Close()

	srv := grpc.NewServer()
	kotgv1.RegisterClusterReadServer(srv, &stubClusterRead{})
	go srv.Serve(lis)
	defer srv.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { t.Fatal(err) }
	defer conn.Close()
	client := kotgv1.NewClusterReadClient(conn)

	got, err := client.GetCluster(context.Background(), &kotgv1.GetClusterRequest{ClusterId: "test"})
	if err != nil { t.Fatal(err) }
	if got.GetClusterId() != "test" {
		t.Fatalf("expected cluster_id=test, got %q", got.GetClusterId())
	}
}
