package manage

import (
	"Graduation/global"
	requ "Graduation/model/common/request"
	"Graduation/model/manage"
	mage "Graduation/model/manage"
	"Graduation/model/manage/request"
	mallservice "Graduation/service/mall"
	"Graduation/utils"
	"Graduation/utils/enum"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type ManageGoodsInfoService struct {
}

// CreateMallGoodsInfo 创建MallGoodsInfo
func (m *ManageGoodsInfoService) CreateMallGoodsInfo(token string, req request.GoodsInfoAddParam) (err error) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&mage.MallAdminUser{}).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	var goodsCategory manage.MallGoodsCategory
	err = global.GVA_DB.Where("category_id=?  AND is_deleted=0", req.GoodsCategoryId).First(&goodsCategory).Error
	if goodsCategory.CategoryLevel != enum.LevelThree.Code() {
		return errors.New("分类数据异常")
	}
	if !errors.Is(global.GVA_DB.Where("goods_name=? AND goods_category_id=?", req.GoodsName, req.GoodsCategoryId).First(&manage.MallGoodsInfo{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同的商品信息")
	}
	originalPrice := utils.Transfer(req.OriginalPrice)
	sellingPrice := utils.Transfer(req.SellingPrice)
	stockNum := utils.Transfer(req.StockNum)
	goodsSellStatus, _ := strconv.Atoi(req.GoodsSellStatus)
	goodsInfo := manage.MallGoodsInfo{
		UUid:               uuid, // 添加商品的管理员id
		GoodsName:          req.GoodsName,
		GoodsIntro:         req.GoodsIntro,
		GoodsCategoryId:    req.GoodsCategoryId,
		GoodsCoverImg:      req.GoodsCoverImg,
		GoodsDetailContent: req.GoodsDetailContent,
		OriginalPrice:      originalPrice,
		SellingPrice:       sellingPrice,
		StockNum:           stockNum,
		Tag:                req.Tag,
		GoodsSellStatus:    goodsSellStatus,
		CommodityStocks:    stockNum, // 库存位
		CreateUser:         uuid,
	}
	if err = utils.Verify(goodsInfo, utils.GoodsAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	err = global.GVA_DB.Create(&goodsInfo).Error
	// 进行商品区块链操作
	var blockChainGroup []manage.MallBlockChain

	// 循环创建创始区块哈希
	for i := 1; i <= stockNum; i++ {
		// 创建创始区块
		tmpBlock := utils.GenerateGenesisBlock(utils.BlockUserInfo{
			Muid:    uuid,
			Buid:    uuid,
			GoodsId: goodsInfo.GoodsId,
			Count:   i,
		})
		// var blockChain manage.MallBlockChain
		blockChainGroup = append(blockChainGroup, manage.MallBlockChain{
			UUid:            uuid,
			Commodity:       goodsInfo.GoodsId, // 商品序号
			CommodityStocks: i,                 // 库存位
			IsSale:          false,
			Number:          0,
			InitBlockHash:   tmpBlock.CurrBlockHash,
			CurrBlockHash:   tmpBlock.CurrBlockHash,
		})
	}
	err = global.GVA_DB.Create(&blockChainGroup).Error
	return err
}

// 查询商品列表信息
// GetMallGoodsInfoInfoList 分页获取MallGoodsInfo记录
func (m *ManageGoodsInfoService) GetMallGoodsInfoInfoList(token string, info request.MallGoodsInfoSearch, goodsName string, goodsSellStatus string) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	var adminUser mage.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户"), list, total
	}
	// 创建db
	db := global.GVA_DB.Model(&manage.MallGoodsInfo{})
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if goodsName != "" {
		db.Where("goods_name =?", goodsName)
	}
	if goodsSellStatus != "" {
		db.Where("goods_sell_status =?", goodsSellStatus)
	}
	var mallGoodsInfos []manage.MallGoodsInfo
	// 超级管理员返回全部信息
	if adminUser.IsSuperAdmin == 1 {
		err = db.Limit(limit).Offset(offset).Order("goods_id desc").Find(&mallGoodsInfos).Error
		return err, mallGoodsInfos, total
	} else { // 管理员返回自己创建的商品信息
		err = db.Limit(limit).Offset(offset).Where("u_uid=?", uuid).Order("goods_id desc").Find(&mallGoodsInfos).Error
		return err, mallGoodsInfos, total
	}

}

