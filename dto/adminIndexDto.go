package dto

type AdminIndexDto struct {
	// IndexDto is the DTO for the index page
	// It contains the title and the description of the page
	LoginUser        string        `json:"loginuser"`
	Title            string        `json:"title"`
	UserTotal        int64         `json:"usertotal"`
	UserToday        int64         `json:"usertoday"`
	ChannelTypeCount int64         `json:"channeltypecount"`
	MealsCount       int64         `json:"mealscount"`
	EpgCount         int64         `json:"epgcount"`
	ChannelCount     int64         `json:"channelcount"`
	ChannelTypeList  []ChannelType `json:"channeltypelist"`
}

type ChannelType struct {
	Num          int64  `json:"num"`
	Name         string `json:"name"`
	ChannelCount int64  `json:"channelcount"`
}
