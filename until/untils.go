package until

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/net/publicsuffix"
)

func Md5(str string) (retMd5 string) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func GetUrl(c *gin.Context) string {
	protocol := "http://"
	if c.Request.TLS != nil ||
		c.Request.Header.Get("X-Forwarded-Proto") == "https" ||
		c.Request.Header.Get("Front-End-Https") == "on" {
		protocol = "https://"
	}
	// 获取主机名
	host := c.Request.Host
	// 获取请求URI
	requestURI := c.Request.RequestURI
	parts := strings.Split(requestURI, "/")
	if len(parts) > 1 {
		parts = parts[:len(parts)-1]
	}
	modifiedPath := strings.Join(parts, "/")
	// 构建完整URL
	return protocol + host + modifiedPath
}

func GetUrlData(url string, ua ...string) string {
	defaultUA := "Go-http-client/1.1"
	useUA := defaultUA

	if len(ua) > 0 && ua[0] != "" {
		useUA = ua[0]
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", useUA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func GetIpRegion(ip string) string {
	city := "局域网"
	if !isPrivateIP(net.ParseIP(ip)) {
		return city
	}
	url := "https://api.mir6.com/api/ip_json?ip=" + ip
	jsonStr := GetUrlData(url)
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return city
	}
	if data, ok := jsonMap["data"].(map[string]interface{}); ok {
		if cityStr, ok := data["location"].(string); ok && cityStr != "" {
			city = cityStr
		}
	}
	return city
}

func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		{
			IP:   net.IPv4(10, 0, 0, 0),
			Mask: net.CIDRMask(8, 32),
		},
		{
			IP:   net.IPv4(172, 16, 0, 0),
			Mask: net.CIDRMask(12, 32),
		},
		{
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.CIDRMask(16, 32),
		},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func DecodeUnicode(s string) string {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		hex := re.FindStringSubmatch(match)[1]
		codePoint, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return match
		}
		return string(rune(codePoint))
	})
}

func GetFileSize(filePath string) string {

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "0 MB"
	}

	// 获取文件大小（字节）
	fileSize := fileInfo.Size()

	// 将文件大小转换为兆字节 (MB)
	fileSizeMB := float64(fileSize) / (1024 * 1024)

	// 输出文件大小（MB）
	return fmt.Sprintf("%.2f MB", fileSizeMB)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CopyFile(from, to string) error {
	// 打开源文件
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件（如果存在会覆盖）
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 复制内容
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// 同步写入磁盘
	return dst.Sync()
}

func CheckJava() bool {
	log.Println("检查Java版本...")
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Java版本检查失败:", err)
		return false
	}

	// 解析输出结果
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// 输出 Java 版本信息
	if len(lines) > 0 {
		javaVersion := lines[0]
		log.Println("Java版本:", javaVersion)

		// 判断 Java 版本是否为 1.8
		if strings.Contains(javaVersion, "1.8") {
			return true
		} else {
			log.Println("Java版本不是 1.8")
			return false
		}
	} else {
		log.Println("无法确定 Java 版本")
		return false
	}
}

func CheckApktool() bool {
	log.Println("检查Apktool版本...")
	cmd := exec.Command("apktool", "-version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Apktool检查失败:", err)
		return false
	}

	// 解析输出结果
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// 输出 Apktool 版本信息
	if len(lines) > 0 {
		apktoolVersion := lines[0]
		log.Println("Apktool版本:", apktoolVersion)
		return true
	} else {
		log.Println("无法确定 Apktool 版本")
		return false
	}
}

func CheckPort(port string) bool {
	log.Println("检查端口占用...")
	// 尝试监听给定端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("端口" + port + "被占用...")
		return false // 端口被占用
	}
	listener.Close() // 关闭监听器
	return true      // 端口未被占用
}

func FilterEmoji(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if utf8.RuneLen(r) < 4 {
			result = append(result, r)
		}
	}
	return string(result)
}

func GetPngFileNames(dir string) ([]string, error) {
	var names []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".png") {
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			names = append(names, name)
		}
	}

	return names, nil
}

func IsSafeImgName(inputPath string) bool {
	// 清理输入路径
	imgName := filepath.Clean(inputPath)

	// 检查是否包含 ..
	if strings.Contains(imgName, "..") {
		return false
	}

	// 检查是否包含 .
	if strings.Contains(imgName, ".") {
		return false
	}

	// 检查是否是绝对路径
	if filepath.IsAbs(imgName) {
		return false
	}

	// 检查是否包含可疑字符
	if strings.ContainsAny(imgName, `:*?"<>|'`) {
		return false
	}
	pattern := `^[a-fA-F0-9]{32}$`
	match, _ := regexp.MatchString(pattern, imgName)
	return match
}

