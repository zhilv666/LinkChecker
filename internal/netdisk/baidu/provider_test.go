package baidu

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"resty.dev/v3"
)

func TestBaiduProvider_Check(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.Client())
	defer httpmock.DeactivateAndReset()

	basePath := "../../../mockdata/baidu"

	tests := []struct {
		name           string // 测试用例名称
		inputURL       string // 输入链接
		inputPwd       string // 输入密码
		mockVerifyFile string // verify 接口 mock 文件名
		mockHTMLFile   string // html 页面 mock 文件名

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
			inputURL:       "https://pan.baidu.com/s/1r7d6wcu2IX_vbz2f6xAvuA?pwd=keaw",
			inputPwd:       "keaw",
			mockVerifyFile: "verify_0.json",
			mockHTMLFile:   "0.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusValid,
			wantProvider: "百度网盘",
			wantTitle:    "【4742】Mysql优化高级技巧、经典案例与专题 61课",
			wantAuthor:   "九福**02",
			wantSize:     "0 B",
			wantExpired:  "永久有效",
		},
		{
			name:           "文件已过期",
			inputURL:       "https://pan.baidu.com/s/16-ymY0JyvEh64MxOeZ74qw?pwd=8z10",
			inputPwd:       "8z10",
			mockVerifyFile: "verify_0.json",
			mockHTMLFile:   "117.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusExpired,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "分享的文件被取消",
			inputURL:       "https://pan.baidu.com/s/1ngl8pnjQHTiGVI3h5OIG4w?pwd=a01d",
			inputPwd:       "a01d",
			mockVerifyFile: "verify_0.json",
			mockHTMLFile:   "-7.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusDeleted,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "涉及侵权、色情、反动、低俗等信息",
			inputURL:       "https://pan.baidu.com/s/10KsuTVkwIzDeTAIbJL7w0w?pwd=10z5",
			inputPwd:       "10z5",
			mockVerifyFile: "verify_0.json",
			mockHTMLFile:   "115.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusBanned,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "分享的文件被删除",
			inputURL:       "https://pan.baidu.com/s/1SoAy2TUpD4bxJZjsIZt7kw?pwd=2jbr",
			inputPwd:       "2jbr",
			mockVerifyFile: "verify_0.json",
			mockHTMLFile:   "145.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusDeleted,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "密码错误",
			inputURL:       "https://pan.baidu.com/s/12e8VjUcC_gDC_dbO0nWfhw?pwd=Spri",
			inputPwd:       "Spri",
			mockVerifyFile: "verify_-9.json",
			mockHTMLFile:   "0.html",

			wantErr:      false,
			wantStatus:   netdisk.StatusNeedPassword,
			wantProvider: "",
			wantTitle:    "",
			wantAuthor:   "",
			wantSize:     "",
			wantExpired:  "",
		},
		{
			name:           "链接错误",
			inputURL:       "https://pan.baidu.com/s/-TRXBxWFenCd8wwVLkk_Q?pwd=csnv",
			inputPwd:       "csnv",
			mockVerifyFile: "verify_105.json",
			mockHTMLFile:   "0.html",

			wantErr:      true,
			wantStatus:   netdisk.StatusUnknown,
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

			mockVerify, err := httpmock.NewJsonResponder(200, httpmock.File(fmt.Sprintf("%s/%s", basePath, tt.mockVerifyFile)))
			assert.NoError(t, err)
			httpmock.RegisterResponder("POST", "https://pan.baidu.com/share/verify", mockVerify)

			mockHTML := httpmock.NewStringResponder(200, httpmock.File(fmt.Sprintf("%s/%s", basePath, tt.mockHTMLFile)).String())
			httpmock.RegisterResponder("GET", "=~https://pan.baidu.com/s/", mockHTML)

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
					now := time.Now()

					// 重新构建预期的 "100年后 1月1日"
					expectedTime := time.Date(now.Year()+100, 1, 1, 0, 0, 0, 0, now.Location())

					// 比较年份、月份、日期即可，避免微妙的时间差（虽然 Date 构造应该很稳）
					assert.Equal(t, expectedTime.Year(), info.ExpiredAt.Year())
					assert.Equal(t, expectedTime.Month(), info.ExpiredAt.Month())
					assert.Equal(t, expectedTime.Day(), info.ExpiredAt.Day())
					return
				}

				ts, err := time.Parse("2006-01-02 15:04:05", tt.wantExpired)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)

				assert.Equal(t, ts.Year(), info.ExpiredAt.Year())
				assert.Equal(t, ts.Month(), info.ExpiredAt.Month())
				assert.Equal(t, ts.Day(), info.ExpiredAt.Day())
			}
		})
	}

}
