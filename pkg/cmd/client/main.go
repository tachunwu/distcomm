package main

import (
	"context"
	"log"

	pebblepb "github.com/tachunwu/distcomm/pkg/proto/pebble"
	rpc "google.golang.org/grpc"
)

func main() {
	conn, err := rpc.Dial("localhost:30000", rpc.WithInsecure())
	if err != nil {
		switch err {
		case context.Canceled:
			return
		default:
			log.Println("error when dialing: ", err)
		}
	}

	defer conn.Close()
	client := pebblepb.NewPebbleServiceClient(conn)

	// Test
	Set(client, "/table/primary/0", "0")
	Set(client, "/table/primary/1", "1")
	Set(client, "/table/primary/2", "2")
	Set(client, "/table/primary/3", "3")
	Get(client, "/table/primary/0")
	Delete(client, "/table/primary/3")
	Scan(client, "/table/primary/0", "/table/primary/4")

}

func Set(client pebblepb.PebbleServiceClient, key string, value string) {
	req := &pebblepb.SetRequest{
		Key:   []byte(key),
		Value: []byte(value),
	}
	res, err := client.Set(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}

func Get(client pebblepb.PebbleServiceClient, key string) {
	req := &pebblepb.GetRequest{
		Key: []byte(key),
	}
	res, err := client.Get(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}

func Delete(client pebblepb.PebbleServiceClient, key string) {
	req := &pebblepb.DeleteRequest{
		Key: []byte(key),
	}
	res, err := client.Delete(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}

func Scan(client pebblepb.PebbleServiceClient, start string, end string) {
	req := &pebblepb.ScanRequest{
		StartKey: []byte(start),
		EndKey:   []byte(end),
	}
	res, err := client.Scan(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
