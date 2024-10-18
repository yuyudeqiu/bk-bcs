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
	"time"

	"github.com/go-redis/redis/v8"
)

// ClusterClient Redis client for cluster mode
type ClusterClient struct {
	cli *redis.ClusterClient
}

// NewClusterClient init ClusterClient from config
func NewClusterClient(config Config) (*ClusterClient, error) {
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		Password:     config.Password,
		DialTimeout:  config.DialTimeout * time.Second,
		ReadTimeout:  config.ReadTimeout * time.Second,
		WriteTimeout: config.WriteTimeout * time.Second,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout * time.Second,
	})
	return &ClusterClient{cli: cli}, nil
}

// GetCli returns the underlying Redis universal client
func (c *ClusterClient) GetCli() redis.UniversalClient {
	return c.cli
}

// Ping sends a PING command to Redis to check connectivity
func (c *ClusterClient) Ping(ctx context.Context) (string, error) {
	return c.cli.Ping(ctx).Result()
}

// Get retrieves the value of the specified key
func (c *ClusterClient) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

// Exists checks if the specified keys exist
func (c *ClusterClient) Exists(ctx context.Context, key ...string) (int64, error) {
	return c.cli.Exists(ctx, key...).Result()
}

// Set sets the value of the specified key with an optional expiration time
func (c *ClusterClient) Set(
	ctx context.Context, key string, value interface{}, duration time.Duration) (string, error) {
	return c.cli.Set(ctx, key, value, duration).Result()
}

// SetNX sets the value of the specified key only if it does not already exist
func (c *ClusterClient) SetNX(
	ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.cli.SetNX(ctx, key, value, expiration).Result()
}

// SetEX sets the value of the specified key with an expiration time
func (c *ClusterClient) SetEX(
	ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return c.cli.SetEX(ctx, key, value, expiration).Result()
}

// Del deletes the specified key
func (c *ClusterClient) Del(ctx context.Context, key string) (int64, error) {
	return c.cli.Del(ctx, key).Result()
}

// Expire sets an expiration time on the specified key
func (c *ClusterClient) Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return c.cli.Expire(ctx, key, duration).Result()
}
