package service

import (
	"errors"
	"fmt"
	"go-iptv/crontab"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func AdminGetChannels(params url.Values) string {

	var res string = ""

	category := params.Get("category")
	if category == "" {
		return res
	}

	var channels []models.IptvChannel

	dao.DB.Model(&models.IptvChannel{}).Where("category = ?", category).Order("id ASC").Find(&channels)

	for _, channel := range channels {
		res += fmt.Sprintf("%s,%s\n", channel.Name, channel.Url)
	}

	return res
}

func UpdateInterval(params url.Values) dto.ReturnJsonDto {
	updateinterval := params.Get("updateinterval")
	autoupdate := params.Get("autoupdate")

	if updateinterval == "" || updateinterval == "0" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入更新时间", Type: "danger"}
	}

	if !until.IsSafe(updateinterval) || !until.IsSafe(autoupdate) {
		return dto.ReturnJsonDto{Code: 0, Msg: "输入不合法", Type: "danger"}
	}

	if autoupdate == "" || autoupdate == "0" {
		autoupdate = "0"
	} else {
		autoupdate = "1"
	}

	autoInt, err := strconv.ParseInt(autoupdate, 10, 64)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入数字", Type: "danger"}
	}

	interval, err := strconv.ParseInt(updateinterval, 10, 64)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入数字", Type: "danger"}
	}

	cfg := dao.GetConfig()

	cfg.Channel.Auto = autoInt
	cfg.Channel.Interval = interval
	dao.SetConfig(cfg)

	if autoInt == 1 && interval > 0 {
		go crontab.Crontab()
	}

	return dto.ReturnJsonDto{Code: 1, Msg: "更新成功", Type: "success"}
}

func AddList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("listname")
	url := params.Get("listurl")
	autocategory := params.Get("autocategory")

	if listName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入频道列表", Type: "danger"}
	}

	if !until.IsSafe(listName) || !until.IsSafe(autocategory) {
		return dto.ReturnJsonDto{Code: 0, Msg: "输入不合法", Type: "danger"}
	}

	iptvCategory := models.IptvCategory{Name: listName, Url: url}
	if autocategory == "on" || autocategory == "1" || autocategory == "true" {
		iptvCategory.AutoCategory = 1
	}

	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Delete(&models.IptvCategory{})

	resp, err := http.Get(url)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	urlData := until.FilterEmoji(string(body))

	if !strings.Contains(urlData, "#genre#") && iptvCategory.AutoCategory == 1 {
		return dto.ReturnJsonDto{Code: 0, Msg: "列表并非DIYP格式,请勿启用DIYP格式自动分类", Type: "danger"}
	}

	var maxSort int
	dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
	iptvCategory.Sort = maxSort + 1

	if iptvCategory.AutoCategory == 1 {
		res := GenreChannels(listName, urlData)
		if res.Code != 0 {
			iptvCategory.Type = "import"
			dao.DB.Model(&models.IptvCategory{}).Create(&iptvCategory)
		}
		return res
	} else {
		repeat, err := AddChannelList(listName, urlData)
		if err == nil && repeat != -1 {
			iptvCategory.Type = "add"
			dao.DB.Model(&models.IptvCategory{}).Create(&iptvCategory)
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
		} else {
			return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
		}
	}
}

func UpdateList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("listname")
	if listName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入频道列表", Type: "danger"}
	}

	var iptvCategory models.IptvCategory
	res := dao.DB.Model(&models.IptvCategory{}).Select("url, autocategory").Where("name = ?", listName).First(&iptvCategory)

	if res.RowsAffected == 0 {
		return dto.ReturnJsonDto{Code: 0, Msg: "频道列表不存在", Type: "danger"}
	}

	resp, err := http.Get(iptvCategory.Url)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	urlData := until.FilterEmoji(string(body)) // 过滤emoji表情

	if iptvCategory.AutoCategory == 1 {
		if !strings.Contains(urlData, "#genre#") {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Update("autocategory", 0)
			repeat, err := AddChannelList(listName, urlData)
			if err == nil && repeat != -1 {
				return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
			} else {
				return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
			}
		}
		return GenreChannels(listName, urlData)
	} else {
		repeat, err := AddChannelList(listName, urlData)
		if err == nil && repeat != -1 {
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
		} else {
			return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
		}
	}
}

