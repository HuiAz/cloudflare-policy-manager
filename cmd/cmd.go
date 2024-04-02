package cmd

import (
	"cloudflare-policy-manager/cloudflare"
	"cloudflare-policy-manager/config"
	"context"
	"flag"

	"github.com/sirupsen/logrus"
)

func Run() {
	parseConf()
	ctx := context.Background()
	cloudCli := cloudflare.NewCloudflareCli(config.Conf.CloudflareConfig)
	d := cloudflare.Device{
		Cli: cloudCli,
	}
	d.ListDeviceProfileSettings(ctx, config.Conf.CloudflareAccountId)

	// 设置设备策略
	for _, p := range config.Conf.Devices {
		d.SetExcludePolicyDeviceProfileSettings(ctx, config.Conf.CloudflareAccountId, p, config.Conf.Mode, config.Conf.Rules)
	}

}

// parseConf 解析配置文件
func parseConf() {
	// 声明命令行参数变量
	helpFlag := flag.Bool("h", false, "显示帮助信息")
	configFlag := flag.String("c", "", "指定配置文件路径")

	// 解析命令行参数
	flag.Parse()

	// 如果使用了-h参数，则打印帮助信息并退出
	if *helpFlag {
		flag.Usage()
		return
	}

	// 输出-c参数的值
	logrus.Infof("config file path: %s", *configFlag)

	config.Parse(configFlag)
}
