PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE iptv_admin (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE iptv_category (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    enable INTEGER NOT NULL DEFAULT 1,
    type TEXT NOT NULL DEFAULT 'add',
    url TEXT DEFAULT NULL,
    autocategory TEXT DEFAULT NULL,
    latesttime TEXT DEFAULT NULL,
    sort INTEGER
);
INSERT INTO iptv_category VALUES(1,'央视频道(自动聚合)',1,'add','',0,'',-2);
INSERT INTO iptv_category VALUES(2,'卫视频道(自动聚合)',1,'add','',0,'',-1);
CREATE TABLE iptv_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    url TEXT DEFAULT NULL,
    category TEXT NOT NULL,
    sort INTEGER
);
CREATE TABLE iptv_epg (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    content TEXT DEFAULT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    remarks TEXT DEFAULT NULL
);
INSERT INTO iptv_epg VALUES(1,'cntv-CCTV1','CCTV1',1,'CCTV1|CCTV-1');
INSERT INTO iptv_epg VALUES(2,'cntv-CCTV2','CCTV2',1,'CCTV2|CCTV-2');
INSERT INTO iptv_epg VALUES(3,'cntv-CCTV3','CCTV3',1,'CCTV3|CCTV-3');
INSERT INTO iptv_epg VALUES(4,'cntv-CCTV4','CCTV4',1,'CCTV4|CCTV-4');
INSERT INTO iptv_epg VALUES(5,'cntv-CCTV5','CCTV5',1,'CCTV5|CCTV-5');
INSERT INTO iptv_epg VALUES(6,'cntv-CCTV5+','CCTV5+',1,'CCTV5+|CCTV-5+');
INSERT INTO iptv_epg VALUES(7,'cntv-CCTV6','CCTV6',1,'CCTV6|CCTV-6');
INSERT INTO iptv_epg VALUES(8,'cntv-CCTV7','CCTV7',1,'CCTV7|CCTV-7');
INSERT INTO iptv_epg VALUES(9,'cntv-CCTV8','CCTV8',1,'CCTV8|CCTV-8');
INSERT INTO iptv_epg VALUES(10,'cntv-CCTV9','CCTV9',1,'CCTV9|CCTV-9');
INSERT INTO iptv_epg VALUES(11,'cntv-CCTV10','CCTV10',1,'CCTV10|CCTV-10');
INSERT INTO iptv_epg VALUES(12,'cntv-CCTV11','CCTV11',1,'CCTV11|CCTV-11');
INSERT INTO iptv_epg VALUES(13,'cntv-CCTV12','CCTV12',1,'CCTV12|CCTV-12');
INSERT INTO iptv_epg VALUES(14,'cntv-CCTV13','CCTV13',1,'CCTV13|CCTV-13');
INSERT INTO iptv_epg VALUES(15,'cntv-CCTV14','CCTV14',1,'CCTV14|CCTV-14');
INSERT INTO iptv_epg VALUES(16,'cntv-CCTV15','CCTV15',1,'CCTV15|CCTV-15');
INSERT INTO iptv_epg VALUES(17,'cntv-CCTV16','CCTV16',1,'CCTV16|CCTV-16');
INSERT INTO iptv_epg VALUES(18,'cntv-CCTV17','CCTV17',1,'CCTV17|CCTV-17');
INSERT INTO iptv_epg VALUES(19,'51zmt-CCTV1','','1','CCTV1|CCTV-1');
INSERT INTO iptv_epg VALUES(20,'51zmt-CCTV2','','1','CCTV2|CCTV-2');
INSERT INTO iptv_epg VALUES(21,'51zmt-CCTV3','','1','CCTV3|CCTV-3');
INSERT INTO iptv_epg VALUES(22,'51zmt-CCTV4','','1','CCTV4|CCTV-4');
INSERT INTO iptv_epg VALUES(23,'51zmt-CCTV5','','1','CCTV5|CCTV-5');
INSERT INTO iptv_epg VALUES(24,'51zmt-CCTV5+','','1','CCTV5+|CCTV-5+');
INSERT INTO iptv_epg VALUES(25,'51zmt-CCTV6','','1','CCTV6|CCTV-6');
INSERT INTO iptv_epg VALUES(26,'51zmt-CCTV7','','1','CCTV7|CCTV-7');
INSERT INTO iptv_epg VALUES(27,'51zmt-CCTV8','','1','CCTV8|CCTV-8');
INSERT INTO iptv_epg VALUES(28,'51zmt-CCTV9','','1','CCTV9|CCTV-9');
INSERT INTO iptv_epg VALUES(29,'51zmt-CCTV10','','1','CCTV10|CCTV-10');
INSERT INTO iptv_epg VALUES(30,'51zmt-CCTV11','','1','CCTV11|CCTV-11');
INSERT INTO iptv_epg VALUES(31,'51zmt-CCTV12','','1','CCTV12|CCTV-12');
INSERT INTO iptv_epg VALUES(32,'51zmt-CCTV13','','1','CCTV13|CCTV-13');
INSERT INTO iptv_epg VALUES(33,'51zmt-CCTV14','','1','CCTV14|CCTV-14');
INSERT INTO iptv_epg VALUES(34,'51zmt-CCTV15','','1','CCTV15|CCTV-15');
INSERT INTO iptv_epg VALUES(35,'51zmt-CCTV16','','1','CCTV16|CCTV-16');
INSERT INTO iptv_epg VALUES(36,'51zmt-CCTV17','','1','CCTV17|CCTV-17');
INSERT INTO iptv_epg VALUES(37,'51zmt-CGTN','','1','CGTN');
INSERT INTO iptv_epg VALUES(38,'51zmt-CCTV4EUO','','1','CCTV4EUO');
INSERT INTO iptv_epg VALUES(39,'51zmt-CCTV4AME','','1','CCTV4AME');
INSERT INTO iptv_epg VALUES(40,'51zmt-湖南卫视','','1','湖南卫视');
INSERT INTO iptv_epg VALUES(41,'51zmt-浙江卫视','','1','浙江卫视');
INSERT INTO iptv_epg VALUES(42,'51zmt-江苏卫视','','1','江苏卫视');
INSERT INTO iptv_epg VALUES(43,'51zmt-北京卫视','','1','北京卫视');
INSERT INTO iptv_epg VALUES(44,'51zmt-东方卫视','','1','东方卫视');
INSERT INTO iptv_epg VALUES(45,'51zmt-安徽卫视','','1','安徽卫视');
INSERT INTO iptv_epg VALUES(46,'51zmt-广东卫视','','1','广东卫视');
INSERT INTO iptv_epg VALUES(47,'51zmt-深圳卫视','','1','深圳卫视');
INSERT INTO iptv_epg VALUES(48,'51zmt-辽宁卫视','','1','辽宁卫视');
INSERT INTO iptv_epg VALUES(49,'51zmt-旅游卫视','','1','旅游卫视');
INSERT INTO iptv_epg VALUES(50,'51zmt-山东卫视','','1','山东卫视');
INSERT INTO iptv_epg VALUES(51,'51zmt-天津卫视','','1','天津卫视');
INSERT INTO iptv_epg VALUES(52,'51zmt-重庆卫视','','1','重庆卫视');
INSERT INTO iptv_epg VALUES(53,'51zmt-东南卫视','','1','东南卫视');
INSERT INTO iptv_epg VALUES(54,'51zmt-甘肃卫视','','1','甘肃卫视');
INSERT INTO iptv_epg VALUES(55,'51zmt-广西卫视','','1','广西卫视');
INSERT INTO iptv_epg VALUES(56,'51zmt-贵州卫视','','1','贵州卫视');
INSERT INTO iptv_epg VALUES(57,'51zmt-河北卫视','','1','河北卫视');
INSERT INTO iptv_epg VALUES(58,'51zmt-黑龙江卫视','','1','黑龙江卫视');
INSERT INTO iptv_epg VALUES(59,'51zmt-河南卫视','','1','河南卫视');
INSERT INTO iptv_epg VALUES(60,'51zmt-湖北卫视','','1','湖北卫视');
INSERT INTO iptv_epg VALUES(61,'51zmt-江西卫视','','1','江西卫视');
INSERT INTO iptv_epg VALUES(62,'51zmt-吉林卫视','','1','吉林卫视');
INSERT INTO iptv_epg VALUES(63,'51zmt-内蒙古卫视','','1','内蒙古卫视');
INSERT INTO iptv_epg VALUES(64,'51zmt-宁夏卫视','','1','宁夏卫视');
INSERT INTO iptv_epg VALUES(65,'51zmt-山西卫视','','1','山西卫视');
INSERT INTO iptv_epg VALUES(66,'51zmt-陕西卫视','','1','陕西卫视');
INSERT INTO iptv_epg VALUES(67,'51zmt-四川卫视','','1','四川卫视');
INSERT INTO iptv_epg VALUES(68,'51zmt-新疆卫视','','1','新疆卫视');
INSERT INTO iptv_epg VALUES(69,'51zmt-云南卫视','','1','云南卫视');
INSERT INTO iptv_epg VALUES(70,'51zmt-青海卫视','','1','青海卫视');
INSERT INTO iptv_epg VALUES(71,'51zmt-兵团卫视','','1','兵团卫视');
INSERT INTO iptv_epg VALUES(72,'51zmt-哈哈炫动','','1','哈哈炫动');
INSERT INTO iptv_epg VALUES(73,'51zmt-延边卫视','','1','延边卫视');
INSERT INTO iptv_epg VALUES(74,'51zmt-黄河卫视','','1','黄河卫视');
INSERT INTO iptv_epg VALUES(75,'51zmt-卡酷动画','','1','卡酷动画');
INSERT INTO iptv_epg VALUES(76,'51zmt-厦门卫视','','1','厦门卫视');
INSERT INTO iptv_epg VALUES(77,'51zmt-金鹰卡通','','1','金鹰卡通');
INSERT INTO iptv_epg VALUES(78,'51zmt-康巴卫视','','1','康巴卫视');
INSERT INTO iptv_epg VALUES(79,'51zmt-西藏卫视','','1','西藏卫视');
INSERT INTO iptv_epg VALUES(80,'51zmt-三沙卫视','','1','三沙卫视');
INSERT INTO iptv_epg VALUES(81,'51zmt-中国教育1台','','1','中国教育1台');
INSERT INTO iptv_epg VALUES(82,'51zmt-中国教育2台','','1','中国教育2台');
INSERT INTO iptv_epg VALUES(83,'51zmt-中国教育3台','','1','中国教育3台');
INSERT INTO iptv_epg VALUES(84,'51zmt-3D电视试验频道','','1','3D电视试验频道');
INSERT INTO iptv_epg VALUES(85,'51zmt-外汇理财','','1','外汇理财');
INSERT INTO iptv_epg VALUES(86,'51zmt-电竞天堂','','1','电竞天堂');
INSERT INTO iptv_epg VALUES(87,'51zmt-IPTV5+','','1','IPTV5+');
INSERT INTO iptv_epg VALUES(88,'51zmt-IPTV6+','','1','IPTV6+');
INSERT INTO iptv_epg VALUES(89,'51zmt-IPTV经典电影','','1','IPTV经典电影');
INSERT INTO iptv_epg VALUES(90,'51zmt-IPTV热播剧场','','1','IPTV热播剧场');
INSERT INTO iptv_epg VALUES(91,'51zmt-IPTV少儿动画','','1','IPTV少儿动画');
INSERT INTO iptv_epg VALUES(92,'51zmt-IPTV魅力时尚','','1','IPTV魅力时尚');
INSERT INTO iptv_epg VALUES(93,'51zmt-CCTV4K','','1','CCTV4K');
INSERT INTO iptv_epg VALUES(94,'51zmt-GTV游戏竞技','','1','GTV游戏竞技');
INSERT INTO iptv_epg VALUES(95,'51zmt-新动漫','','1','新动漫');
INSERT INTO iptv_epg VALUES(96,'51zmt-山东齐鲁','','1','山东齐鲁');
INSERT INTO iptv_epg VALUES(97,'51zmt-山东体育','','1','山东体育');
INSERT INTO iptv_epg VALUES(98,'51zmt-山东农科','','1','山东农科');
INSERT INTO iptv_epg VALUES(99,'51zmt-山东公共','','1','山东公共');
INSERT INTO iptv_epg VALUES(100,'51zmt-山东少儿','','1','山东少儿');
INSERT INTO iptv_epg VALUES(101,'51zmt-山东影视','','1','山东影视');
INSERT INTO iptv_epg VALUES(102,'51zmt-山东综艺','','1','山东综艺');
INSERT INTO iptv_epg VALUES(103,'51zmt-山东生活','','1','山东生活');
INSERT INTO iptv_epg VALUES(104,'51zmt-环宇电影','','1','环宇电影');
INSERT INTO iptv_epg VALUES(105,'51zmt-湖北综合频道','','1','湖北综合频道');
INSERT INTO iptv_epg VALUES(106,'51zmt-湖北影视频道','','1','湖北影视频道');
INSERT INTO iptv_epg VALUES(107,'51zmt-湖北教育频道','','1','湖北教育频道');
INSERT INTO iptv_epg VALUES(108,'51zmt-湖北生活频道','','1','湖北生活频道');
INSERT INTO iptv_epg VALUES(109,'51zmt-湖北公共·新闻','','1','湖北公共·新闻');
INSERT INTO iptv_epg VALUES(110,'51zmt-湖北经济频道','','1','湖北经济频道');
INSERT INTO iptv_epg VALUES(111,'51zmt-湖北垄上频道','','1','湖北垄上频道');
INSERT INTO iptv_epg VALUES(112,'51zmt-武汉新闻综合频道','','1','武汉新闻综合频道');
INSERT INTO iptv_epg VALUES(113,'51zmt-武汉电视剧频道','','1','武汉电视剧频道');
INSERT INTO iptv_epg VALUES(114,'51zmt-武汉科技生活频道','','1','武汉科技生活频道');
INSERT INTO iptv_epg VALUES(115,'51zmt-武汉经济频道','','1','武汉经济频道');
INSERT INTO iptv_epg VALUES(116,'51zmt-武汉文体频道','','1','武汉文体频道');
INSERT INTO iptv_epg VALUES(117,'51zmt-武汉外语频道','','1','武汉外语频道');
INSERT INTO iptv_epg VALUES(118,'51zmt-武汉少儿频道','','1','武汉少儿频道');
INSERT INTO iptv_epg VALUES(119,'51zmt-武汉教育电视台','','1','武汉教育电视台');
INSERT INTO iptv_epg VALUES(120,'51zmt-北京纪实','','1','北京纪实');
INSERT INTO iptv_epg VALUES(121,'51zmt-BTV文艺','','1','BTV文艺');
INSERT INTO iptv_epg VALUES(122,'51zmt-BTV科教','','1','BTV科教');
INSERT INTO iptv_epg VALUES(123,'51zmt-BTV影视','','1','BTV影视');
INSERT INTO iptv_epg VALUES(124,'51zmt-BTV财经','','1','BTV财经');
INSERT INTO iptv_epg VALUES(125,'51zmt-BTV体育','','1','BTV体育');
INSERT INTO iptv_epg VALUES(126,'51zmt-BTV生活','','1','BTV生活');
INSERT INTO iptv_epg VALUES(127,'51zmt-BTV新闻','','1','BTV新闻');
INSERT INTO iptv_epg VALUES(128,'51zmt-SCTV2','','1','SCTV2');
INSERT INTO iptv_epg VALUES(129,'51zmt-SCTV3','','1','SCTV3');
INSERT INTO iptv_epg VALUES(130,'51zmt-SCTV4','','1','SCTV4');
INSERT INTO iptv_epg VALUES(131,'51zmt-SCTV5','','1','SCTV5');
INSERT INTO iptv_epg VALUES(132,'51zmt-SCTV7','','1','SCTV7');
INSERT INTO iptv_epg VALUES(133,'51zmt-峨嵋电影','','1','峨嵋电影');
INSERT INTO iptv_epg VALUES(134,'51zmt-SCTV9','','1','SCTV9');
INSERT INTO iptv_epg VALUES(135,'51zmt-SCTV8','','1','SCTV8');
INSERT INTO iptv_epg VALUES(136,'51zmt-CDTV1','','1','CDTV1');
INSERT INTO iptv_epg VALUES(137,'51zmt-CDTV2','','1','CDTV2');
INSERT INTO iptv_epg VALUES(138,'51zmt-CDTV3','','1','CDTV3');
INSERT INTO iptv_epg VALUES(139,'51zmt-CDTV4','','1','CDTV4');
INSERT INTO iptv_epg VALUES(140,'51zmt-CDTV5','','1','CDTV5');
INSERT INTO iptv_epg VALUES(141,'51zmt-CDTV6','','1','CDTV6');
INSERT INTO iptv_epg VALUES(142,'51zmt-风云足球','','1','风云足球');
INSERT INTO iptv_epg VALUES(143,'51zmt-辽宁都市','','1','辽宁都市');
INSERT INTO iptv_epg VALUES(144,'51zmt-辽宁影视剧','','1','辽宁影视剧');
INSERT INTO iptv_epg VALUES(145,'51zmt-辽宁青少','','1','辽宁青少');
INSERT INTO iptv_epg VALUES(146,'51zmt-辽宁生活','','1','辽宁生活');
INSERT INTO iptv_epg VALUES(147,'51zmt-辽宁公共','','1','辽宁公共');
INSERT INTO iptv_epg VALUES(148,'51zmt-辽宁北方','','1','辽宁北方');
INSERT INTO iptv_epg VALUES(149,'51zmt-辽宁体育','','1','辽宁体育');
INSERT INTO iptv_epg VALUES(150,'51zmt-辽宁经济','','1','辽宁经济');
INSERT INTO iptv_epg VALUES(151,'51zmt-沈阳新闻','','1','沈阳新闻');
INSERT INTO iptv_epg VALUES(152,'51zmt-求索科学','','1','求索科学');
INSERT INTO iptv_epg VALUES(153,'51zmt-CHC高清电影','','1','CHC高清电影');
INSERT INTO iptv_epg VALUES(154,'51zmt-求索动物','','1','求索动物');
INSERT INTO iptv_epg VALUES(155,'51zmt-求索记录','','1','求索记录');
INSERT INTO iptv_epg VALUES(156,'51zmt-CHC动作电影','','1','CHC动作电影');
INSERT INTO iptv_epg VALUES(157,'51zmt-CHC家庭电影','','1','CHC家庭电影');
INSERT INTO iptv_epg VALUES(158,'51zmt-梨园','','1','梨园');
INSERT INTO iptv_epg VALUES(159,'51zmt-风云音乐','','1','风云音乐');
INSERT INTO iptv_epg VALUES(160,'51zmt-第一剧场','','1','第一剧场');
INSERT INTO iptv_epg VALUES(161,'51zmt-风云剧场','','1','风云剧场');
INSERT INTO iptv_epg VALUES(162,'51zmt-世界地理','','1','世界地理');
INSERT INTO iptv_epg VALUES(163,'51zmt-怀旧剧场','','1','怀旧剧场');
INSERT INTO iptv_epg VALUES(164,'51zmt-兵器科技','','1','兵器科技');
INSERT INTO iptv_epg VALUES(165,'51zmt-女性时尚','','1','女性时尚');
INSERT INTO iptv_epg VALUES(166,'51zmt-CCTV-娱乐','','1','CCTV-娱乐');
INSERT INTO iptv_epg VALUES(167,'51zmt-CCTV-戏曲','','1','CCTV-戏曲');
INSERT INTO iptv_epg VALUES(168,'51zmt-高尔夫网球','','1','高尔夫网球');
INSERT INTO iptv_epg VALUES(169,'51zmt-央视精品','','1','央视精品');
INSERT INTO iptv_epg VALUES(170,'51zmt-彩民在线','','1','彩民在线');
INSERT INTO iptv_epg VALUES(171,'51zmt-汽摩','','1','汽摩');
INSERT INTO iptv_epg VALUES(172,'51zmt-留学世界','','1','留学世界');
INSERT INTO iptv_epg VALUES(173,'51zmt-青年学苑','','1','青年学苑');
INSERT INTO iptv_epg VALUES(174,'51zmt-摄影频道','','1','摄影频道');
INSERT INTO iptv_epg VALUES(175,'51zmt-天元围棋','','1','天元围棋');
INSERT INTO iptv_epg VALUES(176,'51zmt-现代女性','','1','现代女性');
INSERT INTO iptv_epg VALUES(177,'51zmt-早期教育','','1','早期教育');
INSERT INTO iptv_epg VALUES(178,'51zmt-证券资讯','','1','证券资讯');
INSERT INTO iptv_epg VALUES(179,'51zmt-央视台球','','1','央视台球');
INSERT INTO iptv_epg VALUES(180,'51zmt-茶频道','','1','茶频道');
INSERT INTO iptv_epg VALUES(181,'51zmt-武术世界','','1','武术世界');
INSERT INTO iptv_epg VALUES(182,'51zmt-发现之旅','','1','发现之旅');
INSERT INTO iptv_epg VALUES(183,'51zmt-环球奇观','','1','环球奇观');
INSERT INTO iptv_epg VALUES(184,'51zmt-国学','','1','国学');
INSERT INTO iptv_epg VALUES(185,'51zmt-文物宝库','','1','文物宝库');
INSERT INTO iptv_epg VALUES(186,'51zmt-新科动漫','','1','新科动漫');
INSERT INTO iptv_epg VALUES(187,'51zmt-幼儿教育','','1','幼儿教育');
INSERT INTO iptv_epg VALUES(188,'51zmt-老故事','','1','老故事');
INSERT INTO iptv_epg VALUES(189,'51zmt-快乐垂钓','','1','快乐垂钓');
INSERT INTO iptv_epg VALUES(190,'51zmt-书画频道','','1','书画频道');
INSERT INTO iptv_epg VALUES(191,'51zmt-先锋乒羽','','1','先锋乒羽');
INSERT INTO iptv_epg VALUES(192,'51zmt-车迷频道','','1','车迷频道');
INSERT INTO iptv_epg VALUES(193,'51zmt-四海钓鱼','','1','四海钓鱼');
INSERT INTO iptv_epg VALUES(194,'51zmt-环球旅游','','1','环球旅游');
INSERT INTO iptv_epg VALUES(195,'51zmt-京视剧场','','1','京视剧场');
INSERT INTO iptv_epg VALUES(196,'51zmt-弈坛春秋','','1','弈坛春秋');
INSERT INTO iptv_epg VALUES(197,'51zmt-央广健康','','1','央广健康');
INSERT INTO iptv_epg VALUES(198,'51zmt-时代家居','','1','时代家居');
INSERT INTO iptv_epg VALUES(199,'51zmt-时代出行','','1','时代出行');
INSERT INTO iptv_epg VALUES(200,'51zmt-时代风尚','','1','时代风尚');
INSERT INTO iptv_epg VALUES(201,'51zmt-财富天下','','1','财富天下');
INSERT INTO iptv_epg VALUES(202,'51zmt-百姓健康','','1','百姓健康');
INSERT INTO iptv_epg VALUES(203,'51zmt-精品剧场','','1','精品剧场');
INSERT INTO iptv_epg VALUES(204,'51zmt-少儿动漫','','1','少儿动漫');
INSERT INTO iptv_epg VALUES(205,'51zmt-欧美影院','','1','欧美影院');
INSERT INTO iptv_epg VALUES(206,'51zmt-中国教育4台','','1','中国教育4台');

CREATE TABLE iptv_meals (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    content TEXT DEFAULT NULL,
    status INTEGER NOT NULL DEFAULT 1
);
INSERT INTO iptv_meals VALUES(1000,'默认套餐','央视频道(自动聚合)_卫视频道(自动聚合)',1);
INSERT INTO iptv_meals VALUES(1001,'卧室套餐','',1);
CREATE TABLE iptv_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name BIGINT NOT NULL,
    mac TEXT NOT NULL,
    deviceid TEXT NOT NULL,
    model TEXT NOT NULL,
    ip TEXT NOT NULL,
    region TEXT DEFAULT NULL,
    exp BIGINT NOT NULL,
    vpn INTEGER NOT NULL DEFAULT 0,
    idchange INTEGER NOT NULL DEFAULT 0,
    author TEXT DEFAULT NULL,
    authortime BIGINT NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT -1,
    lasttime BIGINT NOT NULL,
    marks TEXT DEFAULT NULL,
    meal INTEGER NOT NULL DEFAULT 1000
);
CREATE TABLE IF NOT EXISTS "iptv_movie"  (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,`name` text,api TEXT DEFAULT NULL,`state` integer);
COMMIT;
