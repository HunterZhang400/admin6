package utils

import (
	"context"
	"net"
	"net/http"
	"time"
)

// DefaultHTTPClient 默认httpclient，大多数场景使用此实例。
// 全局共享 HTTP 客户端是 Go 的最佳实践，因为 http.Client 是线程安全的，且内部的连接池可以有效复用连接。
var DefaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		// 连接池配置
		MaxIdleConns:        200,              // 总的空闲连接数
		MaxIdleConnsPerHost: 20,               // 每个 host 的空闲连接数
		MaxConnsPerHost:     200,              // 每个 host 的最大连接数
		IdleConnTimeout:     30 * time.Second, // 空闲连接超时
		// TCP 连接配置
		DialContext: createDialContext(10*time.Second, 30*time.Second),
		// TLS 配置
		TLSHandshakeTimeout: 10 * time.Second,
		// HTTP/2 配置
		ForceAttemptHTTP2: true,
		// 响应头超时
		ResponseHeaderTimeout: 10 * time.Second,
		// 期望继续超时
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 10 * time.Second, // 整体请求超时
}

// HighConcurrencyClient 高并发请求使用实例
var HighConcurrencyClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       30 * time.Second,
		DialContext:           createDialContext(10*time.Second, 30*time.Second),
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// FileTransferClient 配置了用于文件传输的HTTP客户端，设置了连接数、超时等参数以优化性能。
var FileTransferClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   2,
		MaxConnsPerHost:       5,
		IdleConnTimeout:       120 * time.Second,
		DialContext:           createDialContext(10*time.Second, 30*time.Second),
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	},
	Timeout: 60 * time.Minute, // 文件传输需要更长超时
}

var DeepseekHTTPClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   2,
		MaxConnsPerHost:       5,
		DialContext:           createDialContext(10*time.Second, 30*time.Second),
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	},
	Timeout: 3 * time.Minute,
}

// timeout: 连接超时，通常 3-30s;
// KeepAlive: 通常 15-60s
func createDialContext(timeout, keepalive time.Duration) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepalive,
		}
		return dialer.DialContext(ctx, network, addr)
	}
}
