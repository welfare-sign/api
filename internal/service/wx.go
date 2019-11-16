package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/pkg/wsgin"
)

// GetWXConfig 获取微信配置
func (s *Service) GetWXConfig(ctx context.Context, url string) (*model.WXConfigResp, wsgin.APICode, error) {
	var c model.WXConfigResp
	appID := viper.GetString(config.KeyWxAppID)
	appSecret := viper.GetString(config.KeyWxAppSecret)
	c.Appid = appID
	c.Timestamp = time.Now().Unix()
	c.Noncestr = uuid.NewV4().String()

	// get jsapi_ticket
	ticket, err := s.dao.GetWXJSTicket()
	if err != nil {
		log.Warn(ctx, "GetWXConfig.GetWXJSTicket() error", zap.Error(err))
		return nil, apicode.ErrGetWXConfig, err
	}
	if ticket == "" {
		accessToken, err := s.dao.GetWXAccessToken()
		if err != nil {
			log.Warn(ctx, "GetWXConfig.GetWXAccessToken() error", zap.Error(err))
			return nil, apicode.ErrGetWXConfig, err
		}
		if accessToken == "" {
			// get access_token
			resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appId=%s&secret=%s", appID, appSecret))
			if err != nil {
				return nil, apicode.ErrGetWXConfig, err
			}
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, apicode.ErrGetWXConfig, err
			}
			accessToken = gjson.GetBytes(bytes, "access_token").String()
			if accessToken == "" {
				return nil, wsgin.APICodeSuccess, errors.New("获取access_token失败")
			}
			accessTokenExpire := time.Second * time.Duration(gjson.GetBytes(bytes, "expires_in").Int()-60)
			if err := s.dao.StoreWXAccessToken(accessToken, accessTokenExpire); err != nil {
				log.Warn(ctx, "GetWXConfig.StoreWXAccessToken() error", zap.Error(err))
				return nil, apicode.ErrGetWXConfig, err
			}
		}

		resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken))
		if err != nil {
			return nil, apicode.ErrGetWXConfig, err
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, apicode.ErrGetWXConfig, err
		}
		ticket = gjson.GetBytes(bytes, "ticket").String()
		if ticket == "" {
			return nil, wsgin.APICodeSuccess, errors.New("获取jsapi_ticket失败")
		}
		ticketExpire := time.Second * time.Duration(gjson.GetBytes(bytes, "expires_in").Int()-60)
		if err := s.dao.StoreWXJSTicket(ticket, ticketExpire); err != nil {
			log.Warn(ctx, "GetWXConfig.StoreWXJSTicket() error", zap.Error(err))
			return nil, apicode.ErrGetWXConfig, err
		}
	}
	sign := util.GetWXSignString(map[string]string{
		"noncestr":     c.Noncestr,
		"jsapi_ticket": ticket,
		"timestamp":    strconv.FormatInt(c.Timestamp, 10),
		"url":          url,
	})
	c.Signature = sign
	return &c, wsgin.APICodeSuccess, nil
}
