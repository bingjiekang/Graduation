package manage

import (
	"Graduation/global"
	requ "Graduation/model/common/request"
	"Graduation/model/manage"
	"Graduation/model/manage/request"
	"Graduation/utils"
	"errors"

	"gorm.io/gorm"
)

type ManageIndexConfigService struct {
}

// CreateMallIndexConfig 创建MallIndexConfig记录
func (m *ManageIndexConfigService) CreateMallIndexConfig(token string, req request.MallIndexConfigAddParams) (err error) {
	var goodsInfo manage.MallGoodsInfo
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&manage.MallAdminUser{}).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	if errors.Is(global.GVA_DB.Where("goods_id=?", req.GoodsId).First(&goodsInfo).Error, gorm.ErrRecordNotFound) {
		return errors.New("商品不存在")
	}
	if !errors.Is(global.GVA_DB.Where("config_type =? and goods_id=? and is_deleted=0", req.ConfigType, req.GoodsId).First(&manage.MallIndexConfig{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同的首页配置项")
	}
	goodsId := utils.Transfer(req.GoodsId)
	configRank := utils.Transfer(req.ConfigRank)
	configType := utils.Transfer(req.ConfigType)
	mallIndexConfig := manage.MallIndexConfig{
		ConfigName:  req.ConfigName,
		ConfigType:  configType,
		GoodsId:     goodsId,
		RedirectUrl: req.RedirectUrl,
		ConfigRank:  configRank,
		CreateUser:  uuid,
	}
	if err = utils.Verify(mallIndexConfig, utils.IndexConfigAddParamVerify); err != nil {
		return errors.New(err.Error())
	}

	err = global.GVA_DB.Create(&mallIndexConfig).Error
	return err
}

// GetMallIndexConfigInfoList 获取MallIndexConfig记录
func (m *ManageIndexConfigService) GetMallIndexConfigInfoList(info request.MallIndexConfigSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&manage.MallIndexConfig{})
	// todo 有没有更好的方式实现？
	if utils.NumsInList(info.ConfigType, []int{1, 2, 3, 4, 5}) {
		db.Where("config_type=?", info.ConfigType)
	}
	var mallIndexConfigs []manage.MallIndexConfig
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("config_rank desc").Find(&mallIndexConfigs).Error
	return err, mallIndexConfigs, total
}

// GetMallIndexConfig 根据id获取MallIndexConfig记录
func (m *ManageIndexConfigService) GetMallIndexConfig(id uint) (err error, mallIndexConfig manage.MallIndexConfig) {
	err = global.GVA_DB.Where("config_id = ?", id).First(&mallIndexConfig).Error
	return
}

// UpdateMallIndexConfig 更新MallIndexConfig记录
func (m *ManageIndexConfigService) UpdateMallIndexConfig(token string, req request.MallIndexConfigUpdateParams) (err error) {
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&manage.MallAdminUser{}).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	//更新indexConfig
	if errors.Is(global.GVA_DB.Where("goods_id = ?", req.GoodsId).First(&manage.MallGoodsInfo{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("商品不存在！")
	}
	if errors.Is(global.GVA_DB.Where("config_id=?", req.ConfigId).First(&manage.MallIndexConfig{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("未查询到记录！")
	}
	goodsId := utils.Transfer(req.GoodsId)
	configRank := utils.Transfer(req.ConfigRank)
	configType := utils.Transfer(req.ConfigType)
	// configRank, _ := strconv.Atoi(req.ConfigRank)
	mallIndexConfig := manage.MallIndexConfig{
		ConfigId:    req.ConfigId,
		ConfigType:  configType,
		ConfigName:  req.ConfigName,
		RedirectUrl: req.RedirectUrl,
		GoodsId:     goodsId,
		ConfigRank:  configRank,
		UpdateUser:  uuid,
	}
	if err = utils.Verify(mallIndexConfig, utils.IndexConfigUpdateParamVerify); err != nil {
		return errors.New(err.Error())
	}
	var newIndexConfig manage.MallIndexConfig
	err = global.GVA_DB.Where("config_type=? and goods_id=?", mallIndexConfig.ConfigType, mallIndexConfig.GoodsId).First(&newIndexConfig).Error
	if err != nil && newIndexConfig.ConfigId == mallIndexConfig.ConfigId {
		return errors.New("已存在相同的首页配置项")
	}
	err = global.GVA_DB.Where("config_id=?", mallIndexConfig.ConfigId).Updates(&mallIndexConfig).Error
	return err
}

// DeleteMallIndexConfig 删除MallIndexConfig记录
func (m *ManageIndexConfigService) DeleteMallIndexConfig(ids requ.IdsReq) (err error) {
	err = global.GVA_DB.Where("config_id in ?", ids.Ids).Delete(&manage.MallIndexConfig{}).Error
	return err
}
