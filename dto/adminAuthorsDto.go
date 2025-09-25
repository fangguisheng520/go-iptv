package dto

import "go-iptv/models"

type AdminAuthorsDto struct {
	LoginUser         string             `json:"loginuser"`
	Title             string             `json:"title"`
	UnAuthorUserTotal int64              `json:"unauthorusertotal"`
	NewUserToday      int64              `json:"newusertoday"`
	UserTodayAuthor   int64              `json:"usertodayauth"`
	Page              int64              `json:"page"`      // 当前页
	Order             string             `json:"order"`     // 排序字段
	Keywords          string             `json:"keywords"`  // 搜索关键字
	Users             []models.IptvUser  `json:"users"`     // 用户列表
	Meals             []models.IptvMeals `json:"meals"`     // 会员套餐列表
	RecCounts         int64              `json:"recCounts"` // 每页显示条数
	PageCount         int64              `json:"pageCount"` // 总页数
}
