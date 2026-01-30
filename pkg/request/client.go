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

var DefaultTimeout = 10 * time.Second

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0"

func InitClient() {
	NoRedirectClient = resty.New().
		SetRedirectPolicy(
			resty.RedirectPolicyFunc(
				func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}),
		).SetTLSClientConfig(
		&tls.Config{
			InsecureSkipVerify: false,
		})
	// SetProxy("http://127.0.0.1:9000")
	NoRedirectClient.SetHeader("user-agent", UserAgent)

	RestyClient = NewRestyClient()
}

func NewRestyClient() *resty.Client {
	client := resty.New().
		SetHeader("user-agent", UserAgent).
		SetRetryCount(3).
		SetTimeout(DefaultTimeout).
		SetTLSClientConfig(
			&tls.Config{
				InsecureSkipVerify: false,
			})
		// SetProxy("http://127.0.0.1:9000")
	return client
}
