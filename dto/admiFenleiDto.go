package dto

import "go-iptv/models"

type AdminFenleiDto struct {
	LoginUser      string                    `json:"loginuser"`
	Title          string                    `json:"title"`
	AutoUpdate     bool                      `json:"autoupdate"`
	UpdateInterval int64                     `json:"updateinterval"`
	CategoryList   []models.IptvCategoryList `json:"categorylist"`
	Categorys      []models.IptvCategory     `json:"categorys"`
	Fenlei         []models.IptvFenlei       `json:"fenlei"`
	Epgs           []models.IptvEpg          `json:"epgs"`
	EpgsList       []models.IptvEpgList      `json:"epgsList"`
}
