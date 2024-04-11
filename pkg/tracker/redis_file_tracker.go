package tracker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFileTracker struct {
	r          *redis.ClusterClient
	expiration int
}

func NewRedisFileTracker(r *redis.ClusterClient, invalidationTime int) *RedisFileTracker {
	return &RedisFileTracker{r, invalidationTime}
}

func (t *RedisFileTracker) Create(status *FileStatus) error {
	duration := time.Duration(0)

	if strings.HasPrefix(status.Status, "fail:") {
		duration = time.Duration(t.expiration) * time.Second
	}

	status.UpdatedAt = time.Now()

	progressJson, err := json.Marshal(status)
	if err != nil {
		return err
	}

	if err := t.r.Set(context.Background(), status.Key, progressJson, duration).Err(); err != nil {
		return err
	}

	return nil
}

func (t *RedisFileTracker) Delete(key string) error {
	return t.r.Del(context.Background(), key).Err()
}

func (t *RedisFileTracker) Clear() error {
	return t.r.FlushAll(context.Background()).Err()
}

func (t *RedisFileTracker) Get(key string) (*FileStatus, error) {
	v, err := t.r.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var fileProgress FileStatus
	err = json.Unmarshal([]byte(v), &fileProgress)
	if err != nil {
		return nil, err
	}

	return &fileProgress, nil
}

func (t *RedisFileTracker) List(pat string) ([]*FileStatus, error) {
	var cursor uint64
	var n int
	var result []*FileStatus

	ctx := context.Background()

	for {
		keys, cursor, err := t.r.Scan(ctx, cursor, pat, 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			p, err := t.Get(key)
			if err != nil {
				continue
			}
			result = append(result, p)
		}

		if cursor == 0 {
			break
		}

		n += len(keys)
	}

	return result, nil
}

// Ensure implementation
var _ FileTracker = (*RedisFileTracker)(nil)