// GetMallGoodsInfo 根据id和管理员信息获取MallGoodsInfo记录
func (m *ManageGoodsInfoService) GetMallGoodsInfo(token string, id int) (err error, mallGoodsInfo manage.MallGoodsInfo) {
	var adminUser mage.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户"), mallGoodsInfo
	}
	if adminUser.IsSuperAdmin == 1 {
		err = global.GVA_DB.Where("goods_id = ?", id).First(&mallGoodsInfo).Error
	} else {
		err = global.GVA_DB.Where("goods_id = ? AND u_uid = ?", id, uuid).First(&mallGoodsInfo).Error
	}
	return
}

// UpdateMallGoodsInfo 更新MallGoodsInfo记录
func (m *ManageGoodsInfoService) UpdateMallGoodsInfo(token string, req request.GoodsInfoUpdateParam) (err error) {
	goodsId := utils.Transfer(req.GoodsId)
	originalPrice := utils.Transfer(req.OriginalPrice)
	sellingPrice := utils.Transfer(req.SellingPrice)
	stockNum := utils.Transfer(req.StockNum)
	// 不能更新库存,一旦创建即确定
	var orgGoodsInfo manage.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id = ?", goodsId).First(&orgGoodsInfo).Error
	if orgGoodsInfo.StockNum != stockNum {
		return errors.New("库存无法更改,一旦创建即确定!")
	}
	var adminUser mage.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	goodsInfo := manage.MallGoodsInfo{ // uuid不更新
		GoodsId:            goodsId,
		GoodsName:          req.GoodsName,
		GoodsIntro:         req.GoodsIntro,
		GoodsCategoryId:    req.GoodsCategoryId,
		GoodsCoverImg:      req.GoodsCoverImg,
		GoodsDetailContent: req.GoodsDetailContent,
		OriginalPrice:      originalPrice,
		SellingPrice:       sellingPrice,
		StockNum:           stockNum,
		Tag:                req.Tag,
		GoodsSellStatus:    req.GoodsSellStatus,
		CommodityStocks:    stockNum, // 库存位
		UpdateUser:         uuid,
	}
	if err = utils.Verify(goodsInfo, utils.GoodsAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	if adminUser.IsSuperAdmin == 1 {
		err = global.GVA_DB.Where("goods_id = ?", goodsInfo.GoodsId).Updates(&goodsInfo).Error
	} else {
		err = global.GVA_DB.Where("goods_id = ? AND u_uid = ?", goodsInfo.GoodsId, uuid).Updates(&goodsInfo).Error
	}
	return err
}

// ChangeMallGoodsInfoByIds 修改商品上下架信息
func (m *ManageGoodsInfoService) ChangeMallGoodsInfoByIds(token string, ids requ.IdsReq, sellStatus string) (err error) {
	intSellStatus, _ := strconv.Atoi(sellStatus)
	//更新字段为0时，不能直接UpdateColumns
	var adminUser mage.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	if adminUser.IsSuperAdmin == 1 {
		err = global.GVA_DB.Model(&manage.MallGoodsInfo{}).Where("goods_id in ?", ids.Ids).Update("goods_sell_status", intSellStatus).Error
	} else {
		err = global.GVA_DB.Model(&manage.MallGoodsInfo{}).Where("u_uid = ? AND goods_id in ?", uuid, ids.Ids).Update("goods_sell_status", intSellStatus).Error
	}
	return err
}

// DeleteMallGoodsInfo 删除 MallGoodsInfo 商品信息记录
func (m *ManageGoodsInfoService) DeleteMallGoodsInfo(token string, ids requ.IdsReq) (err error, mallGoodsInfo manage.MallGoodsInfo) {
	var adminUser mage.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&adminUser).Error
	if err != nil || adminUser.IsSuperAdmin != 1 {
		return errors.New("不存在的超级管理员用户"), mallGoodsInfo
	}
	// 删除商品信息
	err = global.GVA_DB.Where("goods_id in ?", ids.Ids).Delete(&manage.MallGoodsInfo{}).Error
	return err, mallGoodsInfo
}
