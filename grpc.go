package istio

import (
	"net"

	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer 服务端
type GrpcServer struct {
	ServiceName string
	address     string
	server      *grpc.Server
	Log         log.Logger
}

// GetSrv 获取rpc server
func (s *GrpcServer) GetSrv() *grpc.Server {
	return s.server
}

// Start 启动
func (s *GrpcServer) Start() {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		s.Log.Errorf("%s grpc server failed to listen: %v", s.ServiceName, err)
		return
	}
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.server)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.Log.Errorf("%s grpc server failed to serve: %v", s.ServiceName, err)
		}
	}()
}
func (s *GrpcServer) Stop() {
	if s.server == nil {
		s.Log.Warningf("stop %s grpc server is nil", s.ServiceName)
		return
	}
	s.server.Stop()
}

// NewGrpcServer 创建Grpc服务端
func NewGrpcServer(name, address string, interceptor ...grpc.UnaryServerInterceptor) *GrpcServer {
	var opts []grpc.ServerOption
	for _, v := range interceptor {
		opts = append(opts, grpc.UnaryInterceptor(v))
	}
	server := grpc.NewServer(opts...)
	if logger.Default() == nil {
		logger.Init()
	}
	return &GrpcServer{
		ServiceName: name,
		server:      server,
		address:     address,
		Log:         logger.Default(),
	}
}

// Client grpc客户端
type GrpcClient struct {
	ServiceName string
	conn        *grpc.ClientConn // 连接
	Log         log.Logger
}

// GetConn 获取客户端连接
func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭
func (c *GrpcClient) Close() {
	if c.conn == nil {
		c.Log.Warningf("close %s grpc client is nil", c.ServiceName)
		return
	}
	err := c.conn.Close()
	if err != nil {
		c.Log.Errorf("close %s grpc client: %v", err)
		return
	}
}

// NewGrpcClient 创建Grpc客户端
func NewGrpcClient(name, serverAddress string, interceptor ...grpc.UnaryClientInterceptor) *GrpcClient {
	var opts []grpc.DialOption
	for _, v := range interceptor {
		opts = append(opts, grpc.WithUnaryInterceptor(v))
	}
	if logger.Default() == nil {
		logger.Init()
	}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		logger.Errorf("%s grpc client dial error: %v", name, err)
	}
	return &GrpcClient{
		ServiceName: name,
		conn:        conn,
		Log:         logger.Default(),
	}
}
