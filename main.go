package main

import (
	"context"
	"log"
	"time"

	"github.com/wvalencia19/blocker/node"
	"github.com/wvalencia19/blocker/proto"
	"google.golang.org/grpc"
)

func main() {
	node := node.NewNode()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeTransaction()
		}

	}()

	log.Fatal(node.Start(":3000"))
}

func makeNode(listedAddr string) *node.Node {
	n := node.NewNode()

	go n.Start(listedAddr)
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version:    "blocker-0.1",
		Height:     1,
		ListedAddr: ":4000",
	}

	_, err = c.HandShake(context.TODO(), version)
	if err != nil {
		log.Fatal(err)
	}
}