func DelList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("listname")
	if listName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入频道列表", Type: "danger"}
	}
	dao.DB.Where("name = ?", listName).Delete(&models.IptvCategory{})
	dao.DB.Where("name like ?", "%("+listName+")").Delete(&models.IptvCategory{})
	dao.DB.Where("category = ?", listName).Delete(&models.IptvChannel{})
	dao.DB.Where("category like ?", "%("+listName+")").Delete(&models.IptvChannel{})
	go BindChannel()
	return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("删除列表 %s 成功\n", listName), Type: "success"}
}

func ForbiddenChannels(params url.Values) dto.ReturnJsonDto {
	category := params.Get("category")
	if category == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	var channels []models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", category).Find(&channels)
	for _, channel := range channels {
		if channel.Enable == 1 {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", category).Updates(map[string]interface{}{
				"enable": 0,
			})
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("禁用频道 %s 成功\n", channel.Name), Type: "success"}
		} else {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", category).Updates(map[string]interface{}{
				"enable": 1,
			})
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("启用频道 %s 成功\n", channel.Name), Type: "success"}
		}
	}
	return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
}

func SubmitAddType(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	var category models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).Find(&category)
	if category.Name != "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "该频道已存在", Type: "danger"}
	}
	var maxSort int
	dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)

	dao.DB.Model(&models.IptvCategory{}).Create(&models.IptvCategory{Name: name, Enable: 1, Type: "add", Sort: maxSort + 1})
	go BindChannel() // 添加频道后重新绑定
	return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("添加频道 %s 成功\n", name), Type: "success"}
}

func SubmitDelType(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).Delete(&models.IptvCategory{})
	dao.DB.Model(&models.IptvChannel{}).Where("category = ?", name).Delete(&models.IptvChannel{})
	go BindChannel()
	return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("删除频道 %s 成功\n", name), Type: "success"}
}

/**
 * 提交编辑类型信息的函数
 * @param params url.Values 包含请求参数的URL值映射
 * @return dto.ReturnJsonDto 返回JSON格式的数据传输对象
 */
func SubmitModifyType(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	oldName := params.Get("old_name")
	if name == "" || oldName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", oldName).Update("name", name)
	dao.DB.Model(&models.IptvChannel{}).Where("category = ?", oldName).Update("category", name)
	return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("修改频道 %s 成功\n", name), Type: "success"}
}

func SubmitMoveUp(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	var current, prev models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type = ?", name, "add").First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}
	if err := dao.DB.Model(&models.IptvCategory{}).
		Where("sort < ? and type = ?", current.Sort, "add").
		Order("sort DESC").
		First(&prev).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到可交换的记录", Type: "danger"}
	}
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 交换 sort
		if err := tx.Model(&models.IptvCategory{}).
			Where("id = ?", current.ID).
			Update("sort", prev.Sort).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.IptvCategory{}).
			Where("id = ?", prev.ID).
			Update("sort", current.Sort).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "交换排序失败", Type: "danger"}
	} else {
		return dto.ReturnJsonDto{Code: 1, Msg: "交换排序成功", Type: "success"}
	}
}

func SubmitMoveDown(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}

	var current, next models.IptvCategory

	// 获取当前记录
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type = ?", name, "add").First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}

	// 获取下一条记录（sort 大于当前记录）
	if err := dao.DB.Model(&models.IptvCategory{}).
		Where("sort > ? and type = ?", current.Sort, "add").
		Order("sort ASC").
		First(&next).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到可交换的记录", Type: "danger"}
	}

	// 交换 sort
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.IptvCategory{}).
			Where("id = ?", current.ID).
			Update("sort", next.Sort).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.IptvCategory{}).
			Where("id = ?", next.ID).
			Update("sort", current.Sort).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "交换排序失败", Type: "danger"}
	}

	return dto.ReturnJsonDto{Code: 1, Msg: "交换排序成功", Type: "success"}
}

func SubmitMoveTop(params url.Values) dto.ReturnJsonDto {
	categoryName := params.Get("category")
	if categoryName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}

	var current models.IptvCategory
	if err := dao.DB.Where("name = ? and type = ?", categoryName, "add").First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}

	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 将所有记录的 sort 增加 1（为当前记录腾出最上位置）
		if err := tx.Model(&models.IptvCategory{}).
			Where("id != ? and type = ?", current.ID, "add").
			Update("sort", gorm.Expr("sort + 1")).Error; err != nil {
			return err
		}

		// 将当前记录的 sort 设置为 1（最上）
		if err := tx.Model(&models.IptvCategory{}).
			Where("id = ?", current.ID).
			Update("sort", 1).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "移动到最上失败", Type: "danger"}
	}

	return dto.ReturnJsonDto{Code: 1, Msg: "已移动到最上", Type: "success"}
}

