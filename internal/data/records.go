package data

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skraio/mini-godis/internal/validator"
)

type Record struct {
	Key   string `json:"key"`
	Value *int64 `json:"value"`
}

func ValidateRecord(v *validator.Validator, record *Record) {
	v.Check(record.Key != "", "key", "must be provided")
	v.Check(len(record.Key) <= 28, "key", "must not be more than 28 bytes long")

	v.Check(record.Value != nil, "value", "must be provided")
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

type RecordModel struct {
	RedisDB *redis.Client
}

func (r RecordModel) Insert(record *Record) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.RedisDB.Set(ctx, record.Key, *record.Value, 0).Err()
}

func (r RecordModel) Get(key string) (*Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	value, err := r.RedisDB.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	record := &Record{
		Key:   key,
		Value: &intValue,
	}

	return record, nil
}

func (r RecordModel) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := r.RedisDB.Del(ctx, key).Result()
	if err != nil {
		return nil
	}

	if result == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type Operation int

const (
	Increase Operation = iota
	Decrease
)

func (r RecordModel) Update(record *Record, op Operation, term int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var err error

	switch op {
	case Increase:
		_, err = r.RedisDB.IncrBy(ctx, record.Key, term).Result()
	case Decrease:
		_, err = r.RedisDB.DecrBy(ctx, record.Key, term).Result()
	}

	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}
