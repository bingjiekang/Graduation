package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallShopCartApi struct {
}

// 获取购物车列表信息
func (m *MallShopCartApi) CartItemList(c *gin.Context) {
	token := c.GetHeader("token")
	if err, shopCartItem := mallShopCartService.GetShopCartItems(token); err != nil {
		global.GVA_LOG.Error("获取购物车失败", zap.Error(err))
		response.FailWithMessage("获取购物车失败:"+err.Error(), c)
	} else {
		response.OkWithData(shopCartItem, c)
	}
}

// 添加购物车
func (m *MallShopCartApi) AddMallShopCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.SaveCartItemParam
	_ = c.ShouldBindJSON(&req)
	if err := mallShopCartService.AddMallCartItem(token, req); err != nil {
		global.GVA_LOG.Error("添加购物车失败", zap.Error(err))
		response.FailWithMessage("添加购物车失败:"+err.Error(), c)
		// return
	} else {
		response.OkWithMessage("添加购物车成功", c)
	}

}

// 更新购物车信息
func (m *MallShopCartApi) UpdateMallShopCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.UpdateCartItemParam
	_ = c.ShouldBindJSON(&req)
	if err := mallShopCartService.UpdateMallCartItem(token, req); err != nil {
		global.GVA_LOG.Error("修改购物车失败", zap.Error(err))
		response.FailWithMessage("修改购物车失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("修改购物车成功", c)
}

// 删除商品
func (m *MallShopCartApi) DelMallShoppingCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	id, _ := strconv.Atoi(c.Param("newBeeMallShoppingCartItemId"))
	if err := mallShopCartService.DeleteMallCartItem(token, id); err != nil {
		global.GVA_LOG.Error("删除购物车商品失败", zap.Error(err))
		response.FailWithMessage("删除购物车商品失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("删除购物车商品成功", c)
	}
}

// 获取购物信息
func (m *MallShopCartApi) ShopTotal(c *gin.Context) {
	cartItemIdsStr := c.Query("cartItemIds")
	token := c.GetHeader("token")
	cartItemIds := utils.StrToList(cartItemIdsStr)
	if err, cartItemRes := mallShopCartService.GetCartItemsTotal(token, cartItemIds); err != nil {
		global.GVA_LOG.Error("获取购物明细异常：", zap.Error(err))
		response.FailWithMessage("获取购物明细异常:"+err.Error(), c)
	} else {
		response.OkWithData(cartItemRes, c)
	}

}
