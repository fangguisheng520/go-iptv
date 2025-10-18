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
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminGetChannels(params url.Values) string {

	var res string = ""

	category := params.Get("category")
	if category == "" {
		return res
	}

	var categoryDb models.IptvCategory

	if err := dao.DB.Where("name = ?", category).First(&categoryDb).Error; err != nil {
		return res
	}

	if categoryDb.Sort == -2 {
		return until.GetCCTVChannelList(false)
	}
	if categoryDb.Sort == -1 {
		return until.GetProvinceChannelList(false)
	}

	var channels []models.IptvChannel

	dao.DB.Model(&models.IptvChannel{}).Where("category = ?", category).Order("sort ASC").Find(&channels)

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
		crontab.StopChan = make(chan struct{})
		go crontab.Crontab()
	}
	if autoInt == 0 {
		close(crontab.StopChan)
		crontab.CrontabStatus = false
	}

	return dto.ReturnJsonDto{Code: 1, Msg: "更新成功", Type: "success"}
}

func AddList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("listname")
	url := params.Get("listurl")
	ua := params.Get("listua")
	cId := params.Get("cId")
	autocategory := params.Get("autocategory")
	repeat := params.Get("repeat")

	if listName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入频道列表", Type: "danger"}
	}

	if !until.IsSafe(listName) || !until.IsSafe(autocategory) || !until.IsSafe(cId) {
		return dto.ReturnJsonDto{Code: 0, Msg: "输入不合法", Type: "danger"}
	}

	iptvCategory := models.IptvCategory{Name: listName, Url: url, UA: ua}

	if cId == "" {
		var category models.IptvCategory
		dao.DB.Model(&models.IptvCategory{}).Where("name = ? or url = ?", listName, url).Find(&category)
		if category.Name != "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "该列表名称或url已存在", Type: "danger"}
		}
	} else {
		id, err := strconv.ParseInt(cId, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "请输入数字", Type: "danger"}
		}
		iptvCategory.ID = id
	}

	if autocategory == "on" || autocategory == "1" || autocategory == "true" {
		iptvCategory.AutoCategory = 1
	}

	var doRepeat bool = false
	if repeat == "on" || repeat == "1" || repeat == "true" {
		iptvCategory.Repeat = 1
		doRepeat = true
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-创建请求错误:" + err.Error(), Type: "danger"}
	}

	// 添加自定义 User-Agent
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-无法访问url:" + err.Error(), Type: "danger"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-状态码:" + strconv.Itoa(resp.StatusCode), Type: "danger"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	urlData := until.FilterEmoji(string(body))

	if until.IsM3UContent(urlData) {
		urlData = until.M3UToGenreTXT(urlData)
	}

	if !strings.Contains(urlData, "#genre#") && iptvCategory.AutoCategory == 1 {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到分组, 无法使用自动分组", Type: "danger"}
	}

	if iptvCategory.ID != 0 {
		var cOld models.IptvCategory
		if err := dao.DB.Model(&models.IptvCategory{}).Where("id = ?", iptvCategory.ID).First(&cOld).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "该频道列表不存在", Type: "danger"}
		}
		dao.DB.Model(&models.IptvChannel{}).Where("category = ? or category like ?", cOld.Name, "%("+cOld.Name+")").Delete(&models.IptvChannel{})
		dao.DB.Model(&models.IptvCategory{}).Where("name like ?", "%("+cOld.Name+")").Delete(&models.IptvCategory{})
	}

	// dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Delete(&models.IptvCategory{})

	var maxSort int64
	dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
	iptvCategory.Sort = maxSort + 1

	if iptvCategory.AutoCategory == 1 {
		res := GenreChannels(listName, urlData, doRepeat)
		if res.Code != 0 {
			iptvCategory.Type = "import"
			iptvCategory.LatestTime = time.Now().Format("2006-01-02 15:04:05")
			if iptvCategory.ID != 0 {
				iptvCategory.Enable = 1
				dao.DB.Model(&models.IptvCategory{}).Where("id = ?", iptvCategory.ID).Save(&iptvCategory)
			} else {
				dao.DB.Model(&models.IptvCategory{}).Create(&iptvCategory)
			}
		}
		return res
	} else {
		dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Updates(map[string]interface{}{
			"latesttime": time.Now().Format("2006-01-02 15:04:05"),
		})
		repeat, err := AddChannelList(listName, urlData, doRepeat)
		if err == nil && repeat != -1 {
			iptvCategory.Type = "add"
			iptvCategory.LatestTime = time.Now().Format("2006-01-02 15:04:05")
			if iptvCategory.ID != 0 {
				iptvCategory.Enable = 1
				dao.DB.Model(&models.IptvCategory{}).Where("id = ?", iptvCategory.ID).Save(&iptvCategory)
			} else {
				dao.DB.Model(&models.IptvCategory{}).Create(&iptvCategory)
			}

			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
		} else {
			return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
		}
	}
}

