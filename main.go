package main

import (
	"context"
	"log"

	"github.com/wvalencia19/blocker/node"
	"github.com/wvalencia19/blocker/proto"
	"google.golang.org/grpc"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})

	select {}

}

func makeNode(listedAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listedAddr)
	if len(bootstrapNodes) > 0 {

		if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
			log.Fatal(err)
		}

	}

	return n
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
