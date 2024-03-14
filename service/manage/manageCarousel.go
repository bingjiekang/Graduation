package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/manage"
	requ "Graduation/model/manage/request"
	"Graduation/utils"
	"errors"

	"gorm.io/gorm"
)

type ManageCarouselService struct {
}

// 创建轮播图
func (m *ManageCarouselService) CreateCarousel(token string, req requ.MallCarouselAddParam) (err error) {
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&manage.MallAdminUser{}).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	carouseRank := utils.Transfer(req.CarouselRank)
	mallCarousel := manage.MallCarousel{
		CarouselUrl:  req.CarouselUrl,
		RedirectUrl:  req.RedirectUrl,
		CarouselRank: carouseRank,
		CreateUser:   uuid,
	}
	// 这个校验理论上应该放在api层，但是因为前端的传值是string，而我们的校验规则是Int,所以只能转换格式后再校验
	if err = utils.Verify(mallCarousel, utils.CarouselAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	err = global.GVA_DB.Create(&mallCarousel).Error
	return err
}

// 查询轮播图列表
func (m *ManageCarouselService) GetCarouselInfoList(info requ.MallCarouselSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&manage.MallCarousel{})
	var mallCarousels []manage.MallCarousel
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("carousel_rank desc").Find(&mallCarousels).Error
	return err, mallCarousels, total
}

// 获取指定 id 的信息
func (m *ManageCarouselService) GetCarousel(id int) (err error, mallCarousel manage.MallCarousel) {
	err = global.GVA_DB.Where("carousel_id = ?", id).First(&mallCarousel).Error
	return
}

// 修改指定轮播图内容
func (m *ManageCarouselService) UpdateCarousel(token string, req requ.MallCarouselUpdateParam) (err error) {
	if errors.Is(global.GVA_DB.Where("carousel_id = ?", req.CarouselId).First(&manage.MallCarousel{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("未查询到记录！")
	}
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&manage.MallAdminUser{}).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	carouseRank := utils.Transfer(req.CarouselRank)
	mallCarousel := manage.MallCarousel{
		CarouselUrl:  req.CarouselUrl,
		RedirectUrl:  req.RedirectUrl,
		CarouselRank: carouseRank,
		CreateUser:   uuid,
	}
	// 这个校验理论上应该放在api层，但是因为前端的传值是string，而我们的校验规则是Int,所以只能转换格式后再校验
	if err = utils.Verify(mallCarousel, utils.CarouselAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	err = global.GVA_DB.Where("carousel_id = ?", req.CarouselId).UpdateColumns(&mallCarousel).Error
	return err
}

// 删除指定轮播图内容
func (m *ManageCarouselService) DeleteCarousel(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&manage.MallCarousel{}, "carousel_id in ?", ids.Ids).Error
	return err
}
