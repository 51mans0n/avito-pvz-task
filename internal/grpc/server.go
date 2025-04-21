package grpcserver

import (
	"context"

	"github.com/51mans0n/avito-pvz-task/internal/db"
	pvz_v1 "github.com/51mans0n/avito-pvz-task/pkg/proto/pvz/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pvz_v1.UnimplementedPVZServiceServer
	repo db.Repository
}

func New(repo db.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetPVZList(ctx context.Context, _ *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	rows, err := s.repo.GetPVZListWithFilter(ctx, nil, nil, 1, 1000)
	if err != nil {
		return nil, err
	}

	resp := &pvz_v1.GetPVZListResponse{}
	for _, r := range rows {
		resp.Pvzs = append(resp.Pvzs, &pvz_v1.PVZ{
			Id:               r.PVZ.ID,
			City:             r.PVZ.City,
			RegistrationDate: timestamppb.New(r.PVZ.RegistrationDate),
		})
	}
	return resp, nil
}