func SubmitSave(params url.Values) dto.ReturnJsonDto {
	srclistStr := params.Get("srclist")
	categoryName := params.Get("categoryname")

	if categoryName == "" || srclistStr == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}

	// srcList := strings.Split(srclistStr, "\n")

	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type = ?", categoryName, "add").First(&models.IptvCategory{}).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}
	AddChannelList(categoryName, srclistStr)
	go BindChannel()
	return dto.ReturnJsonDto{Code: 1, Msg: "保存成功", Type: "success"}
}

func GenreChannels(listName, srclist string) dto.ReturnJsonDto {

	data := until.ConvertDataToMap(srclist)

	for genreName, genreList := range data {
		genreName = strings.TrimSpace(genreName)
		if genreName == "" {
			continue
		}

		categoryName := strings.ReplaceAll(fmt.Sprintf("%s(%s)", genreName, listName), " ", "")

		var count int64
		if err := dao.DB.Model(&models.IptvCategory{}).
			Where("name = ?", categoryName).
			Count(&count).Error; err != nil {
			continue
		}

		if count > 0 {
			AddChannelList(categoryName, genreList)
		} else {
			var maxSort int
			dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
			newCategory := models.IptvCategory{
				Name: categoryName,
				Sort: maxSort + 1,
				Type: "add",
			}

			if err := dao.DB.Create(&newCategory).Error; err != nil {
				return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("新增分类 %s 失败\n", categoryName), Type: "danger"}
			}

			AddChannelList(categoryName, genreList)
		}
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "更新成功", Type: "success"}
}

func AddChannelList(cname, srclist string) (int, error) {
	if cname == "" || srclist == "" {
		return 0, errors.New("参数不能为空")
	}

	// 转换为 "频道,URL" 格式
	srclist = until.ConvertListFormat(srclist)

	// 删除旧的分类数据
	err := dao.DB.Model(&models.IptvChannel{}).Where("category = ?", cname).Delete(&models.IptvChannel{}).Error
	if err != nil {
		return 0, err
	}

	// 取已有的 URL，用 map 去重
	existUrls := make(map[string]struct{})
	var iptvs []models.IptvChannel
	err = dao.DB.Model(&models.IptvChannel{}).Find(&iptvs).Error
	if err != nil {
		return 0, err
	}

	for _, iptv := range iptvs {
		if iptv.Url != "" { // 假设 struct 里字段是 Url
			existUrls[iptv.Url] = struct{}{}
		}
	}

	// 预处理正则清洗
	reSpaces := regexp.MustCompile(`\s+`)
	reGenre := regexp.MustCompile(`#genre#`)
	reVer := regexp.MustCompile(`ver\..*?\.m3u8`)
	reTme := regexp.MustCompile(`t\.me.*?\.m3u8`)
	reBbsok := regexp.MustCompile(`https(.*)www\.bbsok\.cf[^>]*`)

	lines := strings.Split(srclist, "\n")
	repetNum := 0

	for _, line := range lines {
		line = strings.ReplaceAll(line, " ,", ",")
		line = strings.ReplaceAll(line, "\r", "")
		line = reSpaces.ReplaceAllString(line, "")
		line = reGenre.ReplaceAllString(line, "")
		line = reVer.ReplaceAllString(line, "")
		line = reTme.ReplaceAllString(line, "")
		line = reBbsok.ReplaceAllString(line, "")

		if !strings.Contains(line, ",") {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		channelName := parts[0]
		source := parts[1]

		// 多个源分割 #
		srcList := strings.Split(source, "#")

		for _, src := range srcList {
			src2 := strings.NewReplacer(
				`"`, "",
				"'", "",
				"}", "",
				"{", "",
			).Replace(src)

			if src2 == "" || channelName == "" {
				continue
			}

			if _, exists := existUrls[src2]; exists {
				repetNum++
				continue
			}

			channel := models.IptvChannel{
				Name:     channelName,
				Url:      src2,
				Category: cname,
			}

			if err := dao.DB.Model(&models.IptvChannel{}).Create(&channel).Error; err != nil {
				continue
			}
			existUrls[src2] = struct{}{}
		}
	}
	go BindChannel()

	return repetNum, nil
}
