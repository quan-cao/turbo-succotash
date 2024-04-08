package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFileTracker struct {
	r *redis.ClusterClient
}

func NewRedisFileTracker(r *redis.ClusterClient) *RedisFileTracker {
	return &RedisFileTracker{r}
}

func (t *RedisFileTracker) Add(isid string, filename string, sourceLang string, targetLang string, status string) error {
	key := fmt.Sprintf("%s_%s", isid, filename)

	duration := time.Duration(0)

	invalidationTimeStr := os.Getenv("PROGRESS_INVALIDATION_TIME")
	invalidationTime, err := strconv.Atoi(invalidationTimeStr)
	if err != nil {
		return err
	}

	if strings.HasPrefix(status, "fail:") {
		duration = time.Duration(invalidationTime) * time.Second
	}

	now := time.Now()
	progress := &FileProgress{
		Status:     &status,
		UpdatedAt:  &now,
		SourceLang: &sourceLang,
		TargetLang: &targetLang,
	}

	progressJson, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	err = t.r.Set(context.Background(), key, string(progressJson), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisFileTracker) Clear(isid string, filename string) error {
	key := fmt.Sprintf("%s_%s", isid, filename)

	err := t.r.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Ensure implementation
var _ FileTracker = (*RedisFileTracker)(nil)
