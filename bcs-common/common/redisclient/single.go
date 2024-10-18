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

package redisclient

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// SingleClient Redis client for single mode
type SingleClient struct {
	cli *redis.Client
}

// NewSingleClient init SingleClient from config
func NewSingleClient(config Config) (*SingleClient, error) {
	if config.Mode != SingleMode {
		return nil, errors.New("redis mode not supported")
	}
	if len(config.Addrs) == 0 {
		return nil, errors.New("address is empty")
	}
	cli := redis.NewClient(&redis.Options{
		Addr:         config.Addrs[0],
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	return &SingleClient{cli: cli}, nil
}

// NewSingleClientFromDSN init SingleClient by dsn
func NewSingleClientFromDSN(dsn string) (*SingleClient, error) {
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(options)
	return &SingleClient{cli: cli}, nil
}

// GetCli returns the underlying Redis universal client
func (c *SingleClient) GetCli() redis.UniversalClient {
	return c.cli
}

// Ping sends a PING command to Redis to check connectivity
func (c *SingleClient) Ping(ctx context.Context) (string, error) {
	return c.cli.Ping(ctx).Result()
}

// Get retrieves the value of the specified key
func (c *SingleClient) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

// Set sets the value of the specified key with an optional expiration time
func (c *SingleClient) Set(ctx context.Context, key string, value interface{}, duration time.Duration) (string, error) {
	return c.cli.Set(ctx, key, value, duration).Result()
}

// SetNX sets the value of the specified key only if it does not already exist
func (c *SingleClient) SetNX(
	ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.cli.SetNX(ctx, key, value, expiration).Result()
}

// SetEX sets the value of the specified key with an expiration time
func (c *SingleClient) SetEX(
	ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return c.cli.SetEX(ctx, key, value, expiration).Result()
}

// Exists sets the value of the specified key with an expiration time
func (c *SingleClient) Exists(ctx context.Context, key ...string) (int64, error) {
	return c.cli.Exists(ctx, key...).Result()
}

// Del deletes the specified key
func (c *SingleClient) Del(ctx context.Context, key string) (int64, error) {
	return c.cli.Del(ctx, key).Result()
}

// Expire sets an expiration time on the specified key
func (c *SingleClient) Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return c.cli.Expire(ctx, key, duration).Result()
}
