# 清和iptv
>提取自矿神群晖IPTV并大改，由原来的PHP+MySql改为Go+Sqlite     
>添加缺失功能，精简删除非必要功能，修改系统存在的安全漏洞   


# 注意
当前版本与之前PHP版本并不兼容，若要使用PHP版本，请使用`docker pull v1st233/iptv:20250905`

## Change log
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