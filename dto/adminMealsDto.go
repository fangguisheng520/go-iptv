package dto

import "go-iptv/models"

type AdminMealsDto struct {
	LoginUser  string                 `json:"loginuser"`
	Title      string                 `json:"title"`
	Meals      []models.IptvMealsShow `json:"meals"`
	MealsName  string                 `json:"mealsmap"`
	ChannelNum int64                  `json:"channelnum"`
}

type MealsReturnDto struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}
