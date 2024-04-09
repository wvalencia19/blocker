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

func (n *Node) bootstrapNetwork(addrs []string) error {
	n.logger.Debugw("dialing remote node", "we", n.listedAddr, "remote", addrs)
	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}
		c, v, err := n.dialRemote(addr)
		if err != nil {
			return err
		}
		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.listedAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)
	n.logger.Infow("node started...", "port", listenAddr)

	// bootstrap the network with a list of already know nodes
	// in the network

	if len(bootstrapNodes) > 0 {
		go n.bootstrapNetwork(bootstrapNodes)

	}
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
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) canConnectWith(addr string) bool {
	if n.listedAddr == addr {
		return false
	}
	connectedPeers := n.getPeerList()

	for _, connectedAdr := range connectedPeers {
		if addr == connectedAdr {
			return false
		}
	}
	return true
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := []string{}

	for _, version := range n.peers {
		peers = append(peers, version.ListedAddr)
	}

	return peers
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// handle the logic where we decide to accept pr drop
	// the incoming node connection

	n.peers[c] = v

	if len(v.PeerList) > 0 {
		go n.bootstrapNetwork(v.PeerList)
	}
	n.logger.Infof("we[%s] new peer connected remote(%s) -height (%d)", n.listedAddr, v.ListedAddr, v.Height)
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) dialRemote(addr string) (proto.NodeClient, *proto.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.HandShake(context.Background(), n.getVersion())
	if err != nil {
		n.logger.Errorf("handshake error", err)
		return nil, nil, err
	}
	return c, v, nil
}
