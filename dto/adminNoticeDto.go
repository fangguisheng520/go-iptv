package dto

type AdminNoticeDto struct {
	LoginUser string `json:"loginuser"`
	Title     string `json:"title"`
	Ad        Ad     `json:"ad"`
}
