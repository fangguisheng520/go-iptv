package dto

import "go-iptv/models"

type AdminMealsDto struct {
	LoginUser  string             `json:"loginuser"`
	Title      string             `json:"title"`
	Meals      []models.IptvMeals `json:"meals"`
	ChannelNum int64              `json:"channelnum"`
}

type MealsReturnDto struct {
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}
