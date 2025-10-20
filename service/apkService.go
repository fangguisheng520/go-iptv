package service

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Getver() dto.GetverRes {
	var res dto.GetverRes

	var cfg = dao.GetConfig()

	res.AppVer = cfg.Build.Version
	res.UpSets = cfg.App.Update.Set
	res.UpText = cfg.App.Update.Text
	res.AppURL = cfg.ServerUrl + "/app/" + cfg.Build.Name + ".apk"
	res.UpSize = until.GetFileSize("./app/" + cfg.Build.Name + ".apk")
	return res
}

func GetBg() string {
	// 获取指定目录下的所有png文件
	dir := "/config/images/bj"
	files, err := filepath.Glob(filepath.Join(dir, "*.png"))
	if err != nil {
		return ""
	}
	if len(files) == 0 {
		return ""
	}

	pngs := make([]string, len(files))
	for i, file := range files {
		pngs[i] = filepath.Base(file)
	}
	randomIndex := rand.Intn(len(pngs))
	return pngs[randomIndex]
}

func ApkLogin(user models.IptvUser) dto.LoginRes {

	var result dto.LoginRes

	var cfg = dao.GetConfig()

	result.IP = user.IP
	result.ID = user.Name
	result.Status = user.Status
	result.NetType = user.NetType
	result.Location = user.Region

	result.ShowInterval = cfg.Channel.Interval
	result.AdText = cfg.Ad.AdText
	result.Decoder = cfg.App.Decoder
	result.AppVer = cfg.Build.Version
	result.AutoUpdate = cfg.Channel.Auto
	result.UpdateInterval = cfg.Channel.Interval
	result.BuffTimeOut = cfg.App.BuffTimeout
	result.TipLoading = cfg.Tips.Loading
	result.DataURL = cfg.ServerUrl + "/apk/channels"
	result.AppURL = cfg.ServerUrl + "/app/" + cfg.Build.Name + ".apk"
	result.ShowTime = cfg.Ad.ShowTime
	result.TipUserNoReg = "当前账号 " + strconv.FormatInt(user.Name, 10) + " " + cfg.Tips.UserNoReg
	result.TipUserExpired = "当前账号 " + strconv.FormatInt(user.Name, 10) + " " + cfg.Tips.UserExpired
	result.TipUserForbidden = "当前账号 " + strconv.FormatInt(user.Name, 10) + " " + cfg.Tips.UserForbidden
	result.AdInfo = "作者博客: www.qingh.xyz"
	result.RandKey = until.Md5(time.Now().Format("20060102150405") + strconv.FormatInt(user.Name, 10))

	return getUserInfo(user, result)
}

