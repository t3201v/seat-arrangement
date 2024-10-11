package repository

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"github.com/t3201v/seat-arrangement/internal/model"
)

// memory storage
type ICinema interface {
	GetCinema(id string) (*model.Cinema, error)
	InsertCinema(*model.Cinema) (string, error)
	UpdateCinema(string, *model.Cinema) error
}

type Cinema struct {
	mu      sync.RWMutex
	logger  *log.Logger
	counter *int64
	cinemas map[int64]*model.Cinema
}

func (c *Cinema) GetCinema(id string) (*model.Cinema, error) {
	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.cinemas[_id]; !ok {
		return nil, nil
	}
	return c.cinemas[_id].Clone(), nil
}

func (c *Cinema) InsertCinema(entity *model.Cinema) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	counter := atomic.LoadInt64(c.counter)
	c.cinemas[counter] = entity
	atomic.AddInt64(c.counter, 1)
	return strconv.FormatInt(counter, 10), nil
}

func (c *Cinema) UpdateCinema(id string, entity *model.Cinema) error {
	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("invalid id")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cinemas[_id] = entity
	return nil
}

func NewCinema(l *log.Logger) ICinema {
	return &Cinema{
		mu:      sync.RWMutex{},
		logger:  l,
		counter: new(int64),
		cinemas: make(map[int64]*model.Cinema),
	}
}
