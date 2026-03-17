package transport

import (
	"context"
	"time"

	mapper "github.com/L1irik259/TestForOzon/internal/mapper"
	service "github.com/L1irik259/TestForOzon/internal/service"
	pb "github.com/L1irik259/TestForOzon/internal/transport/proto/github.com/L1irik259/TestForOzon/transport/genetation/go/v1"
)

type Server struct {
	ItemService service.ItemService
	pb.UnimplementedOzonServiceServer
}

func (s *Server) GetItem(ctx context.Context, req *pb.ItemRequest) (*pb.ItemResponse, error) {
	date, errDate := time.Parse("02/01/2006", req.Date)
	if errDate != nil {
		return nil, errDate
	}

	items, err := s.ItemService.FindAllItemsByDate(date)
	if err != nil {
		return nil, err
	}

	var pbItems []*pb.Item
	for _, item := range items {
		pbItems = append(pbItems, mapper.MapToProto(*item))
	}

	return &pb.ItemResponse{
		Items: pbItems,
	}, nil
}
