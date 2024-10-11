package controller

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/t3201v/seat-arrangement/gen/cinema"
	"github.com/t3201v/seat-arrangement/internal/model"
	"github.com/t3201v/seat-arrangement/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ICinema interface {
	cinema.CinemaServiceServer
}

type Cinema struct {
	cinema.UnimplementedCinemaServiceServer
	svc    service.ICinema
	logger *log.Logger
}

func (c *Cinema) ConfigureCinema(ctx context.Context, request *cinema.ConfigureCinemaRequest) (*cinema.ConfigureCinemaResponse, error) {
	id, err := c.svc.ConfigureCinema(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &cinema.ConfigureCinemaResponse{
		Id: id,
	}, nil
}

func (c *Cinema) UpdateCinemaConfig(ctx context.Context, request *cinema.UpdateCinemaConfigRequest) (*cinema.SuccessResponse, error) {
	err := c.svc.UpdateCinemaConfig(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &cinema.SuccessResponse{
		Success: true,
	}, nil
}

func (c *Cinema) GetAvailableSeats(ctx context.Context, request *cinema.GetAvailableSeatsRequest) (*cinema.GetAvailableSeatsResponse, error) {
	data, grid, err := c.svc.GetAvailableSeats(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result, err := new(model.Cinema).ToPbSeats(data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &cinema.GetAvailableSeatsResponse{
		AvailableSeats: result,
		Grid:           grid,
	}, nil
}

func (c *Cinema) ReserveSeats(ctx context.Context, request *cinema.ReserveSeatsRequest) (*cinema.SuccessResponse, error) {
	err := c.svc.ReserveSeats(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &cinema.SuccessResponse{Success: true}, nil
}

func (c *Cinema) CancelSeats(ctx context.Context, request *cinema.CancelSeatsRequest) (*cinema.SuccessResponse, error) {
	err := c.svc.CancelSeats(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &cinema.SuccessResponse{Success: true}, nil
}

func NewCinema(l *log.Logger, svc service.ICinema) ICinema {
	return &Cinema{
		logger: l,
		svc:    svc,
	}
}
