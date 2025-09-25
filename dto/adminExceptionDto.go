package dto

import "go-iptv/models"

type AdminExceptionDto struct {
	LoginUser     string            `json:"loginuser"`
	Title         string            `json:"title"`
	MaxSameipUser int               `json:"max_sameip_user"`
	Users         []models.IptvUser `json:"users"`
}
