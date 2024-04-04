package node

import (
	"context"
	"net"
	"sync"

	"github.com/wvalencia19/blocker/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	peerLock   sync.RWMutex
	listedAddr string
	logger     *zap.SugaredLogger
	peers      map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	logger, _ := zap.NewProduction()
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "blocker-0.1",
		logger:  logger.Sugar(),
	}
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Infof("[%s] new peer connected (%s) -height (%d)", n.listedAddr, v.ListedAddr, v.Height)
	n.peers[c] = v
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.HandShake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Errorf("handshake error", err)
			continue
		}
		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) Start(listenAddr string) error {
	n.listedAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)
	n.logger.Infow("node started...", "port", listenAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)

	n.logger.Infof("received tx from", peer)
	return &proto.Ack{}, nil
}

func (n *Node) HandShake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListedAddr)

	if err != nil {
		return nil, err
	}
	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "blocker-0.1",
		Height:     0,
		ListedAddr: n.listedAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