func UpdateList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("updatelist")
	if listName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入频道列表", Type: "danger"}
	}

	crontab.UpdateStatus = true
	defer func() { crontab.UpdateStatus = false }()

	var iptvCategory models.IptvCategory
	res := dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).First(&iptvCategory)

	if res.RowsAffected == 0 {
		return dto.ReturnJsonDto{Code: 0, Msg: "频道列表不存在", Type: "danger"}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", iptvCategory.Url, nil)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-创建请求错误:" + err.Error(), Type: "danger"}
	}

	// 添加自定义 User-Agent
	req.Header.Set("User-Agent", iptvCategory.UA)

	resp, err := client.Do(req)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-无法访问url:" + err.Error(), Type: "danger"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败-状态码:" + strconv.Itoa(resp.StatusCode), Type: "danger"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取频道列表失败", Type: "danger"}
	}

	urlData := until.FilterEmoji(string(body)) // 过滤emoji表情

	if until.IsM3UContent(urlData) {
		urlData = until.M3UToGenreTXT(urlData)
	}

	var doRepeat = false
	log.Println("Repeat:", iptvCategory.Repeat)
	if iptvCategory.Repeat == 1 {
		doRepeat = true
	}

	if iptvCategory.AutoCategory == 1 {
		if !strings.Contains(urlData, "#genre#") {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Updates(map[string]interface{}{
				"latesttime":   time.Now().Format("2006-01-02 15:04:05"),
				"autocategory": 0,
			})
			repeat, err := AddChannelList(listName, urlData, doRepeat)
			if err == nil && repeat != -1 {
				dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Updates(map[string]interface{}{
					"latesttime": time.Now().Format("2006-01-02 15:04:05"),
				})
				return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
			} else {
				return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
			}
		}
		return GenreChannels(listName, urlData, doRepeat)
	} else {
		dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Updates(map[string]interface{}{
			"latesttime": time.Now().Format("2006-01-02 15:04:05"),
		})
		repeat, err := AddChannelList(listName, urlData, doRepeat)
		if err == nil && repeat != -1 {
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
		} else {
			return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
		}
	}
}

func DelList(params url.Values) dto.ReturnJsonDto {
	listName := params.Get("dellist")
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
	go until.CleanMealsXmlCacheAll() // 清除缓存
	go until.ClearAutoChannelCache() // 清除缓存
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
	var maxSort int64
	dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)

	dao.DB.Model(&models.IptvCategory{}).Create(&models.IptvCategory{Name: name, Enable: 1, Type: "user", Sort: maxSort + 1})
	// go BindChannel() // 添加频道后重新绑定
	return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("添加频道 %s 成功\n", name), Type: "success"}
}

func SubmitDelType(params url.Values) dto.ReturnJsonDto {
	name := params.Get("category")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}

	var category models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).First(&category).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "该频道不存在", Type: "danger"}
	}
	if category.ID == 1 || category.ID == 2 {
		return dto.ReturnJsonDto{Code: 0, Msg: "该频道不能删除", Type: "danger"}
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
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type not in ?", name, []string{"import", "auto"}).First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}
	if err := dao.DB.Model(&models.IptvCategory{}).
		Where("sort < ? and type not in ?", current.Sort, []string{"import", "auto"}).
		Order("sort DESC").
		First(&prev).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到可交换的记录", Type: "danger"}
	}

	if prev.Sort < 0 {
		return dto.ReturnJsonDto{Code: 0, Msg: "已在自定义分类最上", Type: "danger"}
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
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type not in ?", name, []string{"import", "auto"}).First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}

	// 获取下一条记录（sort 大于当前记录）
	if err := dao.DB.Model(&models.IptvCategory{}).
		Where("sort > ? and type not in ?", current.Sort, []string{"import", "auto"}).
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
	if err := dao.DB.Where("name = ? and type not in ?", categoryName, []string{"import", "auto"}).First(&current).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}

	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 将所有记录的 sort 增加 1（为当前记录腾出最上位置）
		if err := tx.Model(&models.IptvCategory{}).
			Where("id != ? and type not in ?", current.ID, []string{"import", "auto"}).
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

	if categoryName == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}

	// srcList := strings.Split(srclistStr, "\n")

	var category models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ? and type not in (?)", categoryName, []string{"import"}).First(&category).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到当前记录", Type: "danger"}
	}

	if category.Sort < 0 {
		return dto.ReturnJsonDto{Code: 0, Msg: "默认分类不允许修改", Type: "danger"}
	}

	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", categoryName).Updates(map[string]interface{}{
		"latesttime": time.Now().Format("2006-01-02 15:04:05"),
		"type":       "user",
	})
	AddChannelList(categoryName, srclistStr, false)
	// go BindChannel()
	return dto.ReturnJsonDto{Code: 1, Msg: "保存成功", Type: "success"}
}

