package dto

type AdminClientDto struct {
	LoginUser   string   `json:"loginuser"`
	Title       string   `json:"title"`
	Build       Build    `json:"build"`
	App         App      `json:"app"`
	Tips        Tips     `json:"tips"`
	IconUrl     string   `json:"iconurl"`
	BjUrl       []string `json:"bjurl"`
	UpSize      string   `json:"upsize"`
	BuildStatus int64    `json:"status"` // APK编译状态
}
