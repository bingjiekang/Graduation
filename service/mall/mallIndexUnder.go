package mall

import (
	"Graduation/global"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
	"Graduation/utils"
)

type MallIndexInfomationService struct {
}

// GetIndexInfomation 首页新品/热门/推荐返回相关IndexConfig
func (m *MallIndexInfomationService) GetIndexInfomation(configType int, num int) (err error, list interface{}) {
	var indexConfigs []manage.MallIndexConfig
	err = global.GVA_DB.Where("config_type = ?", configType).Where("is_deleted = 0").Order("config_rank desc").Limit(num).Find(&indexConfigs).Error
	if err != nil {
		return
	}
	// 获取商品id
	var ids []int
	for _, indexConfig := range indexConfigs {
		ids = append(ids, indexConfig.GoodsId)
	}
	// 获取商品信息
	var goodsList []manage.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id in ?", ids).Find(&goodsList).Error
	var indexGoodsList []response.MallIndexConfigGoodsResponse
	// 超出30个字符显示....
	for _, indexGoods := range goodsList {
		res := response.MallIndexConfigGoodsResponse{
			GoodsId:       indexGoods.GoodsId,
			GoodsName:     utils.ReplaceLength(indexGoods.GoodsName, 30),
			GoodsIntro:    utils.ReplaceLength(indexGoods.GoodsIntro, 30),
			GoodsCoverImg: indexGoods.GoodsCoverImg,
			SellingPrice:  indexGoods.SellingPrice,
			Tag:           indexGoods.Tag,
		}
		indexGoodsList = append(indexGoodsList, res)
	}
	return err, indexGoodsList
}
