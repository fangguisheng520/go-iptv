# 清和iptv
>提取自矿神群晖IPTV并大改，由原来的PHP+MySql改为Go+Sqlite     
>添加缺失功能，精简删除非必要功能，修改系统存在的安全漏洞   


# 注意
当前版本与之前PHP版本并不兼容，若要使用PHP版本，请使用`docker pull v1st233/iptv:20250905`

## Change log
#### 2025-10-13
- 修改epg订阅链接格式
- 添加频道文件导入，支持txt和m3u格式
- 添加epg xml格式epg导入  暂无对应的epg订阅输出

#### 2025-10-12
- 修复epg订阅token错误
- 修复epg Channel错误

#### 2025-10-11
- 添加m3u订阅epg支持
- 修改内网穿透模式无法访问logo (感谢KIKI协助)
- 修复apk epg为空bug
- 改用aes替代jwt，解决订阅链接过长问题
- 修复频道保存bug，默认频道添加修改错误提示
- 更新依赖

#### 2025-10-10
- 修改客户端修改授权显示未授权bug
- 套餐添加txt、m3u格式订阅，用以支持酷9、tvbox等客户端 (暂时只有51zmt的epg，后续会添加更多； 台标暂时只能手动添加到`/config/logo`，后续会添加页面上传功能)

#### 2025-10-9
- 修复点播保存不合法问题
- 精简菜单，精简配置项
- 修复背景图片显示逻辑

#### 2025-9-29
- 修复弹窗提交内容时滚动条失效bug
- 修复时区问题
- 添加客户端是否需要授权开关(默认需要)
- 清理日志，添加登录提示

#### 2025-9-28
- 修复EXTM3U直播源无法保存bug
- 修复登录超时的跳转体验
- 修复定时更新外部链接功能的定时任务重复执行
- 添加epg缓存
- epg暂时只有cntv和51zmt 够我用2333，求打赏给动力      

#### 2025-9-26
- 未安装时全局跳转到安装页面
- 添加编译时无法下载

#### 2025-9-25
- go重构管理页面
- 改为sqlite，更清晰明了的文件映射
- 添加了自动安装及友好的安装提示
- 更友好的页面加载体验，专注家庭使用，删除了订单相关功能
- 更简洁的操作体验

#### 2025-9-18
- 改不动这个史了，归档，go重构了下，看这个吧 [go-iptv](https://github.com/wz1st/go-iptv)

#### 2025-9-5
- 修复了文件上传漏洞
- 修复了任意文件删除漏洞
- 添加定时更新外部列表
- 添加更改应用图标功能
- 添加自动重新编译功能
- 添加修改应用名称、包名、签名key功能
- 修改了系统图标、系统名称、系统版本

#### 2025-8-25
- 添加docker自动构建，添加armv7、arm64、386、amd64版本

#### 2025-8-22
- 修复了SQL注入漏洞
- 改为alpine+nginx+php-fpm+mariadb 精简镜像大小

## 安装
```
docker volume create iptv
docker pull v1st233/iptv:latest
docker run -d --name iptv_server -p <port>:80 -v iptv:/config v1st233/iptv:latest
# username: admin
# password: password
```
或
```
git clone https://github.com/wz1st/go-iptv.git
cd iptv
docker build -f Dockerfile -t image_name:latest .
docker volume create iptv
docker run -d --name iptv_server -p port:80 -v iptv:/config image_name:latest
``` 
## 使用
容器跑起来后访问`http://<ip>:<port>`即可，根据提示安装系统，然后登录添加源->修改套餐->下载安装APK->授权用户即可使用

## 鸣谢
- [我不是矿神](https://imnks.com/)
## 打赏
>如果觉得好用，请打赏支持一下

<div style="display: flex; justify-content: center; gap: 50px;">
  <img src="./static/images/wxpay.jpg" alt="微信" width="300">
  <img src="./static/images/zfbpay.jpg" alt="支付宝" width="300">
</div>



## 小声哔哔
>本程序仅供学习交流使用，请勿用于商业用途，否则后果自负。     
>本程序不保证长期稳定运行，请自行备份。     
>源自己找，有问题自己解决。     
<a id="bottom"></a> 