package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zhilv666/linkchecker/internal/conf"
	"github.com/zhilv666/linkchecker/pkg/log" // 引入你的 log 包
	"gopkg.in/yaml.v3"
)

// InitConfig 初始化配置
func InitConfig() *Config {
	// 1. 配置 Viper 查找路径
	// 策略：优先找当前目录，其次找 data 目录
	viper.SetConfigName(conf.ConfigName) // config
	viper.SetConfigType(conf.ConfigExt)  // yaml
	viper.AddConfigPath(".")             // 搜索路径1: ./
	viper.AddConfigPath(conf.ConfigDir)  // 搜索路径2: ./data

	var config *Config

	// 2. 尝试读取配置文件
	err := viper.ReadInConfig()

	// 3. 如果读不到（文件不存在），则创建默认文件
	if err != nil {
		// 判断是否是因为找不到文件
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 默认生成的路径： ./data/config.yaml
			// filepath.Join 在 Windows 会生成反斜杠，ToSlash 转成正斜杠好看点
			defaultConfigPath := filepath.ToSlash(filepath.Join(conf.ConfigDir, conf.ConfigFile))

			log.Warnf("配置文件未找到，正在生成默认配置: %s", defaultConfigPath)

			// A. 生成默认配置 (建议 DefaultConfig 改为不依赖 pwd，返回相对路径)
			config = DefaultConfig()

			// B. 确保目录存在
			// 创建 DB 目录 (如 data/)
			if err := os.MkdirAll(filepath.Dir(config.Database.DBFile), 0755); err != nil {
				log.Panicf("无法创建数据库目录: %v", err)
			}
			// 创建 Log 目录 (如 logs/)
			if err := os.MkdirAll(filepath.Dir(config.Log.Filepath), 0755); err != nil {
				log.Panicf("无法创建日志目录: %v", err)
			}
			// 确保配置文件的目录 (data/) 存在
			if err := os.MkdirAll(filepath.Dir(defaultConfigPath), 0755); err != nil {
				log.Panicf("无法创建配置目录: %v", err)
			}

			// C. 序列化为 YAML
			data, err := yaml.Marshal(config)
			if err != nil {
				log.Panicf("格式化默认配置失败: %v", err)
			}

			// D. 写入文件
			if err := os.WriteFile(defaultConfigPath, data, 0644); err != nil {
				log.Panicf("写入默认配置文件失败: %v", err)
			}

			log.Info("默认配置已生成")
			return config
		} else {
			// 文件存在但读取出错（可能是格式错了）
			log.Panicf("读取配置文件出错: %v", err)
		}
	}

	// 4. 文件存在，直接反序列化
	if err := viper.Unmarshal(&config); err != nil {
		log.Panicf("解析配置文件失败: %v", err)
	}

	// 打印加载的配置文件路径 (强制转为正斜杠)
	log.Infof("已加载配置: %s", filepath.ToSlash(viper.ConfigFileUsed()))
	return config
}

// GetDSN 将数据库配置组装为一个 DSN
func (d *Database) GetDSN() string {
	switch d.Type {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&%s",
			d.User,
			d.Password,
			d.Host,
			d.Port,
			d.Name,
			d.Params)
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai %s",
			d.Host,
			d.User,
			d.Password,
			d.Name,
			d.Port,
			d.Params)
	case "sqlite3":
		return d.DBFile
	default:
		return ""
	}
}
