package cloudflare

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"gopkg.in/yaml.v3"
)

type Device struct {
	Cli *CloudflareCli
}

type Tunnels struct {
	Rules []Rule `yaml:"rules"`
}
type Rule string

// ListDevices 列出设备
func (d *Device) ListDevices(ctx context.Context, identifier string) ([]cloudflare.DeviceManagedNetwork, error) {

	rc := &cloudflare.ResourceContainer{
		Level:      cloudflare.AccountRouteLevel,
		Identifier: identifier,
	}

	params := cloudflare.ListDeviceManagedNetworksParams{}
	resp, err := d.Cli.Api.ListDeviceManagedNetworks(ctx, rc, params)
	if err != nil {
		HandlerErrors(err)
		return nil, err
	}

	return resp, nil
}

// ListDeviceProfileSettings 列出设备策略
func (d *Device) ListDeviceProfileSettings(ctx context.Context, identifier string) ([]cloudflare.DeviceSettingsPolicy, error) {
	rc := &cloudflare.ResourceContainer{
		Level:      cloudflare.AccountRouteLevel,
		Identifier: identifier,
	}

	params := cloudflare.ListDeviceSettingsPoliciesParams{
		ResultInfo: cloudflare.ResultInfo{
			Page:    1,
			PerPage: 100,
		},
	}
	var policyList []cloudflare.DeviceSettingsPolicy

	for {
		policies, resp, err := d.Cli.Api.ListDeviceSettingsPolicies(ctx, rc, params)
		if err != nil {
			HandlerErrors(err)
			return nil, err
		}
		fmt.Println("policy数量", len(policies))
		for _, policy := range policies {
			fmt.Println("Policy ID: ", policy.PolicyID)
			fmt.Println("Policy Name: ", policy.Name)
			for _, exclude := range *policy.Exclude {
				fmt.Println("Host: ", exclude.Host)
				fmt.Println("Address: ", exclude.Address)
				fmt.Println("Description: ", exclude.Description)
			}
		}
		policyList = append(policyList, policies...)
		if resp.Page >= resp.TotalPages {
			break
		}
		params.Page += 1
	}
	return policyList, nil
}

// SetExcludePolicyDeviceProfileSettings 设置设备策略
func (d *Device) SetExcludePolicyDeviceProfileSettings(ctx context.Context, accountID, policyID string, mode string, files []string) error {
	tunnels := make([]cloudflare.SplitTunnel, 0)
	for _, file := range files {
		// 获取 yaml 文件中的规则
		t, err := d.GenerateSplitTunnelFromYaml(file)
		if err != nil {
			HandlerErrors(err)
			return err
		}
		tunnels = append(tunnels, t...)
	}

	// 更新设备策略
	_, err := d.Cli.Api.UpdateSplitTunnelDeviceSettingsPolicy(ctx, accountID, policyID, mode, tunnels)
	if err != nil {
		HandlerErrors(err)
		return err
	}
	return nil
}

// GenerateSplitTunnelFromYaml 从 yaml 文件中生成 SplitTunnel
func (d *Device) GenerateSplitTunnelFromYaml(file string) ([]cloudflare.SplitTunnel, error) {
	// 读取 yaml 文件
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var tunnels Tunnels
	err = yaml.Unmarshal(yamlFile, &tunnels)
	if err != nil {
		return nil, err
	}

	var splitTunnels []cloudflare.SplitTunnel
	for _, tunnel := range tunnels.Rules {
		if strings.Contains(string(tunnel), "全球直连") {
			_tmpList := strings.Split(string(tunnel), ",")
			suffix := _tmpList[0]
			host := _tmpList[1]

			splitTunnel1 := cloudflare.SplitTunnel{}

			splitTunnel2 := cloudflare.SplitTunnel{}
			switch suffix {
			case "DOMAIN-SUFFIX":
				splitTunnel1.Host = "*." + host
				splitTunnel2.Host = host
				splitTunnel2.Description = "Split Tunnel for " + host
				splitTunnels = append(splitTunnels, splitTunnel2)
			case "IP-CIDR":
				splitTunnel1.Address = host
			case "DOMAIN":
				splitTunnel1.Host = host
			case "IP-CIDR6":
				splitTunnel1.Address = host
			}
			splitTunnel1.Description = "Split Tunnel for " + host
			if splitTunnel1.Host != "" || splitTunnel1.Address != "" {
				splitTunnels = append(splitTunnels, splitTunnel1)
			}
			if splitTunnel2.Host != "" {
				splitTunnels = append(splitTunnels, splitTunnel2)
			}
		}
	}
	return splitTunnels, nil
}
