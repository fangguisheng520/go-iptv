package dto

import "go-iptv/models"

type AdminMovieDto struct {
	LoginUser string             `json:"loginuser"`
	Title     string             `json:"title"`
	Movies    []models.IptvMovie `json:"movies"`
}
