package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/spf13/cobra"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/netdisk/baidu"
	"github.com/zhilv666/linkchecker/internal/netdisk/quark"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/request"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检测网盘链接状态",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")
		parallel, _ := cmd.Flags().GetInt("parallel")

		if parallel < 1 {
			parallel = 1
		}
		cpuNum := runtime.NumCPU()
		maxLimit := cpuNum * 2
		if parallel > maxLimit {
			fmt.Printf("⚠️  提示: 设定的并发数 (%d) 过高，已自动调整为系统建议上限: %d (CPU核数 * 2)\n", parallel, maxLimit)
			parallel = maxLimit
		}

		cache := cache.New(&cache.Config{})

		if input != "" {
			processFileMode(cache, input, output, parallel)
			return
		}

		client := request.NewRestyClient()
		manager := netdisk.NewManager(
			cache,
			baidu.New(client),
			quark.New(client),
		)

		if len(args) > 0 {
			for _, url := range args {
				info, err := manager.Check(url, "")
				if err != nil {
					fmt.Printf("解析出错: %v, rawUrl: %s", err, url)
				}
				fmt.Printf("%+v\n", info)
			}
			return
		}
		cmd.Help()
	},
}

func processFileMode(cache cache.Cache, input, output string, parallel int) {
	outputMap := make(map[string]bool)
	func() {
		if _, err := os.Stat(output); err == nil {
			file, err := os.Open(output)
			if err != nil {
				fmt.Printf("Read Output File %s Failure: %v", output, err)
				return
			}
			csvReader := csv.NewReader(file)
			csvReader.FieldsPerRecord = -1
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Read Output Line %s Failure: %v", output, err)
					continue
				}
				if len(record) > 6 {
					url := record[6]
					outputMap[url] = true
				}
			}
			fmt.Printf(">>> 已加载历史进度, 跳过 %d 个链接\n", len(outputMap))
		}
	}()

	var writer *csv.Writer
	var file *os.File
	var err error

	_, statErr := os.Stat(output)
	if statErr == nil {
		file, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		file, err = os.Create(output)
	}
	if err != nil {
		fmt.Printf("无法打开输出文件: %v\n", err)
		return
	}
	defer file.Close()

	writer = csv.NewWriter(file)
	if statErr != nil {
		writer.Write([]string{"平台", "标题", "分享者", "大小", "状态", "过期时间", "原始链接", "格式化链接", "密码"})
		writer.Flush()
	}

	inputFile, err := os.Open(input)
	if err != nil {
		fmt.Printf("无法打开输入文件: %v\n", inputFile)
		return
	}
	defer inputFile.Close()

	inputReader := csv.NewReader(inputFile)
	inputReader.FieldsPerRecord = -1

	var wg sync.WaitGroup
	var filelock sync.Mutex

	sem := make(chan struct{}, parallel)

	for {
		record, err := inputReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("跳过错误行: %v\n", err)
			continue
		}
		if len(record) == 0 {
			continue
		}

		url := record[0]
		pwd := ""
		if len(record) > 1 {
			pwd = record[1]
		}

		if outputMap[url] {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(targetUrl, targetPwd string) {
			defer wg.Done()
			defer func() { <-sem }()
			client := request.NewRestyClient()
			manager := netdisk.NewManager(
				cache,
				baidu.New(client),
				quark.New(client),
			)

			info, err := manager.Check(targetUrl, targetPwd)
			if err != nil {
				fmt.Printf("[Error] %s: %v\n", targetUrl, err)
				return
			}

			fmt.Printf("[已处理] %s: %s\n", targetUrl, info.Status.String())

			filelock.Lock()
			var expired string
			if info.ExpiredAt != nil {
				info.ExpiredAt.Format("2006-01-02 15-04-05")
			}
			writer.Write([]string{
				info.Provider,
				info.Title,
				info.Author,
				info.Size,
				info.Status.String(),
				expired,
				info.RawUrl,
				info.NormalizedUrl,
				info.Password,
			})
			writer.Flush()
			filelock.Unlock()
		}(url, pwd)
	}

	wg.Wait()
	fmt.Printf("检测完成，结果已保存至: %s\n", output)
}

func init() {
	checkCmd.Flags().StringP("file", "f", "", "输入的 csv 文件路径, 无默认值\ncsv 文件格式如下:\n 链接,密码")
	checkCmd.Flags().StringP("output", "o", "./result.csv", "输出的文件路径")
	checkCmd.Flags().IntP("parallel", "p", 1, "并发检测数量")
	rootCmd.AddCommand(checkCmd)
}
