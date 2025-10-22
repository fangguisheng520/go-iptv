package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/url"
	"strings"
)

func MealsChangeStatus(params url.Values) dto.ReturnJsonDto {
	mealId := params.Get("change_status")

	if mealId == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "套餐ID不能为空", Type: "danger"}
	}

	if mealId == "1000" {
		return dto.ReturnJsonDto{Code: 0, Msg: "默认套餐不能修改状态", Type: "danger"}
	}

	var meals models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ?", mealId).First(&meals).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "套餐 " + mealId + " 不存在", Type: "danger"}
	}
	if meals.Status == 1 {
		go until.CleanMealsXmlCacheOne(meals.ID)
		dao.DB.Model(&models.IptvMeals{}).Where("id = ?", meals.ID).Update("status", 0)
		return dto.ReturnJsonDto{Code: 1, Msg: "套餐 " + meals.Name + " 下线", Type: "success"}
	} else {
		go until.CleanMealsXmlCacheOne(meals.ID)
		dao.DB.Model(&models.IptvMeals{}).Where("id = ?", meals.ID).Update("status", 1)
		return dto.ReturnJsonDto{Code: 1, Msg: "套餐 " + meals.Name + " 上线", Type: "success"}
	}
}

func MealsEdit(params url.Values, editType int) dto.ReturnJsonDto {
	if editType == 1 {
		//编辑套餐
		mealId := params.Get("editmeal")
		if mealId == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "套餐id不能为空", Type: "danger"}
		}
		var meal models.IptvMeals
		if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ?", mealId).First(&meal).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "套餐 " + mealId + " 不存在", Type: "danger"}
		}
		var categoryList []models.IptvCategory
		if err := dao.DB.Model(&models.IptvCategory{}).Where("type != ?", "import").Find(&categoryList).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "没有频道分类信息，无法生成套餐", Type: "danger"}
		}

		mealList := strings.Split(meal.Content, "_")

		var dataList []dto.MealsReturnDto
		for _, v := range categoryList {
			var data dto.MealsReturnDto
			data.Name = v.Name
			data.Checked = false
			for _, v2 := range mealList {
				if v.Name == v2 {
					data.Checked = true
				}
			}
			dataList = append(dataList, data)
		}
		return dto.ReturnJsonDto{Code: 1, Data: dataList, Msg: "获取成功", Type: "success"}
	} else {
		var categoryList []models.IptvCategory
		if err := dao.DB.Model(&models.IptvCategory{}).Where("type != ?", "import").Find(&categoryList).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "没有频道分类信息，无法生成套餐", Type: "danger"}
		}
		var dataList []dto.MealsReturnDto
		for _, v := range categoryList {
			var data dto.MealsReturnDto
			data.Name = v.Name
			data.Checked = true

			dataList = append(dataList, data)
		}
		return dto.ReturnJsonDto{Code: 1, Data: dataList, Msg: "获取成功", Type: "success"}
	}
}

func MealsDel(params url.Values) dto.ReturnJsonDto {
	mealId := params.Get("delmeal")
	if mealId == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "没有获取到套餐ID", Type: "danger"}
	}
	if mealId == "1000" {
		return dto.ReturnJsonDto{Code: 0, Msg: "默认套餐无法删除", Type: "danger"}
	}
	if err := dao.DB.Where("id = ?", mealId).Delete(&models.IptvMeals{}).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "删除失败", Type: "danger"}
	}
	go until.CleanMealsXmlCacheAll()
	return dto.ReturnJsonDto{Code: 1, Msg: "删除成功", Type: "success"}
}

func MealsSubmit(params url.Values) dto.ReturnJsonDto {
	mealId := params.Get("mealId")
	mealName := params.Get("mealName")
	namesList := params["names[]"]

	if mealName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "套餐名称不能为空", Type: "danger"}
	}
	if len(namesList) == 0 {
		return dto.ReturnJsonDto{Code: 0, Msg: "请选择频道", Type: "danger"}
	}

	iptvMeals := models.IptvMeals{
		Name:    mealName,
		Content: strings.Join(namesList, "_"),
	}

	if mealId == "" {
		iptvMeals.Status = 1
		if err := dao.DB.Create(&iptvMeals).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "添加失败", Type: "danger"}
		}
		go until.CleanMealsXmlCacheOne(iptvMeals.ID)
		return dto.ReturnJsonDto{Code: 1, Msg: "添加成功", Type: "success"}
	} else {
		if err := dao.DB.Where("id = ?", mealId).Updates(&iptvMeals).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "编辑失败", Type: "danger"}
		}
		go until.CleanMealsXmlCacheOne(iptvMeals.ID)
		return dto.ReturnJsonDto{Code: 1, Msg: "编辑成功", Type: "success"}
	}

}