func IsSafe(input string) bool {

	if input == "" {
		return true
	}

	// 检查是否包含可疑字符
	if strings.ContainsAny(input, `:*?"<>|'+`) {
		return false
	}

	return true
}

func GetFileModTimeStr(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	// 使用自定义格式 YYYY.MM.DD
	return info.ModTime().Format("2006.01.02"), nil
}

func DiffDays(ts1, ts2 int64) int {
	// 转换为 time.Time
	t1 := time.Unix(ts1, 0)
	t2 := time.Unix(ts2, 0)

	// 计算差值
	diff := t2.Sub(t1)
	if diff < 0 {
		diff = -diff // 保证正数
	}

	days := diff.Hours() / 24

	return int(math.Ceil(days))
}

func GetContainerID() (string, error) {
	hostname, err := os.ReadFile("/etc/hostname")
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(hostname))
	return id, nil
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

func CheckLogo(path string) (bool, error) {
	// 检查路径是否存在
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil // 不存在
	}
	if err != nil {
		return false, err // 其他错误
	}

	// 判断是否是目录
	if !info.IsDir() {
		return false, fmt.Errorf("%s 不是目录", path)
	}

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// 读取目录内容（最多读一个即可判断是否为空）
	entries, err := dir.ReadDir(1)
	if err != nil && err != fs.ErrClosed {
		return false, err
	}

	// 如果有内容返回 true，否则 false
	if len(entries) > 0 {
		return true, nil
	}

	return false, nil
}

func GetLogos() []string {
	// 获取指定目录下的所有png文件
	dir := "/config/logo"
	files, err := filepath.Glob(filepath.Join(dir, "*.png"))
	if err != nil {
		return []string{}
	}
	if len(files) == 0 {
		return []string{}
	}

	pngs := make([]string, len(files))
	for i, file := range files {
		pngs[i] = filepath.Base(file)
	}
	return pngs
}

func EpgNameGetLogo(eNmae string) string {
	// 获取指定目录下的所有png文件
	dir := "/config/logo"
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

	epgName := strings.SplitN(eNmae, "-", 2)[1]
	for _, logo := range pngs {
		logoName := strings.Split(logo, ".")[0]
		if strings.EqualFold(epgName, logoName) {
			return "/logo/" + logo
		}
	}
	return ""
}

func IsValidHost(host string) bool {
	if host == "" {
		return false
	}

	// 去掉端口
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err == nil {
			host = h
		}
	}

	// 如果是 IP
	if net.ParseIP(host) != nil {
		return true
	}

	// 匹配域名
	domainPattern := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(domainPattern, host)
	return matched
}

func GetMainDomain(input string) string {
	// 补全 URL 协议前缀，避免解析失败
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}

	u, err := url.Parse(input)
	if err != nil {
		return ""
	}

	domain := u.Hostname()
	if domain == "" {
		return ""
	}

	// 获取主域 + 后缀（如 51zmt.top / bbc.co.uk）
	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return ""
	}

	// 获取后缀（如 top / co.uk）
	suffix, _ := publicsuffix.PublicSuffix(domain)

	// 去掉后缀部分，保留主域
	mainDomain := strings.TrimSuffix(eTLDPlusOne, "."+suffix)
	return mainDomain
}

func MergeAndUnique(a, b []string) []string {
	m := make(map[string]struct{})

	for _, v := range append(a, b...) {
		// 去掉空字符串
		if v == "" {
			continue
		}
		// 去重
		if _, exists := m[v]; !exists {
			m[v] = struct{}{}
		}
	}

	// 转回切片
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func GetVersion() string {
	data, err := os.ReadFile("/app/Version")
	if err != nil {
		return "2.3.0"
	}
	return strings.TrimSpace(string(data))
}

func EqualStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := append([]string(nil), a...)
	bCopy := append([]string(nil), b...)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func Int64InStringSlice(target int64, list []string) bool {
	s := strconv.FormatInt(target, 10) // int64 转字符串
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func InStringSlice(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func CheckRam() bool {
	// 判断可用内存
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return false
	}
	log.Printf("可用内存: %d MB (%d GB)\n", vmStat.Available/1024/1024, vmStat.Available/1024/1024/1024)
	return vmStat.Available < 256*1024*1024
}

func IsLowResource() bool {
	// 判断 ARM 架构
	if runtime.GOARCH == "arm" {
		return true
	}

	// 判断 CPU 核心数
	if runtime.NumCPU() < 2 {
		return true
	}

	// 判断可用内存
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return false
	}
	return vmStat.Available < 256*1024*1024
}

func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
