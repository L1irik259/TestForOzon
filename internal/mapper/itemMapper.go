package mapper

import (
	domain "github.com/L1irik259/TestForOzon/internal/domain"
	pb "github.com/L1irik259/TestForOzon/internal/transport/proto/github.com/L1irik259/TestForOzon/transport/genetation/go/v1"
)

// Структура для маппинга данных из доменной модели в protobuf модель для передачи по gRPC
func MapToProto(item domain.Item) *pb.Item {
	return &pb.Item{
		Id:        item.ItemStaticData.ID,
		NumCode:   item.ItemStaticData.NumCode,
		CharCode:  item.ItemStaticData.CharCode,
		Name:      item.ItemStaticData.Name,
		Nominal:   int32(item.ItemDynamicData.Nominal),
		Value:     item.ItemDynamicData.Value,
		VunitRate: item.ItemDynamicData.VunitRate,
		Date:      item.ItemDynamicData.Date.Format("2006-01-02"),
	}
}