func GenreChannels(listName, srclist string, doRepeat bool) dto.ReturnJsonDto {

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
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", listName).Updates(map[string]interface{}{
				"latesttime": time.Now().Format("2006-01-02 15:04:05"),
			})
			AddChannelList(categoryName, genreList, doRepeat)
		} else {
			var maxSort int64
			dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
			newCategory := models.IptvCategory{
				LatestTime: time.Now().Format("2006-01-02 15:04:05"),
				Name:       categoryName,
				Sort:       maxSort + 1,
				Type:       "add",
			}

			if err := dao.DB.Create(&newCategory).Error; err != nil {
				return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("新增分类 %s 失败\n", categoryName), Type: "danger"}
			}

			AddChannelList(categoryName, genreList, doRepeat)
		}
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "更新成功", Type: "success"}
}

func AddChannelList(cname, srclist string, doRepeat bool) (int, error) {
	if cname == "" {
		return 0, errors.New("参数不能为空")
	}
	if srclist == "" {
		// 如果 srclist 为空，删除当前分类下所有数据
		if err := dao.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Delete(&models.IptvChannel{}, "category = ?", cname).Error
		}); err != nil {
			return 0, err
		}
		go BindChannel()
		// go until.UpdateChannelsId()
		return 0, nil
	}

	// 转换为 "频道,URL" 格式
	srclist = until.ConvertListFormat(srclist)

	// 获取 cname 分类下已有的频道
	var oldChannels []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Where("category = ?", cname).Find(&oldChannels).Error; err != nil {
		return 0, err
	}

	// 当前分类已有 URL -> channelName（大小写敏感）
	existMap := make(map[string]string)
	for _, ch := range oldChannels {
		if ch.Url != "" && ch.Name != "" {
			existMap[ch.Url] = ch.Name
		}
	}

	var handChannels []models.IptvChannel
	existHandMap := make(map[string]string)
	if doRepeat {
		dao.DB.Table(models.IptvChannel{}.TableName()+" AS c").
			Select("c.name, c.url").
			Joins("LEFT JOIN "+models.IptvCategory{}.TableName()+" AS cat ON c.category = cat.name and cat.enable = 1").
			Where("cat.type = ?", "user").
			Scan(&handChannels)

		for _, ch := range handChannels {
			if ch.Url != "" && ch.Name != "" {
				existHandMap[ch.Url] = ch.Name
			}
		}
	}

	// 正则清洗
	reSpaces := regexp.MustCompile(`\s+`)
	reGenre := regexp.MustCompile(`#genre#`)
	reVer := regexp.MustCompile(`ver\..*?\.m3u8`)
	reTme := regexp.MustCompile(`t\.me.*?\.m3u8`)
	reBbsok := regexp.MustCompile(`https(.*)www\.bbsok\.cf[^>]*`)

	lines := strings.Split(srclist, "\n")
	newChannels := make([]models.IptvChannel, 0)
	srclistUrls := make(map[string]struct{})
	repetNum := 0
	delIDs := make([]int64, 0)
	var sortIndex int64 = 1
	// +++ 新增：原始有效频道计数器 +++
	var rawCount int64 = 0

	// 先处理循环，准备新增和标记要删除的旧数据
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

		srcList := strings.Split(source, "#")
		for _, src := range srcList {
			src2 := strings.Trim(strings.NewReplacer(`"`, "", "'", "", "}", "", "{", "").Replace(src), " \r\n\t")
			if src2 == "" || channelName == "" {
				continue
			}
			// +++ 新增：每处理一个有效频道就计数 +++
			rawCount++

			srclistUrls[src2] = struct{}{}

			if doRepeat {
				if _, exists := existHandMap[src2]; exists {
					for _, ch := range oldChannels {
						if ch.Url == src2 {
							delIDs = append(delIDs, ch.ID)
						}
					}
					repetNum++
					continue
				}
			}

			if oldName, exists := existMap[src2]; exists {
				if oldName != channelName {
					// URL 相同但 channelName 不同 → 删除旧数据
					for _, ch := range oldChannels {
						if ch.Url == src2 {
							delIDs = append(delIDs, ch.ID)
						}
					}
				} else {
					// URL + channelName 相同 → 检查顺序
					for _, ch := range oldChannels {
						if ch.Url == src2 && ch.Name == channelName && ch.Sort != sortIndex {
							ch.Sort = sortIndex
							if err := dao.DB.Model(&models.IptvChannel{}).Where("id = ?", ch.ID).Update("sort", sortIndex).Error; err != nil {
								log.Println("更新顺序失败:", err)
							}
							break
						}
					}
					sortIndex++
					continue
				}
			}

			// 新增数据
			newChannels = append(newChannels, models.IptvChannel{
				Name:     channelName,
				Url:      src2,
				Category: cname,
				Sort:     sortIndex,
			})
			existMap[src2] = channelName
			sortIndex++
		}
	}
	log.Println("原始有效频道数量:", rawCount) // 新增日志输出
	dao.DB.Model(&models.IptvCategory{}).Where("name = ?", cname).Update("quantity", rawCount)

	// 批量删除数据库中当前分类但新列表中没有的 URL
	for _, ch := range oldChannels {
		if _, ok := srclistUrls[ch.Url]; !ok {
			delIDs = append(delIDs, ch.ID)
		}
	}

	// 在事务中执行删除和新增
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if len(delIDs) > 0 {
			if err := tx.Delete(&models.IptvChannel{}, delIDs).Error; err != nil {
				return err
			}
		}
		if len(newChannels) > 0 {
			if err := tx.Create(&newChannels).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return repetNum, err
	}

	// 只有当有新增或删除时才执行异步更新
	if len(newChannels) > 0 || len(delIDs) > 0 {
		go BindChannel()
		go until.CleanMealsXmlCacheAll()
		go until.ClearAutoChannelCache()
	}
	return repetNum, nil
}

