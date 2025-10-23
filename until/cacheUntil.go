package until

import (
	"go-iptv/dao"
	"go-iptv/models"
	"log"
	"strconv"
)

func CleanMealsXmlCacheAll() {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Find(&meals)
	for _, meal := range meals {
		dao.Cache.Delete("rssEpgXml_" + strconv.FormatInt(meal.ID, 10))
	}

	MakeMealsXmlCacheAll()
}

func MakeMealsXmlCacheAll() {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Where("status = 1").Find(&meals)
	log.Println("重建套餐EPG订阅缓存")
	for _, meal := range meals {
		GetEpg(meal.ID)
	}
	log.Println("重建套餐EPG订阅缓存完成")
}

func CleanMealsXmlCacheOne(id int64) {
	log.Println("删除套餐EPG订阅缓存: ", id)
	dao.Cache.Delete("rssEpgXml_" + strconv.FormatInt(id, 10))
	GetEpg(id)
}

func CleanMealsTxtCacheAll() {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Find(&meals)
	for _, meal := range meals {
		dao.Cache.Delete("rssMealTxt_" + strconv.FormatInt(meal.ID, 10))
	}

	CleanMealsXmlCacheAll()
}

func CleanMealsTxtCacheOne(id int64) {
	log.Println("删除套餐TXT订阅缓存: ", id)
	dao.Cache.Delete("rssMealTxt_" + strconv.FormatInt(id, 10))
	CleanMealsXmlCacheOne(id)
}

func CleanAutoCacheAll() {
	var ca []models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("enable = 1 and type = ?", "auto").Find(&ca)
	for _, ca := range ca {
		log.Println("删除自动聚合缓存: ", ca.Name)
		dao.Cache.Delete("autoCategory_" + strconv.FormatInt(ca.ID, 10))
	}
	CleanMealsTxtCacheAll()
}
