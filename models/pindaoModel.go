package models

type IptvPindao struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Url       string `gorm:"column:url" json:"url"`
	Status    int64  `gorm:"column:status" json:"status"`
	Sort      int64  `gorm:"column:sort" json:"sort"`
	EId       int64  `gorm:"column:e_id" json:"e_id"`
	CId       int64  `gorm:"column:c_id" json:"c_id"`
	ListId    int64  `gorm:"column:list_id" json:"list_id"`
	Rss_Rules string `gorm:"column:rss_rules" json:"rss_rules"`
	Epg_Rules string `gorm:"column:epg_rules" json:"epg_rules"`
	Rss_List  string `gorm:"column:rss_list" json:"rss_list"`
	Epg_List  string `gorm:"column:epg_list" json:"epg_list"`
}

func (IptvPindao) TableName() string {
	return "iptv_pindao"
}

type IptvPindaoShow struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Url       string `gorm:"column:url" json:"url"`
	Status    int64  `gorm:"column:status" json:"status"`
	Sort      int64  `gorm:"column:sort" json:"sort"`
	EId       int64  `gorm:"column:e_id" json:"e_id"`
	CId       int64  `gorm:"column:c_id" json:"c_id"`
	ListId    int64  `gorm:"column:list_id" json:"list_id"`
	EpgName   string `gorm:"column:epg_name" json:"epg_name"`
	Logo      string `gorm:"-" json:"logo"`
	Rss_Rules string `gorm:"column:rss_rules" json:"rss_rules"`
	Epg_Rules string `gorm:"column:epg_rules" json:"epg_rules"`
	Rss_List  string `gorm:"column:rss_list" json:"rss_list"`
	Epg_List  string `gorm:"column:epg_list" json:"epg_list"`
}

func (IptvPindaoShow) TableName() string {
	return "iptv_pindao"
}
