package server

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStorage interface {
	Store(key string, value string, exp time.Duration) error
	Load(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	List(prefix string) ([]string, error)
}

var tokenStorages = make(map[string](func(string) (TokenStorage, error)))

func RegisterTokenStorage(name string, f func(string) (TokenStorage, error)) {
	tokenStorages[name] = f
}

func NewTokenStorage(uri string) (TokenStorage, error) {
	if !strings.Contains(uri, ":") {
		return nil, errors.New("invalid token storage uri")
	}
	storageType := uri[:strings.Index(uri, ":")]

	if f, ok := tokenStorages[storageType]; ok {
		return f(uri)
	}

	return nil, errors.New("unknown token storage type")
}

type TokenStorageMemory struct {
	data *sync.Map
}

func NewTokenStorageMemory(string) (TokenStorage, error) {
	return &TokenStorageMemory{
		data: &sync.Map{},
	}, nil
}

func (ts *TokenStorageMemory) Store(key string, value string, exp time.Duration) error {
	ts.data.Store(key, value)
	time.AfterFunc(exp, func() {
		ts.Delete(key)
	})
	return nil
}

func (ts *TokenStorageMemory) Load(key string) (string, error) {
	v, ok := ts.data.Load(key)
	if !ok {
		return "", ErrTokenNotFound
	}
	return v.(string), nil
}

func (ts *TokenStorageMemory) Delete(key string) error {
	ts.data.Delete(key)
	return nil
}

func (ts *TokenStorageMemory) Exists(key string) (bool, error) {
	_, ok := ts.data.Load(key)
	return ok, nil
}

func (ts *TokenStorageMemory) List(prefix string) ([]string, error) {
	keys := []string{}
	ts.data.Range(func(k, v interface{}) bool {
		if !strings.HasPrefix(k.(string), prefix) {
			return true
		}
		keys = append(keys, k.(string))
		return true
	})
	return keys, nil
}

var ErrTokenNotFound = errors.New("token not found")

type TokenStorageRedis struct {
	client *redis.Client
	prefix string
}

func NewTokenStorageRedis(uri string) (TokenStorage, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	prefix := parsedURI.Query().Get("prefix")
	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	log.Printf("Connecting to redis using %#v", opts)
	client := redis.NewClient(opts)
	return &TokenStorageRedis{
		client: client,
		prefix: prefix,
	}, nil
}

func (ts *TokenStorageRedis) Store(key string, value string, exp time.Duration) error {
	return ts.client.Set(context.Background(), ts.prefix+key, value, exp).Err()
}

func (ts *TokenStorageRedis) Load(key string) (string, error) {
	v, err := ts.client.Get(context.Background(), ts.prefix+key).Result()
	if err == redis.Nil {
		return "", ErrTokenNotFound
	}
	return v, err
}

func (ts *TokenStorageRedis) Delete(key string) error {
	return ts.client.Del(context.Background(), ts.prefix+key).Err()
}

func (ts *TokenStorageRedis) Exists(key string) (bool, error) {
	_, err := ts.client.Get(context.Background(), ts.prefix+key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return true, err
}

func (ts *TokenStorageRedis) List(prefix string) ([]string, error) {
	var keys []string
	var cursor uint64 = 0
	matchPattern := ts.prefix + prefix + "*"

	for {
		scanKeys, nextCursor, err := ts.client.Scan(context.Background(), cursor, matchPattern, 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	for i, key := range keys {
		keys[i] = strings.TrimPrefix(key, ts.prefix)
	}

	return keys, nil
}

func init() {
	RegisterTokenStorage("memory", NewTokenStorageMemory)
	RegisterTokenStorage("redis", NewTokenStorageRedis)
	RegisterTokenStorage("rediss", NewTokenStorageRedis)
}
