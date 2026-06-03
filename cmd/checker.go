package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/zhilv666/linkchecker/internal/dto"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/netdisk/baidu"
	"github.com/zhilv666/linkchecker/internal/netdisk/quark"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/request"
	"github.com/zhilv666/linkchecker/pkg/response"
	"resty.dev/v3"
)

type AgentChecker struct {
	manager   *netdisk.Manager
	apiClient *resty.Client
	parallel  int

	writer    *csv.Writer
	writeLock sync.Mutex
	errors    []string
	errLock   sync.Mutex
	doneMap   map[string]bool
}

func NewAgentChecker(proxy string, serverAddr, token string, parallel int) (*AgentChecker, error) {
	reqCfg := &request.Config{
		Timeout:   15 * time.Second,
		UserAgent: request.DefaultUserAgent,
		Proxy:     proxy,
	}
	checkClient := request.NewClient(reqCfg)

	var reportClient *resty.Client
	if serverAddr != "" {
		reportClient = request.NewClient(nil).SetBaseURL(serverAddr)
		if token != "" {
			reportClient.SetHeader("token", token)
		}
	}

	memCache := cache.New(&cache.Config{})
	manager := netdisk.NewManager(memCache,
		baidu.New(checkClient),
		quark.New(checkClient),
	)

	return &AgentChecker{
		manager:   manager,
		apiClient: reportClient,
		parallel:  parallel,
		doneMap:   make(map[string]bool),
	}, nil
}

func (a *AgentChecker) RunArgsMode(urls []string) {
	sem := make(chan struct{}, a.parallel)
	var wg sync.WaitGroup

	for _, url := range urls {
		sem <- struct{}{}
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			defer func() { <-sem }()
			fmt.Printf("url: %s\n", u)
			a.processOne(u, "")
		}(url)
	}
	wg.Wait()
	fmt.Println("✅ 所有命令行任务完成")
	a.displayError()
}

func (a *AgentChecker) RunFileMode(inputPath, outputPath string) {
	if a.apiClient != nil && !a.pingServer() {
		fmt.Println("⚠️  无法连接 Server，将跳过上报步骤，仅本地检测。")
	}

	a.doneMap = loadHistory(outputPath)
	if len(a.doneMap) > 0 {
		fmt.Printf(">>> 已加载 %d 条历史记录，将自动跳过\n", len(a.doneMap))
	}

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ 无法打开输出文件: %v\n", err)
		return
	}
	defer f.Close()
	a.writer = csv.NewWriter(f)
	// a.writerRow([]string{"平台", "标题", "大小", "状态", "原始链接", "链接", "密码"})

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("❌ 无法打开输入文件: %v\n", err)
		return
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.FieldsPerRecord = -1

	sem := make(chan struct{}, a.parallel)
	var wg sync.WaitGroup

	fmt.Printf(">>> 开始任务...\n")
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) == 0 {
			continue
		}

		rawUrl := record[0]
		pwd := ""
		if len(record) > 1 {
			pwd = record[1]
		}
		if a.doneMap[rawUrl] {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(u, p string) {
			defer wg.Done()
			defer func() { <-sem }()
			a.processOne(u, p)
		}(rawUrl, pwd)
	}

	wg.Wait()

	fmt.Println("✅ 所有任务完成")
	a.displayError()
}

func (a *AgentChecker) processOne(url, pwd string) {
	// 1. 本地检测 (消耗本地 CPU 和 代理流量)
	fmt.Printf("🔍 [检测中] %s\n", url)
	info, err := a.manager.Check(url, pwd)
	if err != nil {
		fmt.Printf("❌ [检测失败] %s: %v\n", url, err)

		a.errLock.Lock()
		a.errors = append(a.errors, fmt.Sprintf("%s -> %v", url, err))
		a.errLock.Unlock()
		return
	}

	if a.apiClient != nil && info.Status == 1 {
		fmt.Printf("🚀 [上报中] %s -> Server\n", info.Title)

		payload := dto.ReportReq{
			Provider:  info.Provider,
			Title:     info.Title,
			Size:      info.Size,
			Author:    info.Author,
			ExpiredAt: info.ExpiredAt,
			RawURL:    info.RawUrl,
			Status:    info.Status,
			URL:       info.NormalizedUrl,
			PWD:       info.Password,
		}
		var resp response.Response
		_, err := a.apiClient.R().
			SetBody(payload).
			SetResult(&resp).
			Post("/api/v1/report")

		if err != nil {
			fmt.Printf("⚠️ [上报网络错误] %s: %v\n", url, err)
		} else if resp.Code != 0 {
			fmt.Printf("⚠️ [上报被拒] %s Code: %d, Msg: %s\n", url, resp.Code, resp.Msg)
		} else {
			fmt.Printf("☁️ [已上报] %s\n", url)
		}
	}
	fmt.Printf("✅ [完成] %s | %s\n", url, info.Status.String())

	if a.writer != nil {
		a.writerRow([]string{
			info.Provider,
			info.Title,
			info.Size,
			info.Status.String(),
			url,
			info.NormalizedUrl,
			info.Password,
		})
	}
}

func (a *AgentChecker) writerRow(row []string) {
	a.writeLock.Lock()
	defer a.writeLock.Unlock()
	a.writer.Write(row)
	a.writer.Flush()
}

func (a *AgentChecker) pingServer() bool {
	if a.apiClient == nil {
		return false
	}
	_, err := a.apiClient.R().Get("/ping")
	if err != nil {
		fmt.Printf("❌ 连接服务器失败: %v\n", err)
		return false
	}
	return true
}

func (a *AgentChecker) displayError() {
	if len(a.errors) > 0 {
		fmt.Println("\n==============================")
		fmt.Printf("⚠️  共发现 %d 个检测错误:\n", len(a.errors))
		fmt.Println("==============================")

		for _, msg := range a.errors {
			// 使用 %s 打印字符串
			fmt.Printf("\t❌ %s\n", msg)
		}
		fmt.Println("")
	}
}

func loadHistory(output string) map[string]bool {
	m := make(map[string]bool)
	f, err := os.Open(output)
	if err != nil {
		return m
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) > 4 {
			url := record[4]
			m[url] = true
		}
	}
	return m
}