func GetChannels(channel dto.DataReqDto) string {
	resList := []dto.ChannelListDto{{
		Name: "我的收藏",
		Data: []dto.ChannelData{},
		Tmp:  "6L+Z5Y+q5piv5Y2g5L2N77yM5LiN54S25rKh6aKR6YGT5a655piT5Ye6546w6ZSZ6K+v77yM5LirYXBr5Yqg5a+G5pWw5o2u6ZyA6KaB5YWI5YigMTI45a2X6IqC77yM5rKh6aKR6YGT5bCx5LiN5aSfMTI4",
	}}

	var dbUser models.IptvUser
	err := dao.DB.Where("mac = ?", channel.Mac).First(&dbUser).Error
	if err != nil {

		resList = append(resList, dto.ChannelListDto{})
	}

	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	userExp := int64(until.DiffDays(todayZero.Unix(), dbUser.Exp))
	if userExp <= 0 {
		resList = append(resList, dto.ChannelListDto{})
	}

	var meal models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Where("status = ? and id = ?", 1, dbUser.Meal).First(&meal)
	cList := strings.Split(meal.Content, "_")

	var channelList []models.IptvChannel

	if len(cList) > 1 || (len(cList) == 1 && cList[0] != "") {
		dao.DB.Model(&models.IptvChannel{}).Where("category in ?", cList).Order("sort asc").Find(&channelList)
	} else {
		resList = append(resList, dto.ChannelListDto{
			Name: "该套餐无频道",
		})

		jsonData, _ := json.Marshal(resList)
		jsonStr := until.DecodeUnicode(string(jsonData))
		return encrypt(jsonStr, channel.Rand)
	}

	var categoryList []models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("name in ? and enable = ? and type != ?", cList, 1, "import").Order("sort asc").Find(&categoryList)

	for _, v := range categoryList {
		var tmpData []dto.ChannelData
		var i int64 = 1
		var dataMap = make(map[string][]string)
		var tmpMap = make(map[string]int64)

		if v.Sort == -2 {
			text := until.GetCCTVChannelList(false)
			text = strings.TrimSpace(text)     // 去掉结尾多余换行
			parts := strings.Split(text, "\n") // 按换行符分割
			if len(parts) == 1 && parts[0] == "" {
				resList = append(resList, dto.ChannelListDto{
					ID:   int64(v.Sort + 3),
					Name: v.Name,
					Data: tmpData,
				})
				continue
			}
			for _, part := range parts {
				data := strings.Split(part, ",")
				dataMap[data[0]] = append(dataMap[data[0]], strings.TrimSpace(data[1]))
				if _, ok := tmpMap[data[0]]; !ok {
					tmpMap[data[0]] = i
					i++
				}
			}
		} else if v.Sort == -1 {
			text := until.GetProvinceChannelList(false)
			text = strings.TrimSpace(text)     // 去掉结尾多余换行
			parts := strings.Split(text, "\n") // 按换行符分割
			if len(parts) == 1 && parts[0] == "" {
				resList = append(resList, dto.ChannelListDto{
					ID:   int64(v.Sort + 3),
					Name: v.Name,
					Data: tmpData,
				})
				continue
			}
			for _, part := range parts {
				data := strings.Split(part, ",")
				dataMap[data[0]] = append(dataMap[data[0]], strings.TrimSpace(data[1]))
				if _, ok := tmpMap[data[0]]; !ok {
					tmpMap[data[0]] = i
					i++
				}
			}
		} else {
			for _, channel := range channelList {
				if v.Name == channel.Category {
					dataMap[channel.Name] = append(dataMap[channel.Name], strings.TrimSpace(channel.Url))
					if _, ok := tmpMap[channel.Name]; !ok {
						tmpMap[channel.Name] = i
						i++
					}
				}
			}
		}

		for k, v1 := range tmpMap {
			tmpData = append(tmpData, dto.ChannelData{
				Num:    v1,
				Name:   k,
				Source: dataMap[k],
			})
		}

		sort.Slice(tmpData, func(i, j int) bool {
			return tmpData[i].Num < tmpData[j].Num
		})

		resList = append(resList, dto.ChannelListDto{
			ID:   int64(v.Sort + 3),
			Name: v.Name,
			Data: tmpData,
		})
	}
	sort.Slice(resList, func(i, j int) bool {
		return resList[i].ID < resList[j].ID
	})
	jsonData, _ := json.Marshal(resList)
	jsonStr := until.DecodeUnicode(string(jsonData))

	return encrypt(jsonStr, channel.Rand)
}

func encrypt(str string, randkey string) string {
	encoded, _ := CompressString(str)

	// Step 2: MD5 加密 key

	hashedKey := until.Md5(until.GetAesKey() + randkey)

	// Step 3: 截取 hashedKey 的一部分
	subKey := hashedKey[7:23]

	// Step 3: AES 加密
	aes := until.NewAes(subKey, "AES-128-ECB", "")
	encrypted, err := aes.Encrypt(encoded)

	if err != nil {
		return ""
	}

	// Step 4: 替换字符
	// encrypted := string(ciphertext)
	encrypted = strings.ReplaceAll(encrypted, "f", "&")
	encrypted = strings.ReplaceAll(encrypted, "b", "f")
	encrypted = strings.ReplaceAll(encrypted, "&", "b")
	encrypted = strings.ReplaceAll(encrypted, "t", "#")
	encrypted = strings.ReplaceAll(encrypted, "y", "t")
	encrypted = strings.ReplaceAll(encrypted, "#", "y")

	// Step 5: 反转和截取
	start := 44
	length := 128
	end := start + length

	// 防止越界
	if end > len(encrypted) {
		end = len(encrypted)
	}

	coded := encrypted[start:end]
	reversed := until.ReverseString(coded)
	finalEncrypted := reversed + encrypted

	return finalEncrypted
}

