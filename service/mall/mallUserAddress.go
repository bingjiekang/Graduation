package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	requ "Graduation/model/mall/request"
	"Graduation/utils"
	"errors"

	"github.com/jinzhu/copier"
)

type MallUserAddressService struct {
}

// 用户地址保存
func (m *MallUserAddressService) AddUserAddress(token string, req requ.AddAddressParam) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 用户地址信息
	var userAddress mall.MallUserAddress
	err = copier.Copy(&userAddress, &req)
	if err != nil {
		return err
	}
	uuid, _, _ := utils.UndoToken(token)
	userAddress.Uuid = uuid
	// 判断是否为默认地址
	if req.DefaultFlag == 1 { // 新增默认地址
		// 查询是否已有默认地址
		if err = UpdateUserDefaultAddress(uuid); err != nil {
			return err
		}
	}
	// 创建新的地址
	if err = global.GVA_DB.Create(&userAddress).Error; err != nil {
		return
	}
	return
}

// GetUserAddressList 获取全部收货地址
func (m *MallUserAddressService) GetUserAddressList(token string) (err error, userAddress []mall.MallUserAddress) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在"), userAddress
	}
	uuid, _, _ := utils.UndoToken(token)
	// 得到用户全部收货地址信息
	global.GVA_DB.Where("u_uid=? and is_deleted=0", uuid).Find(&userAddress)
	return
}

// 查询对应标识地址的信息
func (m *MallUserAddressService) GetUserAddress(token string, id int64) (err error, userAddress mall.MallUserAddress) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在"), userAddress
	}
	uuid, _, _ := utils.UndoToken(token)
	// 得到用户对应标识id的收货地址信息
	global.GVA_DB.Where("u_uid = ? and address_id = ? and is_deleted = 0", uuid, id).Find(&userAddress)
	return
}

// 修改对应标识地址的信息
func (m *MallUserAddressService) UpdateUserAddress(token string, req requ.UpdateAddressParam) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 解析 token 并获得uuid
	uuid, _, _ := utils.UndoToken(token)
	// 修改对应标识的 地址信息
	var reqUserAddr mall.MallUserAddress
	// 获取对应uid和对应标识的地址信息
	if err = global.GVA_DB.Where("address_id = ? and u_uid = ?", req.AddressId, uuid).First(&reqUserAddr).Error; err != nil {
		// 如果查询不到 证明传入的地址信息标识错误
		return errors.New("用户地址不存在")
	}
	// 成功获取到对应标识的地址信息
	if uuid != reqUserAddr.Uuid { // 不是用户本人
		return errors.New("非用户本人 禁止该操作！")
	}
	// 将对应修改的信息保存到 准备修改的这个对应标识的结构体里
	err = copier.Copy(&reqUserAddr, &req)
	if err != nil {
		return
	}
	// 修改对应数据信息 并保存
	if req.DefaultFlag == 1 { // 如果是默认收货地址
		// 查询是否已有默认地址
		if err = UpdateUserDefaultAddress(uuid); err != nil {
			return err
		}
	}
	// 将对应的uuid赋值进去
	reqUserAddr.Uuid = uuid
	err = global.GVA_DB.Save(&reqUserAddr).Error
	return
}

// 查询是否有默认地址信息 有默认地址信息则直接修改,
func UpdateUserDefaultAddress(uuid int64) (err error) {
	// 查询是否已有默认地址
	var defaultUserAddress mall.MallUserAddress
	global.GVA_DB.Where("u_uid=? and default_flag =1 and is_deleted = 0", uuid).First(&defaultUserAddress)
	// 已有默认地址(将原来默认地址取消)
	if defaultUserAddress != (mall.MallUserAddress{}) {
		defaultUserAddress.DefaultFlag = 0 // 设为非默认
		err = global.GVA_DB.Save(&defaultUserAddress).Error
		if err != nil {
			return
		}
	}
	return nil
}

// 删除对应标识的用户地址信息
func (m *MallUserAddressService) DeleteUserAddress(token string, id int64) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 解析 token 并获得uuid
	uuid, _, _ := utils.UndoToken(token)
	// 修改对应标识的 地址信息
	var reqUserAddr mall.MallUserAddress
	// 获取对应uid和对应标识的地址信息
	if err = global.GVA_DB.Where("address_id = ? and u_uid = ?", id, uuid).First(&reqUserAddr).Error; err != nil {
		// 如果查询不到 证明传入的地址信息标识错误
		return errors.New("用户地址不存在")
	}
	// 成功获取到对应标识的地址信息
	if uuid != reqUserAddr.Uuid { // 不是用户本人
		return errors.New("非用户本人 禁止该操作！")
	}
	err = global.GVA_DB.Delete(&reqUserAddr).Error
	return
}

// 获取用户默认地址信息
func (m *MallUserAddressService) GetUserDefaultAddress(token string) (err error, userAddress mall.MallUserAddress) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在"), userAddress
	}
	// 解析 token 并获得uuid
	uuid, _, _ := utils.UndoToken(token)
	if err = global.GVA_DB.Where("u_uid =? and default_flag =1 and is_deleted = 0 ", uuid).First(&userAddress).Error; err != nil {
		return errors.New("不存在默认地址失败"), userAddress
	}
	return
}
