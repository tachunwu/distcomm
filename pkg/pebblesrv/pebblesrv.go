package pebblesrv

import (
	"context"
	"log"
	"net"

	"github.com/cockroachdb/pebble"
	pebblepb "github.com/tachunwu/distcomm/pkg/proto/pebble"
	rpc "google.golang.org/grpc"
)

type PebbleServer interface {
	Init(addr string, dirname string)
	Serve()
	Close()
}

func NewPebbleServer() PebbleServer {
	return &pebbleServer{}
}

type pebbleServer struct {
	pebblepb.PebbleServiceServer
	addr       string
	rpc        *rpc.Server
	dialCtx    context.Context
	dialCancel func()
	db         *pebble.DB
}

func (srv *pebbleServer) Init(addr string, dirname string) {
	srv.addr = addr
	srv.dialCtx, srv.dialCancel = context.WithCancel(context.Background())
	srv.rpc = rpc.NewServer()
	pebblepb.RegisterPebbleServiceServer(srv.rpc, srv)

	db, err := pebble.Open(dirname, &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	srv.db = db
}

func (srv *pebbleServer) Serve() {

	var lis net.Listener

	lis, err := net.Listen("tcp", srv.addr)
	if err != nil {
		log.Println(err)
	}

	if err := srv.rpc.Serve(lis); err != nil {
		switch err {
		case rpc.ErrServerStopped:
		default:
			log.Fatal(err)
		}
	}
}

func (srv *pebbleServer) Close() {
	srv.rpc.Stop()
	srv.dialCancel()
	if err := srv.db.Close(); err != nil {
		log.Println(err)
	}
}

// Service logic

func (srv *pebbleServer) Set(ctx context.Context, req *pebblepb.SetRequest) (*pebblepb.SetResponse, error) {
	err := srv.db.Set(
		req.GetKey(),
		req.GetValue(),
		pebble.Sync,
	)

	if err != nil {
		log.Println(err)
		return &pebblepb.SetResponse{
			Status: "set operation fail",
		}, err
	}
	log.Println(
		"Set: ",
		string(req.GetKey()),
		string(req.GetValue()),
	)
	return &pebblepb.SetResponse{
		Status: "set operation success",
	}, nil
}

func (srv *pebbleServer) Get(ctx context.Context, req *pebblepb.GetRequest) (*pebblepb.GetResponse, error) {
	value, closer, err := srv.db.Get(req.GetKey())
	if err != nil {
		log.Println(err)
		return &pebblepb.GetResponse{
			Status: "get operation fail",
		}, nil
	}

	if err := closer.Close(); err != nil {
		log.Println(err)
		return &pebblepb.GetResponse{
			Status: "get operation fail",
		}, nil
	}

	log.Println(
		"Get: ",
		string(req.GetKey()),
		string(value),
	)

	return &pebblepb.GetResponse{
		Value:  value,
		Status: "get operation success",
	}, nil
}

func (srv *pebbleServer) Delete(ctx context.Context, req *pebblepb.DeleteRequest) (*pebblepb.DeleteResponse, error) {
	err := srv.db.Delete(req.GetKey(), pebble.Sync)
	if err != nil {
		log.Println(err)
		return &pebblepb.DeleteResponse{
			Status: "delete operation fail",
		}, err
	}

	log.Println(
		"Delete: ",
		string(req.GetKey()),
	)
	return &pebblepb.DeleteResponse{
		Status: "delete operation success",
	}, nil
}

func (srv *pebbleServer) Scan(ctx context.Context, req *pebblepb.ScanRequest) (*pebblepb.ScanResponse, error) {
	values := [][]byte{}
	iter := srv.db.NewIter(scanIterOptions(req.GetStartKey(), req.GetEndKey()))
	for iter.First(); iter.Valid(); iter.Next() {
		values = append(values, iter.Value())
	}
	if err := iter.Close(); err != nil {
		log.Println(err)
		return &pebblepb.ScanResponse{
			Status: "scan operation fail",
		}, err
	}

	log.Println(
		"Scan: ",
		string(req.StartKey),
		string(req.EndKey),
	)

	return &pebblepb.ScanResponse{
		Status: "scan operation success",
		Values: values,
	}, nil
}

func scanIterOptions(start []byte, end []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	}
}
