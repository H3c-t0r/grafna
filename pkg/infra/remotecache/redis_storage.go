package remotecache

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/grafana/grafana/pkg/setting"
)

const redisCacheType = "redis"

type redisStorage struct {
	c     *redis.Client
	codec codec
}

// parseRedisConnStr parses k=v pairs in csv and builds a redis Options object
func parseRedisConnStr(connStr string) (*redis.Options, error) {
	keyValueCSV := strings.Split(connStr, ",")
	options := &redis.Options{Network: "tcp"}
	setTLSIsTrue := false
	for _, rawKeyValue := range keyValueCSV {
		keyValueTuple := strings.SplitN(rawKeyValue, "=", 2)
		if len(keyValueTuple) != 2 {
			if strings.HasPrefix(rawKeyValue, "password") {
				// don't log the password
				rawKeyValue = "password" + setting.RedactedPassword
			}
			return nil, fmt.Errorf("incorrect redis connection string format detected for '%v', format is key=value,key=value", rawKeyValue)
		}
		connKey := keyValueTuple[0]
		connVal := keyValueTuple[1]
		switch connKey {
		case "addr":
			options.Addr = connVal
		case "password":
			options.Password = connVal
		case "db":
			i, err := strconv.Atoi(connVal)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", "value for db in redis connection string must be a number", err)
			}
			options.DB = i
		case "pool_size":
			i, err := strconv.Atoi(connVal)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", "value for pool_size in redis connection string must be a number", err)
			}
			options.PoolSize = i
		case "ssl":
			if connVal != "true" && connVal != "false" && connVal != "insecure" {
				return nil, fmt.Errorf("ssl must be set to 'true', 'false', or 'insecure' when present")
			}
			if connVal == "true" {
				setTLSIsTrue = true // Needs addr already parsed, so set later
			}
			if connVal == "insecure" {
				options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			}
		default:
			return nil, fmt.Errorf("unrecognized option '%v' in redis connection string", connKey)
		}
	}
	if setTLSIsTrue {
		// Get hostname from the Addr property and set it on the configuration for TLS
		sp := strings.Split(options.Addr, ":")
		if len(sp) < 1 {
			return nil, fmt.Errorf("unable to get hostname from the addr field, expected host:port, got '%v'", options.Addr)
		}
		options.TLSConfig = &tls.Config{ServerName: sp[0]}
	}
	return options, nil
}

func newRedisStorage(opts *setting.RemoteCacheOptions, codec codec) (*redisStorage, error) {
	opt, err := parseRedisConnStr(opts.ConnStr)
	if err != nil {
		return nil, err
	}
	return &redisStorage{c: redis.NewClient(opt), codec: codec}, nil
}

// Set sets value to given key in session.
func (s *redisStorage) Set(ctx context.Context, key string, val interface{}, expires time.Duration) error {
	item := &cachedItem{Val: val}
	value, err := s.codec.Encode(ctx, item)
	if err != nil {
		return err
	}

	return s.SetByteArray(ctx, key, value, expires)
}

// Set sets value to a given key
func (s *redisStorage) SetByteArray(ctx context.Context, key string, data []byte, expires time.Duration) error {
	status := s.c.Set(ctx, key, data, expires)
	return status.Err()
}

// Get gets value by given key in session.
func (s *redisStorage) Get(ctx context.Context, key string) (interface{}, error) {
	v, err := s.GetByteArray(ctx, key)

	if err.Error() == "EOF" {
		return nil, ErrCacheItemNotFound
	}

	if err != nil {
		return nil, err
	}

	item := &cachedItem{}
	err = s.codec.Decode(ctx, v, item)

	if err == nil {
		return item.Val, nil
	}

	return nil, err
}

// GetByteArray returns the value as byte array
func (s *redisStorage) GetByteArray(ctx context.Context, key string) ([]byte, error) {
	return s.c.Get(ctx, key).Bytes()
}

// Delete delete a key from session.
func (s *redisStorage) Delete(ctx context.Context, key string) error {
	cmd := s.c.Del(ctx, key)
	return cmd.Err()
}
