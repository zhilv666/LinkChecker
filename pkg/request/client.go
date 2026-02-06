package request

import (
	"crypto/tls"
	"net/http"
	"time"

	"resty.dev/v3"
)

var (
	NoRedirectClient *resty.Client
	RestyClient      *resty.Client
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0"
	DefaultTimeout   = 10 * time.Second
)

type Config struct {
	Timeout   time.Duration
	Proxy     string
	Debug     bool
	UserAgent string
	VerifySSL bool
}

func DefaultConfg() *Config {
	return &Config{
		Timeout:   DefaultTimeout,
		UserAgent: DefaultUserAgent,
		VerifySSL: true,
		Debug:     false,
	}
}

// NewClient 创建一个标准的 Resty 客户端
func NewClient(cfg *Config) *resty.Client {
	if cfg == nil {
		cfg = DefaultConfg()
	}

	client := resty.New()
	client.SetHeader("User-Agent", cfg.UserAgent)
	client.SetTimeout(cfg.Timeout)
	client.SetRetryCount(3)

	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: !cfg.VerifySSL,
	})

	if cfg.Proxy != "" {
		client.SetProxy(cfg.Proxy)
	}

	return client
}

// NewNoRedirectClient 创建一个禁止重定向的客户端
func NewNoRedirectClient(cfg *Config) *resty.Client {
	client := NewClient(cfg)

	client.SetRedirectPolicy(
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	)

	return client
}
