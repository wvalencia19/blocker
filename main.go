package main

import (
	"context"
	"log"
	"time"

	"github.com/wvalencia19/blocker/crypto"
	"github.com/wvalencia19/blocker/node"
	"github.com/wvalencia19/blocker/proto"
	"github.com/wvalencia19/blocker/util"
	"google.golang.org/grpc"
)

func main() {
	makeNode(":3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"}, false)

	time.Sleep(time.Second)
	makeNode(":5001", []string{":4000"}, false)

	for {
		time.Sleep(time.Second * 2)
		makeTransaction()
	}
	select {}

}

func makeNode(listedAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "Blocker-1",
		ListenAddr: listedAddr,
	}
	if isValidator {
		cfg.PrivateKey = crypto.GeneratePrivatekey()
	}
	n := node.NewNode(cfg)
	go n.Start(listedAddr, bootstrapNodes)

	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)
	privKey := crypto.GeneratePrivatekey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
