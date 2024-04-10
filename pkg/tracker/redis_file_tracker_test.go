package tracker

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestRedisFileTracker_Create(t *testing.T) {
	client, mock := redismock.NewClusterMock()

	key := "new_key"
	status := &FileStatus{Status: "status"}
	val, _ := json.Marshal(status)

	mock.
		CustomMatch(func(expected, actual []interface{}) error {
			var want, got FileStatus

			if err := json.Unmarshal([]byte(expected[2].(string)), &want); err != nil {
				return err
			}

			if err := json.Unmarshal(actual[2].([]byte), &got); err != nil {
				return err
			}

			got.UpdatedAt = want.UpdatedAt
			if !reflect.DeepEqual(want, got) {
				return fmt.Errorf("expected %v, got %v", want, got)
			}

			return nil
		}).
		ExpectSet(key, string(val), time.Duration(0)).
		SetVal(string(val))

	tracker := NewRedisFileTracker(client, 0)

	err := tracker.Create(key, status)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisFileTracker_Get(t *testing.T) {
	client, mock := redismock.NewClusterMock()

	key := "new_key"
	want := &FileStatus{Status: "status"}
	val, _ := json.Marshal(want)

	mock.ExpectGet(key).SetVal(string(val))

	tracker := NewRedisFileTracker(client, 0)

	got, err := tracker.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}
