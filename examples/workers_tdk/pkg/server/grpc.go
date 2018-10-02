package server

import (
	"workers/examples/workers_tdk/pkg/pb"
	"workers/examples/workers_tdk/pkg/usecase"

	"github.com/tokopedia/tdk/go/app/grpc"
	"golang.org/x/net/context"
)

func NewGrpcServer() *grpc.GrpcService {
	svc := grpc.New(&grpc.Config{
		Address: ":5000",
	})

	pb.RegisterAppServer(svc.Server(), &OrderServer{})
	svc.RegisterGeniHttp(&OrderServer{})

	return svc
}

type OrderServer struct {
}

// lets create order via gRPC
func (*OrderServer) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderReply, error) {
	order := usecase.Order{
		UserID:    int(req.UserId),
		Quantity:  int(req.Quantity),
		ProductID: int(req.ProductId),
	}

	invoice, err := orderUsecase.PutNewOrder(order)
	if err != nil {
		return nil, err
	}

	return &pb.OrderReply{
		Invoice: invoice,
	}, nil
}
