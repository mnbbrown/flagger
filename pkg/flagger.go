package flagger

import (
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/go-redis/redis"
)

var (
	globalDefault = &Flag{
		Type:          BOOL,
		InternalValue: 1,
	}
)

type flagType string

// flagTypes are types of flags
const (
	BOOL    flagType = "BOOL"    // BOOL is a simple boolean flag type
	PERCENT flagType = "PERCENT" // PERCENT is a boolean that is true ${PERCENT}% of the time
)

// Flag is a flag ;)
type Flag struct {
	Type          flagType `json:"type"`
	InternalValue int      `json:"value"`
}

// MarshalBinary implements encoding support for redis
func (f *Flag) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

// UnmarshalBinary implements encoding support for redis
func (f *Flag) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	return nil
}

// Value returns the calculated flag value
func (f *Flag) Value() bool {
	switch f.Type {
	case BOOL:
		return f.InternalValue != 0
	case PERCENT:
		return (rand.Float32() * 100) < float32(f.InternalValue)
	}
	return globalDefault.Value()
}

// SaveFlag saves a flag to redis
func SaveFlag(redisClient *redis.Client, name, environment string, flag *Flag) error {
	if err := redisClient.HSet(name, environment, flag).Err(); err != nil {
		return err
	}
	return nil
}

// ErrFlagNotFound means the flag was not found
var ErrFlagNotFound = errors.New("Flag not found")

// GetFlag loads a flag from a redis client
func GetFlag(redisClient *redis.Client, name, environment string) (*Flag, error) {
	var result *redis.StringCmd
	if result = redisClient.HGet(name, environment); result.Err() == redis.Nil {
		result = redisClient.HGet(name, "default")
	}

	if result.Err() == redis.Nil {
		return nil, ErrFlagNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	f := &Flag{}
	err := result.Scan(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
