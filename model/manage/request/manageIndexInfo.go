package request

import (
	"Graduation/model/common/request"
	"Graduation/model/manage"
)

type MallIndexConfigSearch struct {
	manage.MallIndexConfig
	request.PageInfo
}

type MallIndexConfigAddParams struct {
	ConfigName  string      `json:"configName"`
	ConfigType  interface{} `json:"configType"`
	GoodsId     interface{} `json:"goodsId"`
	RedirectUrl string      `json:"redirectUrl"`
	ConfigRank  interface{} `json:"configRank"`
}

type MallIndexConfigUpdateParams struct {
	ConfigId    int         `json:"configId"`
	ConfigName  string      `json:"configName"`
	RedirectUrl string      `json:"redirectUrl"`
	ConfigType  interface{} `json:"configType"`
	GoodsId     interface{} `json:"goodsId"`
	ConfigRank  interface{} `json:"configRank"`
}
