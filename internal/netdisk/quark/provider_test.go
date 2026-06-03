package quark

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"resty.dev/v3"
)

var quarkTestLocation = time.FixedZone("UTC+8", 8*60*60)

func TestBaiduProvider_Check(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.Client())
	defer httpmock.DeactivateAndReset()

	basePath := "../../../mockdata/quark"

	tests := []struct {
		name           string // 测试用例名称
		inputURL       string // 输入链接
		inputPwd       string // 输入密码
		mockTokenFile  string // token 接口 mock 文件名
		mockDetailFile string // detail 接口 mock 文件名

		// 预期结果
		wantErr      bool
		wantStatus   netdisk.ShareStatus
		wantProvider string
		wantTitle    string
		wantAuthor   string
		wantSize     string
		wantExpired  string
	}{
		{
			name:           "正常数据",
			inputURL:       "https://pan.quark.cn/s/dd811422b7fb",
			inputPwd:       "",
			mockTokenFile:  "token_0.json",
			mockDetailFile: "detail_0.json",

			wantErr:      false,
			wantStatus:   netdisk.StatusValid,
			wantProvider: "夸克网盘",
			wantTitle:    "图灵爬虫14期",
			wantAuthor:   "夸父*010",
			wantSize:     "55.79 GB",
			wantExpired:  "2100-01-01 00:00:00",
		},
		{
			name:           "文件不存在",
			inputURL:       "https://pan.quark.cn/s/57c091bc5dcf",
			inputPwd:       "",
			mockTokenFile:  "token_41004.json",
			mockDetailFile: "detail_0.json",

			wantErr:      false,
			wantStatus:   netdisk.StatusDeleted,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "文件涉及违规内容",
			inputURL:       "https://pan.quark.cn/s/889f179a4678",
			inputPwd:       "",
			mockTokenFile:  "token_41010.json",
			mockDetailFile: "detail_0.json",

			wantErr:      false,
			wantStatus:   netdisk.StatusBanned,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "分享地址已失效",
			inputURL:       "https://pan.quark.cn/s/4321fd4b2044",
			inputPwd:       "",
			mockTokenFile:  "token_41011.json",
			mockDetailFile: "detail_0.json",

			wantErr:      false,
			wantStatus:   netdisk.StatusExpired,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "需要提取码",
			inputURL:       "https://pan.quark.cn/s/8f1426afc706",
			inputPwd:       "",
			mockTokenFile:  "token_41008.json",
			mockDetailFile: "detail_0.json",

			wantErr:      false,
			wantStatus:   netdisk.StatusNeedPassword,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			mockToken, err := httpmock.NewJsonResponder(200, httpmock.File(fmt.Sprintf("%s/%s", basePath, tt.mockTokenFile)))
			assert.NoError(t, err)
			httpmock.RegisterResponder("POST", "https://drive-h.quark.cn/1/clouddrive/share/sharepage/token", mockToken)

			mockDetail, err := httpmock.NewJsonResponder(200, httpmock.File(fmt.Sprintf("%s/%s", basePath, tt.mockDetailFile)))
			assert.NoError(t, err)
			httpmock.RegisterResponder("GET", "https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail", mockDetail)

			provider := New(client)
			info, err := provider.Check(tt.inputURL, tt.inputPwd)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, info)
			assert.Equal(t, tt.wantStatus, info.Status, "状态码不匹配")
			assert.Equal(t, tt.inputURL, info.RawUrl, "原始链接不匹配")

			if tt.wantStatus == netdisk.StatusValid {
				assert.Equal(t, tt.inputPwd, info.Password)
			}

			assert.Equal(t, tt.wantProvider, info.Provider, "提供商不匹配")

			if tt.wantTitle != "" {
				assert.Equal(t, tt.wantTitle, info.Title, "标题不匹配")
			}
			if tt.wantAuthor != "" {
				assert.Equal(t, tt.wantAuthor, info.Author, "作者不匹配")
			}
			if tt.wantSize != "" {
				assert.Equal(t, tt.wantSize, info.Size, "文件大小不匹配")
			}
			if tt.wantExpired != "" {
				assert.NotNil(t, info.ExpiredAt, "过期时间为空")

				if tt.wantExpired == "永久有效" {
					now := time.Now().In(quarkTestLocation)

					// 重新构建预期的 "100年后 1月1日"
					expectedTime := time.Date(now.Year()+100, 1, 1, 0, 0, 0, 0, quarkTestLocation)
					actualTime := info.ExpiredAt.In(quarkTestLocation)

					// 比较年份、月份、日期即可，避免微妙的时间差（虽然 Date 构造应该很稳）
					assert.Equal(t, expectedTime.Year(), actualTime.Year())
					assert.Equal(t, expectedTime.Month(), actualTime.Month())
					assert.Equal(t, expectedTime.Day(), actualTime.Day())
					return
				}

				ts, err := time.ParseInLocation("2006-01-02 15:04:05", tt.wantExpired, quarkTestLocation)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				actualTime := info.ExpiredAt.In(quarkTestLocation)

				assert.Equal(t, ts.Year(), actualTime.Year())
				assert.Equal(t, ts.Month(), actualTime.Month())
				assert.Equal(t, ts.Day(), actualTime.Day())
			}
		})
	}

}
