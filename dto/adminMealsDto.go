package dto

import "go-iptv/models"

type AdminMealsDto struct {
	LoginUser string             `json:"loginuser"`
	Title     string             `json:"title"`
	Meals     []models.IptvMeals `json:"meals"`
}

type MealsReturnDto struct {
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}
