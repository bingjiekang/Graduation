package mall

import (
	"Graduation/global"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
	"Graduation/utils/enum"

	"github.com/jinzhu/copier"
)

type MallGoodsCategoryService struct {
}

// 获取分类页 一二三级 分类信息(一级分类是左边栏,二级是内部标题分类,三级是各个品牌,三级点进去是商品)
func (m *MallGoodsCategoryService) GetGoodsCategories() (err error, MallIndexCategoryVOS []response.MallIndexCategoryVO) {

	// 获取并添加一级分类的固定数量的数据
	_, firstLevelCategories := selectByLevelAndParentIdsAndNumber([]int{0}, enum.LevelOne.Code())
	if len(firstLevelCategories) != 0 {
		// fmt.Println("一级分类", firstLevelCategories)
		var firstLevelCategoryIds []int
		for _, firstLevelCategory := range firstLevelCategories {
			firstLevelCategoryIds = append(firstLevelCategoryIds, firstLevelCategory.CategoryId)
		}
		// 获取并添加二级分类的数据
		_, secondLevelCategories := selectByLevelAndParentIdsAndNumber(firstLevelCategoryIds, enum.LevelTwo.Code())
		var secondLevelCategoryVOS []response.SecondLevelCategoryVO
		if len(secondLevelCategories) != 0 {
			// fmt.Println("二级分类", secondLevelCategories)
			var secondLevelCategoryIds []int
			for _, secondLevelCategory := range secondLevelCategories {
				secondLevelCategoryIds = append(secondLevelCategoryIds, secondLevelCategory.CategoryId)
			}
			// 获取并添加三级分类的数据
			_, thirdLevelCategories := selectByLevelAndParentIdsAndNumber(secondLevelCategoryIds, enum.LevelThree.Code())

			thirdLevelCategoryMap := make(map[int][]manage.MallGoodsCategory)
			if len(thirdLevelCategories) != 0 {
				// fmt.Println("三级分类", thirdLevelCategories)
				for _, thirdLevelCategory := range thirdLevelCategories {
					thirdLevelCategoryMap[thirdLevelCategory.ParentId] = []manage.MallGoodsCategory{}
				}
				// 将三级数据都放入到 thirdLevelCategoryMap
				for k, v := range thirdLevelCategoryMap {
					for _, third := range thirdLevelCategories {
						if k == third.ParentId {
							v = append(v, third)
						}
						thirdLevelCategoryMap[k] = v
					}
				}
			}
			// 处理二级分类(将三级分类中数据放入二级分类中)
			for _, secondLevelCategory := range secondLevelCategories {
				var secondLevelCategoryVO response.SecondLevelCategoryVO
				err = copier.Copy(&secondLevelCategoryVO, &secondLevelCategory)
				// 如果该二级分类下有数据则放入 secondLevelCategoryVOS 对象中
				if _, ok := thirdLevelCategoryMap[secondLevelCategory.CategoryId]; ok {
					// 根据二级分类的 id 取出 thirdLevelCategoryMap 分组中的三级分类 list
					tempGoodsCategories := thirdLevelCategoryMap[secondLevelCategory.CategoryId]
					var thirdLevelCategoryRes []response.ThirdLevelCategoryVO
					err = copier.Copy(&thirdLevelCategoryRes, &tempGoodsCategories)
					secondLevelCategoryVO.ThirdLevelCategoryVOS = thirdLevelCategoryRes
				}
				secondLevelCategoryVOS = append(secondLevelCategoryVOS, secondLevelCategoryVO)
			}
		}
		//处理二级分类
		secondLevelCategoryVOMap := make(map[int][]response.SecondLevelCategoryVO)
		if secondLevelCategoryVOS != nil {
			//根据 parentId 将 thirdLevelCategories 分组

			for _, secondLevelCategory := range secondLevelCategoryVOS {
				secondLevelCategoryVOMap[secondLevelCategory.ParentId] = []response.SecondLevelCategoryVO{}
			}
			for k, v := range secondLevelCategoryVOMap {
				for _, second := range secondLevelCategoryVOS {
					if k == second.ParentId {
						var secondLevelCategory response.SecondLevelCategoryVO
						copier.Copy(&secondLevelCategory, &second)
						v = append(v, secondLevelCategory)
					}
					secondLevelCategoryVOMap[k] = v
				}
			}
		}
		// 将二级分类数据放入一级中
		for _, firstCategory := range firstLevelCategories {
			var mallIndexCategoryVO response.MallIndexCategoryVO
			err = copier.Copy(&mallIndexCategoryVO, &firstCategory)
			// 如果该一级分类下有数据则放入 MallIndexCategoryVOS 对象中
			if _, ok := secondLevelCategoryVOMap[firstCategory.CategoryId]; ok {
				// 根据一级分类的id取出secondLevelCategoryVOMap分组中的二级级分类list
				tempGoodsCategories := secondLevelCategoryVOMap[firstCategory.CategoryId]
				mallIndexCategoryVO.SecondLevelCategoryVOS = tempGoodsCategories
			}
			MallIndexCategoryVOS = append(MallIndexCategoryVOS, mallIndexCategoryVO)
		}
	}
	// fmt.Print("结果", MallIndexCategoryVOS)
	return
}

// 获取分类数据(这里把limit去掉,全部获得)
func selectByLevelAndParentIdsAndNumber(ids []int, level int) (err error, categories []manage.MallGoodsCategory) {
	// 获取对应分类数据
	// err = global.GVA_DB.Where("parent_id in ? and category_level =? and is_deleted = 0", ids, level).Order("category_rank desc").Limit(limit).Find(&categories).Error
	err = global.GVA_DB.Where("parent_id in ? and category_level =? and is_deleted = 0", ids, level).Order("category_rank asc").Find(&categories).Error
	// fmt.Println(categories, len(categories))
	return
}
