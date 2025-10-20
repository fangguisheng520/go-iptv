package dto

type AdminClientDto struct {
	LoginUser   string   `json:"loginuser"`
	Title       string   `json:"title"`
	ServerUrl   string   `json:"serverurl"`
	Build       Build    `json:"build"`
	App         App      `json:"app"`
	Tips        Tips     `json:"tips"`
	IconUrl     string   `json:"iconurl"`
	BjUrl       []string `json:"bjurl"`
	UpSize      string   `json:"upsize"`
	ApkUrl      string   `json:"apkurl"`
	ApkName     string   `json:"apkname"`
	BuildStatus int64    `json:"status"` // APK编译状态
}
