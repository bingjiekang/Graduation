package mall

import (
	"Graduation/global"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
)

type MallCarouselService struct {
}

// GetIndexCarousels 首页返回固定数量的轮播图对象
func (m *MallCarouselService) GetIndexCarousels(num int) (err error, mallCarousels []manage.MallCarousel, list interface{}) {
	var carouselIndexs []response.MallIndexCarouselResponse
	err = global.GVA_DB.Where("is_deleted = 0").Order("carousel_rank desc").Limit(num).Find(&mallCarousels).Error
	for _, carousel := range mallCarousels {
		carouselIndex := response.MallIndexCarouselResponse{
			CarouselUrl: carousel.CarouselUrl,
			RedirectUrl: carousel.RedirectUrl,
		}
		carouselIndexs = append(carouselIndexs, carouselIndex)
	}
	return err, mallCarousels, carouselIndexs
}
