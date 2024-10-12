package service

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/t3201v/seat-arrangement/gen/cinema"
	"github.com/t3201v/seat-arrangement/internal/model"
	"github.com/t3201v/seat-arrangement/repository"
)

type ICinema interface {
	ConfigureCinema(ctx context.Context, request *cinema.ConfigureCinemaRequest) (string, error)
	UpdateCinemaConfig(ctx context.Context, request *cinema.UpdateCinemaConfigRequest) error
	GetAvailableSeats(ctx context.Context, request *cinema.GetAvailableSeatsRequest) ([][][]int, string, error)
	ReserveSeats(ctx context.Context, request *cinema.ReserveSeatsRequest) error
	CancelSeats(ctx context.Context, request *cinema.CancelSeatsRequest) error
}

type Cinema struct {
	logger *log.Logger
	repo   repository.ICinema
}

func (c *Cinema) ConfigureCinema(ctx context.Context, request *cinema.ConfigureCinemaRequest) (string, error) {
	entity := model.NewCinema(c.logger, int(request.Rows), int(request.Columns), int(request.MinDistance))
	id, err := c.repo.InsertCinema(entity)
	if err != nil {
		c.logger.Error(err)
		return "", err
	}
	return id, nil
}

func (c *Cinema) UpdateCinemaConfig(ctx context.Context, request *cinema.UpdateCinemaConfigRequest) error {
	entity, err := c.repo.GetCinema(request.Id)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	if entity == nil {
		return fmt.Errorf("not found id %s", request.Id)
	}

	entity.UpdateConfig(int(request.Rows), int(request.Columns), int(request.MinDistance))
	err = c.repo.UpdateCinema(request.Id, entity)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	return nil
}

func (c *Cinema) GetAvailableSeats(ctx context.Context, request *cinema.GetAvailableSeatsRequest) ([][][]int, string, error) {
	entity, err := c.repo.GetCinema(request.Id)
	if err != nil {
		c.logger.Error(err)
		return nil, "", err
	}
	if entity == nil {
		return nil, "", fmt.Errorf("not found id %s", request.Id)
	}
	return entity.ListAvailableSeatsGrouped(), entity.String(), nil
}

func (c *Cinema) ReserveSeats(ctx context.Context, request *cinema.ReserveSeatsRequest) error {
	entity, err := c.repo.GetCinema(request.Id)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	if entity == nil {
		return fmt.Errorf("not found id %s", request.Id)
	}

	seats := make([][]int, 0)
	for _, seat := range request.SeatCoords {
		if seat == nil {
			return errors.New("malformed seat data")
		}
		seats = append(seats, []int{int(seat.Row), int(seat.Column)})
	}
	err = entity.ReserveSeats(seats, request.GroupName)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	err = c.repo.UpdateCinema(request.Id, entity)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	return nil
}

func (c *Cinema) CancelSeats(ctx context.Context, request *cinema.CancelSeatsRequest) error {
	entity, err := c.repo.GetCinema(request.Id)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	if entity == nil {
		return fmt.Errorf("not found id %s", request.Id)
	}

	seats := make([][]int, 0)
	for _, seat := range request.SeatCoords {
		if seat == nil {
			return errors.New("malformed seat data")
		}
		seats = append(seats, []int{int(seat.Row), int(seat.Column)})
	}
	err = entity.CancelSeats(seats)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	err = c.repo.UpdateCinema(request.Id, entity)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	return nil
}

func NewCinema(l *log.Logger, repo repository.ICinema) ICinema {
	return &Cinema{
		logger: l,
		repo:   repo,
	}
}
