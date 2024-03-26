package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/wvalencia19/blocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version  string
	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "blocker-0.1",
	}
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	fmt.Printf("new peer connected (%s) -height (%d)\n", v.ListedAddr, v.Height)
	n.peers[c] = v
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) Start(listenAddr string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)
	fmt.Printf("node running on port :%s", listenAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from", peer)
	return &proto.Ack{}, nil
}

func (n *Node) HandShake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}

	peer, _ := peer.FromContext(ctx)
	c, err := makeNodeClient(v.ListedAddr)

	if err != nil {
		return nil, err
	}
	n.addPeer(c, v)

	fmt.Printf("received version from %s: %+v\n", v, peer.Addr)
	return ourVersion, nil
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
