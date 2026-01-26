package quark

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zhilv666/linkchecker/internal/net"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"resty.dev/v3"
)

type QuarkProvider struct {
	pattern *regexp.Regexp
	client  *resty.Client
}

func New(client *resty.Client) *QuarkProvider {
	if client == nil {
		client = net.NewRestyClient()
	}
	return &QuarkProvider{
		pattern: regexp.MustCompile(`\/s\/([^/?#]+)`),
		client:  client,
	}
}

func (q *QuarkProvider) Name() string {
	return "夸克网盘"
}

func (q *QuarkProvider) Match(url string) bool {
	return strings.Contains(url, "pan.quark.cn")
}

func (q *QuarkProvider) Check(rawUrl, password string) (*netdisk.ShareInfo, error) {
	// 1. 提取 ID
	matches := q.pattern.FindStringSubmatch(rawUrl)
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法解析夸克分享ID")
	}
	shareID := matches[1]

	// client := net.NewRestyClient()
	client := q.client
	client.SetHeaders(map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0",
		"Content-Type": "application/json",
		"origin":       "https://pan.quark.cn",
		"referer":      "https://pan.quark.cn/",
	})

	// 获取 shard token
	var tokenResp quarkTokenResp
	_, err := client.
		R().
		SetBody(map[string]any{
			"pwd_id":   shareID,
			"passcode": password,
		}).
		SetResult(&tokenResp).
		SetError(&tokenResp).
		Post("https://drive-h.quark.cn/1/clouddrive/share/sharepage/token")
	if err != nil {
		return nil, err
	}

	if tokenResp.Code != 0 {
		if tokenResp.Code == 41004 {
			return &netdisk.ShareInfo{Status: netdisk.StatusDeleted, RawUrl: rawUrl}, nil
		}
		if tokenResp.Code == 41010 {
			return &netdisk.ShareInfo{Status: netdisk.StatusBanned, RawUrl: rawUrl}, nil
		}
		if tokenResp.Code == 41011 {
			return &netdisk.ShareInfo{Status: netdisk.StatusExpired, RawUrl: rawUrl}, nil
		}
		if tokenResp.Code == 41008 {
			return &netdisk.ShareInfo{Status: netdisk.StatusNeedPassword, RawUrl: rawUrl}, nil
		}
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, fmt.Errorf("获取 shard token 失败: [%d] %s", tokenResp.Code, tokenResp.Message)
	}
	sToken := tokenResp.Data.Stoken

	// 获取详情信息
	var detailResp quarkDetailResp
	_, err = client.
		R().
		SetQueryParams(map[string]string{
			"pwd_id":       shareID,
			"stoken":       sToken,
			"pdir_fid":     "0",
			"_fetch_share": "1",
		}).
		SetResult(&detailResp).
		SetError(&tokenResp).
		Get("https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail")
	if err != nil {
		return nil, err
	}
	if detailResp.Code != 0 {
		if detailResp.Code == 41004 {
			return &netdisk.ShareInfo{Status: netdisk.StatusDeleted, RawUrl: rawUrl}, nil
		}
		if detailResp.Code == 14001 {
			return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, nil
		}
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, fmt.Errorf("获取详情失败: [%d] %s", tokenResp.Code, tokenResp.Message)
	}

	shareData := detailResp.Data.Share
	t := time.UnixMilli(shareData.ExpiredAt)
	info := &netdisk.ShareInfo{
		Status:        netdisk.StatusValid,
		Provider:      q.Name(),
		Title:         shareData.Title,
		Size:          netdisk.FormatSize(shareData.Size),
		Author:        tokenResp.Data.Author.NickName,
		ExpiredAt:     &t,
		RawUrl:        rawUrl,
		Password:      password,
		NormalizedUrl: fmt.Sprintf("https://pan.quark.cn/s/%s", shareID),
	}

	return info, nil
}
