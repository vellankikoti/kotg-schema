// Stub ClusterRead + ClusterAction implementation. Proves consumers can
// implement the kotg-schema service interfaces. Returns hardcoded data.
//
// Run: go run ./examples/stub_cluster_server
package main

import (
	"context"
	"errors"
	"log"
	"net"

	kotgv1 "github.com/vellankikoti/kotg-schema/gen/go/kotg/v1"
	"google.golang.org/grpc"
)

type stubClusterRead struct {
	kotgv1.UnimplementedClusterReadServer
}

func (s *stubClusterRead) GetCluster(_ context.Context, req *kotgv1.GetClusterRequest) (*kotgv1.Cluster, error) {
	if req.GetClusterId() == "" {
		return nil, errors.New("cluster_id required")
	}
	return &kotgv1.Cluster{
		ClusterId:      req.GetClusterId(),
		Name:           "stub-cluster",
		Distribution:   "in-cluster",
		K8SVersion:     "v1.33.0",
		NodeCount:      1,
		NamespaceCount: 5,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil { log.Fatal(err) }
	srv := grpc.NewServer()
	kotgv1.RegisterClusterReadServer(srv, &stubClusterRead{})
	log.Println("stub_cluster_server listening on 127.0.0.1:50051")
	if err := srv.Serve(lis); err != nil { log.Fatal(err) }
}
