package dto

import "go-iptv/models"

type AdminChannelsDto struct {
	LoginUser      string                `json:"loginuser"`
	Title          string                `json:"title"`
	AutoUpdate     bool                  `json:"autoupdate"`
	UpdateInterval int64                 `json:"updateinterval"`
	CategoryList   []models.IptvCategory `json:"categorylist"`
}