func CategoryChangeStatus(params url.Values) dto.ReturnJsonDto {
	name := params.Get("change_status")
	if name == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "Category name不能为空", Type: "danger"}
	}

	var cateData models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).First(&cateData).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询Category失败", Type: "danger"}
	}

	if cateData.Enable == 1 {
		dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).Update("enable", 0)
		dao.DB.Model(&models.IptvCategory{}).Where("name like ?", "%("+name+")").Update("enable", 0)
	} else {
		dao.DB.Model(&models.IptvCategory{}).Where("name = ?", name).Update("enable", 1)
		dao.DB.Model(&models.IptvCategory{}).Where("name like ?", "%("+name+")").Update("enable", 1)
	}
	go until.CleanMealsXmlCacheAll() // 清除缓存
	go until.ClearAutoChannelCache() // 清除缓存
	return dto.ReturnJsonDto{Code: 1, Msg: "Category " + cateData.Name + "状态修改成功", Type: "success"}
}

func UpdateListAll() dto.ReturnJsonDto {
	if crontab.UpdateStatus {
		return dto.ReturnJsonDto{Code: 0, Msg: "后台更新中", Type: "danger"}
	}

	crontab.UpdateStatus = true
	defer func() { crontab.UpdateStatus = false }()

	go crontab.UpdateList() // 更新所有频道列表
	return dto.ReturnJsonDto{Code: 1, Msg: "开始后台更新", Type: "success"}
}

func UploadPayList(c *gin.Context) dto.ReturnJsonDto {
	file, err := c.FormFile("paylistfile")
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取文件失败:" + err.Error(), Type: "danger"}
	}

	listName := "文件导入" + time.Now().Format("20060102")

	f, err := file.Open()
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "打开文件失败: " + err.Error(), Type: "danger"}
	}
	defer f.Close()

	// 读取内容
	data, err := io.ReadAll(f)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "读取文件失败: " + err.Error(), Type: "danger"}
	}

	// 转为字符串
	urlData := until.FilterEmoji(string(data)) // 过滤emoji表情

	if until.IsM3UContent(urlData) {
		urlData = until.M3UToGenreTXT(urlData)
	}

	if !strings.Contains(urlData, "#genre#") {
		repeat, err := AddChannelList(listName, urlData, false)
		if err == nil && repeat != -1 {
			var maxSort int64
			dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
			dao.DB.Model(&models.IptvCategory{}).Create(&models.IptvCategory{Name: listName, Type: "user", Sort: maxSort + 1, LatestTime: time.Now().Format("2006-01-02 15:04:05")})
			return dto.ReturnJsonDto{Code: 1, Msg: fmt.Sprintf("更新列表 %s 成功，重复 %d 条\n", listName, repeat), Type: "success"}
		} else {
			return dto.ReturnJsonDto{Code: 0, Msg: fmt.Sprintf("更新列表 %s 失败\n", listName), Type: "danger"}
		}
	}
	return GenreChannels(listName, urlData, false)
}
