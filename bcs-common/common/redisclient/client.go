/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package redisclient xxx
package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

// RedisMode xxx
type RedisMode string

const (
	// SingleMode redis single mode
	SingleMode RedisMode = "single"
	// SentinelMode redis sentinel mode
	SentinelMode RedisMode = "sentinel"
	// ClusterMode redis cluster mode
	ClusterMode RedisMode = "cluster"
)

// Client defines the interface for Redis operations
type Client interface {
	GetCli() redis.UniversalClient
	Ping(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key ...string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) (string, error)
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Del(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, duration time.Duration) (bool, error)
}

// NewClient creates a Redis client based on the configuration for different deployment modes
func NewClient(config Config) (Client, error) {
	switch config.Mode {
	case SingleMode:
		return NewSingleClient(config)
	case SentinelMode:
		return NewSentinelClient(config)
	case ClusterMode:
		return NewClusterClient(config)
	}
	return nil, fmt.Errorf("invalid config mode: %s", config.Mode)
}

// NewTestClient creates a Redis client for unit testing
func NewTestClient() (Client, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return &SingleClient{cli: client}, nil
}
