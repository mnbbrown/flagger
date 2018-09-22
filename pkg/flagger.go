package flagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

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
	Name          string
	Type          flagType `json:"type"`
	InternalValue int      `json:"value"`
	Namespace     string   `json:"namespace"`
	Tags          []string `json:"tags"`
}

// MarshalBinary implements encoding support for redis
func (f *Flag) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

// UnmarshalBinary implements encoding support for redis
func (f *Flag) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
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

// Flagger is implemented by each of the storage backends
type Flagger interface {
	SaveFlag(flag *Flag) error
	GetFlag(name string) (*Flag, error)
	GetFlagWithTags(name string, tags []string) (*Flag, error)
	ListFlags() ([]*Flag, error)
}

// RedisFlagger is a redis backed flagger
type RedisFlagger struct {
	client    *redis.Client
	namespace string
}

// NewRedisFlagger creates a new flagger backed by redis
func NewRedisFlagger(host string, db int) (Flagger, error) {
	client := redis.NewClient(&redis.Options{
		Addr: host,
		DB:   db,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return &RedisFlagger{client: client, namespace: "flagger"}, nil
}

func (rf *RedisFlagger) getKeyName(key ...string) string {
	path := strings.Join(key, ":")
	return fmt.Sprintf("%s:%s", rf.namespace, path)
}

// SaveFlag saves a flag to redis
func (rf *RedisFlagger) SaveFlag(flag *Flag) error {
	for _, tag := range flag.Tags {
		if err := rf.client.SAdd(rf.getKeyName("TAGS", tag), flag.Name).Err(); err != nil {
			return err
		}
	}
	if err := rf.client.Set(rf.getKeyName("IDS", flag.Name), flag, 0).Err(); err != nil {
		return err
	}
	return nil
}

// ErrFlagNotFound means the flag was not found
var ErrFlagNotFound = errors.New("Flag not found")

func flagInResults(flag string, tags []string) bool {
	for _, f := range tags {
		if f == flag {
			return true
		}
	}
	return false
}

// GetFlag loads a flag without tags
func (rf *RedisFlagger) GetFlag(name string) (*Flag, error) {
	return rf.GetFlagWithTags(name, []string{})
}

// GetFlagWithTags loads a flag from a redis client
func (rf *RedisFlagger) GetFlagWithTags(name string, tags []string) (*Flag, error) {

	var flagsWithTags []string
	if len(tags) > 0 {
		prefixed := []string{}
		for _, tag := range tags {
			prefixed = append(prefixed, rf.getKeyName("TAGS", tag))
		}

		res := rf.client.SInter(prefixed...)
		if res.Err() != nil {
			return nil, res.Err()
		}
		var err error
		flagsWithTags, err = res.Result()
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(flagsWithTags, tags)

	if len(tags) > 0 && !flagInResults(name, flagsWithTags) {
		return nil, ErrFlagNotFound
	}

	f := &Flag{}
	if err := rf.client.Get(rf.getKeyName("IDS", name)).Scan(f); err != nil {
		return nil, err
	}
	return f, nil
}

// ListFlags returns a list of flags grouped by name and environment
func (rf *RedisFlagger) ListFlags() (flags []*Flag, err error) {
	return nil, nil
}
