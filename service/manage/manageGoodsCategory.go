package manage

import (
	"Graduation/global"
	"Graduation/model/manage"
	"Graduation/model/manage/request"
	"Graduation/utils"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type ManageGoodsCategoryService struct {
}

// AddCategory 添加商品分类
func (m *ManageGoodsCategoryService) AddCategory(req request.MallGoodsCategoryReq) (err error) {
	if !errors.Is(global.GVA_DB.Where("category_level=? AND category_name=? AND is_deleted=0", req.CategoryLevel, req.CategoryName).First(&manage.MallGoodsCategory{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同分类")
	}
	// fmt.Println(req)
	// 将对应数据进行字符串转数字
	rank, _ := strconv.Atoi(req.CategoryRank)
	level, _ := strconv.Atoi(req.CategoryLevel)
	pid, _ := strconv.Atoi(req.ParentId)
	category := manage.MallGoodsCategory{
		CategoryLevel: level,
		CategoryName:  req.CategoryName,
		CategoryRank:  rank,
		ParentId:      pid,
		IsDeleted:     0,
	}
	// 这个校验理论上应该放在api层，但是因为前端的传值是string，而我们的校验规则是Int,所以只能转换格式后再校验
	if err = utils.Verify(category, utils.GoodsCategoryVerify); err != nil {
		return errors.New(err.Error())
	}
	return global.GVA_DB.Create(&category).Error
}

// SelectCategoryPage 获取分类分页数据
func (m *ManageGoodsCategoryService) SelectCategoryPage(req request.SearchCategoryParams) (err error, list interface{}, total int64) {
	limit := req.PageSize
	if limit > 1000 {
		limit = 1000
	}
	offset := req.PageSize * (req.PageNumber - 1)
	db := global.GVA_DB.Model(&manage.MallGoodsCategory{})
	var categoryList []manage.MallGoodsCategory

	if utils.NumsInList(req.CategoryLevel, []int{1, 2, 3}) {
		db.Where("category_level=?", req.CategoryLevel)
	}
	if req.ParentId >= 0 {
		db.Where("parent_id=?", req.ParentId)
	}
	err = db.Where("is_deleted=0").Count(&total).Error

	if err != nil {
		return err, categoryList, total

	} else {
		db = db.Where("is_deleted=0").Order("category_rank desc").Limit(limit).Offset(offset)
		err = db.Find(&categoryList).Error
	}
	return err, categoryList, total
}