func CompressString(input string) (string, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	_, err := w.Write([]byte(input))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func getUserInfo(user models.IptvUser, result dto.LoginRes) dto.LoginRes {
	var cfg = dao.GetConfig()

	var movie []models.IptvMovie
	dao.DB.Model(&models.IptvMovie{}).Where("state = ?", 1).Order("id desc").Find(&movie)

	result.MovieEngine.Model = movie

	if cfg.App.NeedAuthor == 0 {
		result = getMealName(user, result)
		log.Printf("用户: %d 登录成功,IP: %s 设备ID: %s 套餐: %s \n", result.ID, result.IP, user.DeviceID, result.MealName)
	} else if cfg.App.NeedAuthor == 1 && user.Status != -1 {
		result = getMealName(user, result)
		log.Printf("用户: %d 登录成功,IP: %s 设备ID: %s 套餐: %s\n", result.ID, result.IP, user.DeviceID, result.MealName)
	} else {
		log.Printf("用户: %d 登录成功,IP: %s 设备ID: %s 未授权 \n", result.ID, result.IP, user.DeviceID)
	}

	return result
}

func getMealName(user models.IptvUser, result dto.LoginRes) dto.LoginRes {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Where("status = ?", 1).Find(&meals)

	for _, v := range meals {
		if v.ID == 1000 && result.MealName == "" {
			result.MealName = v.Name
			result.ProvList = strings.Split(v.Content, "_")
		}
		if v.ID == user.Meal {
			result.MealName = v.Name
			result.ProvList = strings.Split(v.Content, "_")
		}
	}
	return result
}

func CheckUserDb(user dto.ApkUser, ip string) models.IptvUser {
	var dbUser models.IptvUser
	res := dao.DB.Where("mac = ?", user.Mac).Find(&dbUser)
	if res.RowsAffected == 0 {
		return AddUser(user, ip)
	}

	dbUser.LastTime = time.Now().Unix()
	dbUser.IP = ip
	dbUser.Region = until.GetIpRegion(user.IP)
	dbUser.NetType = user.NetType
	dbUser.DeviceID = user.DeviceID

	dao.DB.Model(&models.IptvUser{}).Where("mac = ?", user.Mac).Updates(dbUser)

	return dbUser
}

func AddUser(user dto.ApkUser, ip string) models.IptvUser {
	user.IP = ip
	user.Region = until.GetIpRegion(ip)
	var cfg = dao.GetConfig()

	dbData := models.IptvUser{
		Name:     int64(genName()),
		Mac:      user.Mac,
		DeviceID: user.DeviceID,
		Model:    user.Model,
		IP:       user.IP,
		Region:   user.Region,
		LastTime: time.Now().Unix(),
		Meal:     1000,
	}

	switch cfg.App.NeedAuthor {
	case 0:
		dbData.Status = 999
		dbData.Marks = "自动授权"
		dbData.Exp = 0
	case 1:
		dbData.Status = -1
		dbData.Marks = "未授权"
		dbData.Exp = 0
	}

	dao.DB.Model(&models.IptvUser{}).Create(&dbData)
	return dbData
}

func genName() int {
	name := rand.Intn(999999-1000+1) + 1000 // 生成 1000~999999 之间的随机数
	var count int64
	err := dao.DB.Model(&models.IptvUser{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return name
	} else {
		return genName() // 递归调用
	}
}
