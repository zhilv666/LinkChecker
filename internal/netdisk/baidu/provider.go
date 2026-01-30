package baidu

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/pkg/request"
	"resty.dev/v3"

	"github.com/antchfx/htmlquery"
)

type BaiduProvider struct {
	patternS       *regexp.Regexp
	patternInit    *regexp.Regexp
	patternContent *regexp.Regexp
	client         *resty.Client
}

func New(client *resty.Client) *BaiduProvider {
	if client == nil {
		client = request.NewRestyClient()
	}
	return &BaiduProvider{
		patternS:       regexp.MustCompile(`baidu\.com\/s\/([a-zA-Z0-9_-]+)`),
		patternInit:    regexp.MustCompile(`surl=([a-zA-Z0-9_-]+)`),
		patternContent: regexp.MustCompile(`locals\.mset\(([\s\S]*?)\);`),
		client:         client,
	}
}

func (b *BaiduProvider) Name() string {
	return "百度网盘"
}

func (b *BaiduProvider) Match(url string) bool {
	return strings.Contains(url, "pan.baidu.com")
}

func (b *BaiduProvider) Check(rawUrl, password string) (*netdisk.ShareInfo, error) {
	u, err := url.Parse(rawUrl)
	if err == nil {
		pwd := u.Query().Get("pwd")
		if pwd != "" {
			password = pwd
		}
	}

	var shareID string
	if matches := b.patternS.FindStringSubmatch(rawUrl); len(matches) > 1 {
		shareID = matches[1]
	} else if matches := b.patternInit.FindStringSubmatch(rawUrl); len(matches) > 1 {
		shareID = "1" + matches[1]
	}

	if shareID == "" {
		return nil, fmt.Errorf("无法解析百度网盘 ID")
	}

	// client := net.NewRestyClient()
	client := b.client
	client.SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0",
	})

	// 验证密码
	var verifyResp baiduVerifyResp
	_, err = client.
		R().
		SetQueryParams(map[string]string{
			"surl": shareID[1:],
		}).
		SetFormData(map[string]string{
			"pwd":       password,
			"vcode":     "",
			"vcode_str": "",
		}).
		SetHeader("Referer", fmt.Sprintf("https://pan.baidu.com/share/init?surl=%s", shareID[1:])).
		SetResult(&verifyResp).
		Post("https://pan.baidu.com/share/verify")
	if err != nil {
		fmt.Println("err: ", err)
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, nil
	}
	if verifyResp.Errno != 0 {
		// 密码错误
		if verifyResp.Errno == -9 {
			return &netdisk.ShareInfo{Status: netdisk.StatusNeedPassword, RawUrl: rawUrl}, nil
		}
		// 链接错误
		if verifyResp.Errno == 105 {
			return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, fmt.Errorf("百度错误码: %d (%s)  -- 可能是链接错误", verifyResp.Errno, verifyResp.ErrMsg)
		}
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, fmt.Errorf("百度错误码: %d (%s)", verifyResp.Errno, verifyResp.ErrMsg)
	}

	// 获取数据
	resp, err := client.
		R().
		Get("https://pan.baidu.com/s/" + shareID)
	if err != nil {
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, nil
	}

	html := resp.String()
	matches := b.patternContent.FindStringSubmatch(string(resp.Bytes()))
	if len(matches) < 2 {
		return nil, fmt.Errorf("解析页面数据失败，未找到 locals.mset")
	}

	var expiredAt *time.Time

	if doc, err := htmlquery.Parse(strings.NewReader(html)); err == nil {
		// 安全查找节点
		if validDateNode := htmlquery.FindOne(doc, `//div[contains(@class, "share-valid-check")]`); validDateNode != nil {
			text := htmlquery.InnerText(validDateNode)
			// 安全分割字符串 (防止 panic)
			parts := strings.Split(text, "：")
			if len(parts) > 1 {
				dateStr := strings.TrimSpace(parts[1])
				if dateStr != "永久有效" {
					if t, err := time.Parse("2006-01-02 15:04", dateStr); err == nil {
						expiredAt = &t
					}
				} else {
					// 1. 获取当前时间上下文（主要是为了获取 Location，防止时区问题）
					now := time.Now()

					// 2. 构造新时间：(当前年份+100), 1月, 1日, 0时, 0分, 0秒, 0纳秒, 当前时区
					t := time.Date(now.Year()+100, 1, 1, 0, 0, 0, 0, now.Location())

					// 3. 赋值指针
					expiredAt = &t
				}
			}
		}
	}

	var dataResp baiduDataResp
	if err = json.Unmarshal([]byte(matches[1]), &dataResp); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if dataResp.Errno != 0 {
		//啊哦，你来晚了，分享的文件已经被取消了，下次要早点哟。
		if dataResp.Errno == -7 {
			return &netdisk.ShareInfo{Status: netdisk.StatusDeleted, RawUrl: rawUrl}, nil
		}
		// 此链接分享内容可能因为涉及侵权、色情、反动、低俗等信息，无法访问！
		if dataResp.Errno == 115 {
			return &netdisk.ShareInfo{Status: netdisk.StatusBanned, RawUrl: rawUrl}, nil
		}
		// 啊哦，来晚了，该分享文件已过期
		if dataResp.Errno == 117 {
			return &netdisk.ShareInfo{Status: netdisk.StatusExpired, RawUrl: rawUrl}, nil
		}
		// 啊哦，你来晚了，分享的文件已经被删除了，下次要早点哟。
		if dataResp.Errno == 145 {
			return &netdisk.ShareInfo{Status: netdisk.StatusDeleted, RawUrl: rawUrl}, nil
		}
		return &netdisk.ShareInfo{Status: netdisk.StatusUnknown, RawUrl: rawUrl}, fmt.Errorf("百度错误码: %d", dataResp.Errno)
	}

	if len(dataResp.FileList) == 0 {
		return &netdisk.ShareInfo{Status: netdisk.StatusDeleted, RawUrl: rawUrl}, nil
	}

	file := dataResp.FileList[0]

	var size int64 = 0
	for _, file := range dataResp.FileList {
		size += file.Size
	}

	info := &netdisk.ShareInfo{
		Status:        netdisk.StatusValid,
		Provider:      b.Name(),
		Title:         file.ServerFilename,
		Size:          netdisk.FormatSize(size),
		Author:        dataResp.Linkusername,
		ExpiredAt:     expiredAt,
		RawUrl:        rawUrl,
		Password:      password,
		NormalizedUrl: fmt.Sprintf("https://pan.baidu.com/s/%s", shareID),
	}

	return info, nil
}
