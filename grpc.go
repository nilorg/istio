package istio

import (
	"fmt"
	"net"

	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer 服务端
type GrpcServer struct {
	serviceName string
	address     string
	server      *grpc.Server
	Log         log.Logger
}

// GetSrv 获取rpc server
func (s *GrpcServer) GetSrv() *grpc.Server {
	return s.server
}

func (s *GrpcServer) register() {
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.server)
}

// Run ...
func (s *GrpcServer) Run() {
	s.register()

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		s.Log.Errorf("%s grpc server failed to listen: %v", s.serviceName, err)
		return
	}
	if err := s.server.Serve(lis); err != nil {
		s.Log.Errorf("%s grpc server failed to serve: %v", s.serviceName, err)
	}
}

// Start 启动
func (s *GrpcServer) Start() {
	go func() {
		s.Run()
	}()
}

// Stop ...
func (s *GrpcServer) Stop() {
	if s.server == nil {
		s.Log.Warningf("stop %s grpc server is nil", s.serviceName)
		return
	}
	s.server.Stop()
}

// NewGrpcServer 创建Grpc服务端
func NewGrpcServer(name string, address string, interceptor ...grpc.UnaryServerInterceptor) *GrpcServer {
	var opts []grpc.ServerOption
	for _, v := range interceptor {
		opts = append(opts, grpc.UnaryInterceptor(v))
	}
	server := grpc.NewServer(opts...)
	if logger.Default() == nil {
		logger.Init()
	}
	return &GrpcServer{
		serviceName: name,
		server:      server,
		address:     address,
		Log:         logger.Default(),
	}
}

// GrpcClient grpc客户端
type GrpcClient struct {
	serviceName string
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
		c.Log.Warningf("close %s grpc client is nil", c.serviceName)
		return
	}
	err := c.conn.Close()
	if err != nil {
		c.Log.Errorf("close %s grpc client: %v", err)
		return
	}
}

// NewGrpcClient 创建Grpc客户端
func NewGrpcClient(serviceName string, port int, interceptor ...grpc.UnaryClientInterceptor) *GrpcClient {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	for _, v := range interceptor {
		opts = append(opts, grpc.WithUnaryInterceptor(v))
	}
	if logger.Default() == nil {
		logger.Init()
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serviceName, port), opts...)
	if err != nil {
		logger.Errorf("%s grpc client dial error: %v", serviceName, err)
	}
	return &GrpcClient{
		serviceName: serviceName,
		conn:        conn,
		Log:         logger.Default(),
	}
}
