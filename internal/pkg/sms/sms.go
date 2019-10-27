package sms

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

const (
	poolSize         = 2 // 最大并发数
	maxTaskQueueSize = 5 // 可缓存的最大请求数
)

var client *sdk.Client

func init() {
	var err error
	client, err = sdk.NewClientWithAccessKey(viper.GetString(config.KeySMSRegion), viper.GetString(config.KeySMSAK), viper.GetString(config.KeySMSAS))
	if err != nil {
		panic(errors.WithMessage(err, "sms client init error"))
	}
	client.EnableAsync(poolSize, maxTaskQueueSize)
}

// Send send sms
// templateValue 中 key 值全小写
// 模版地址：https://dysms.console.aliyun.com/dysms.htm?spm=5176.2020520001.106.d20dysms.27c24bd3dNl0g6#/domestic/text/template
func Send(mobile string, templateName string, templateValue map[string]string) error {
	param, err := json.Marshal(templateValue)
	if err != nil {
		return err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = viper.GetString(config.KeySMSDomain)
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = viper.GetString(config.KeySMSRegion)
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = viper.GetString(config.KeySMSSignName)
	request.QueryParams["TemplateCode"] = templateName
	request.QueryParams["TemplateParam"] = string(param)

	if _, err = client.ProcessCommonRequest(request); err != nil {
		return err
	}
	return nil
}
