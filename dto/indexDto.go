package dto

type IndexDto struct {
	ApkTime    string `json:"apk_time"`
	ApkVersion string `json:"apk_version"`
	ApkUrl     string `json:"apk_url"`
	ApkSize    string `json:"apk_size"`
	ApkName    string `json:"apk_name"`
	Content    string `json:"content"`
}
