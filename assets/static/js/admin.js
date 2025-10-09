function updateYear() {
	$('#year').text(new Date().getFullYear());
}
updateYear();

$(function() {
	$(".categorylist-btn").first().click();
});

// 替换 main 和 copyright
function replaceContent(html) {
	var $temp = $('<div>').html(html);
	var $newMain = $temp.find('main');
	if ($newMain.length) $('main').replaceWith($newMain);
	// $(".categorylist-btn").first().click();

	// var $newCopyright = $temp.find('p.copyright');
	// if ($newCopyright.length) $('p.copyright').replaceWith($newCopyright);

	updateYear();
}

function getPath(url) {
    try {
        const u = new URL(url, window.location.origin); // 兼容相对路径
        return u.pathname; // 只保留路径
    } catch (e) {
        // 如果是纯相对路径，直接处理
        return url.split(/[?#]/)[0];
    }
}

// 更新菜单 active 状态
function updateMenuActive(url) {
	url = getPath(url);

	// 移除所有已有 active
	$('.nav-drawer a').removeClass('active');
	$('.nav-drawer a').closest('li').removeClass('active');
	$('.nav-item-has-subnav').removeClass('open');
	$('.nav-drawer .nav-item').removeClass('active');

	// 给当前链接及其父菜单添加 active
	var $link = $('.nav-drawer a[href="' + url + '"]');
	$link.parent("li").addClass('active');
	$link.closest(".nav-item-has-subnav").addClass("open");
	$link.closest('.nav-item').addClass('active'); // 父菜单高亮
}

// 加载页面函数
function loadPage(url, pushHistory = true, showName = "") {
	if (!url || url === "javascript:;" || url === "javascript:void(0)") return;
	$.ajax({
		url: url,
		method: 'GET',
		success: function(res) {
			let resStr = typeof res === 'string' ? res : JSON.stringify(res);
			if (resStr.includes('/admin/login')) {
				window.location.href = "/admin/login";
				return;
			}
			replaceContent(res);
			if (url.hash === '#bottom') {
                const images = document.querySelectorAll('img');
                let loadedCount = 0;

                if (images.length === 0) {
                    // 没有图片直接滚动
                    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
                } else {
                    images.forEach(img => {
                        if (img.complete) {
                            loadedCount++;
                        } else {
                            img.addEventListener('load', () => {
                                loadedCount++;
                                if (loadedCount === images.length) {
                                    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
                                }
                            });
                            img.addEventListener('error', () => {
                                loadedCount++;
                                if (loadedCount === images.length) {
                                    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
                                }
                            });
                        }
                    });

                    // 如果所有图片已经加载完成
                    if (loadedCount === images.length) {
                        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
                    }
                }
            }

			if (pushHistory) history.pushState(null, '', url);
			updateMenuActive(url);
			if (showName !== "") {
				showlist(showName);
			}else{
				$(".categorylist-btn").first().click();
			}
		},
		error: function(err) {
			console.error('请求失败:', err);
			return false; // 返回 false 表示失败
		}
	});
}

function submitFormPOST(btn, showChannel = false) {

	// 获取按钮所在的 form
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
		return;
	}

	var action = form.action || window.location.pathname; // 默认当前路径

	let showName = "";

	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		if (key === "iconfile" || key === "bjfile") {
			return;
		}
		params.append(key, value);
	});
	params.append(btn.name, "");
	if (params.has("submitappinfo") && !params.has("up_sets")) {
		params.append("up_sets", 0);
	}

	if(showChannel && (btn.name === "submit_addtype" || 
		btn.name  === "submit_modifytype" || 
		btn.name  === "submit_moveup" ||
		btn.name  === "submit_movedown" ||
		btn.name  === "submit_movetop"
	)){
		showName = params.get("category");
	}

	if(showChannel && btn.name === "submitsave"){
		showName = params.get("categoryname");
	}
	// 使用 fetch AJAX 提交
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); // 先拿到响应内容
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); // 或 response.json()
	})
	// .then(response => response.json())
	.then(data => {
		lightyear.notify(data.msg, data.type, 3000);
		if (data.type === "success") {
			loadPage(action, true, showName);
			if ($('.modal-backdrop').length > 0) {
				document.querySelector('.modal-backdrop').remove();
				document.body.classList.remove('modal-open');
				document.body.style.overflow = ''; // 恢复滚动条
			}
		}
	})
	.catch(err => {
		// console.error("提交失败:", err);
		lightyear.notify("提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
	});
}

function submitFormGET(btn) {
	// 获取按钮所在的 form
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
	}

	// 获取 form 的 action，默认当前路径
	var action = form.action || window.location.pathname;

	// 把表单数据拼接成 query string
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		params.append(key, value);
	});
	// 把当前按钮的 name 加进去（通常用于区分提交动作）
	if (btn.name) {
		params.append(btn.name, "");
	}

	// 拼接到 URL
	const url = action + (action.includes("?") ? "&" : "?") + params.toString();

	loadPage(url);
}



// 菜单点击事件
$('.nav-drawer .nav-item > a, .nav-drawer .nav-subnav a').click(function(e) {
	e.preventDefault();
	var url = $(this).attr('href');
	loadPage(url);
});

// 页面首次加载时高亮当前菜单
var path = window.location.pathname;
$('.nav-drawer a').each(function() {
	var href = $(this).attr('href');
	if (href !== "javascript:;" && path.startsWith(href)) {
		updateMenuActive(href);
	}
});

// 处理浏览器前进/后退
window.addEventListener('popstate', function() {
	loadPage(location.pathname, false);
});

// 注销
$('#logout').click(function() {
	$.get('/admin/logout', function() {
		location.reload();
	});
});

document.addEventListener("click", function(e) {
    // 找到点击的 <a>，并且包含 loaduser 类
	let  link = e.target.closest("a.aboutbottom");
    if (link) {
		e.preventDefault();
		const url = new URL(link.href, window.location.origin);
		
		loadPage(url);
		
		return; // 不是目标元素，忽略
	}
    link = e.target.closest("a.loaduser");
    if (!link) return; // 不是目标元素，忽略

    // 确保 <a> 在 main 容器内
    const mainContainer = document.querySelector("main.lyear-layout-content");
    if (!mainContainer || !mainContainer.contains(link)) return;

    e.preventDefault(); // 阻止默认跳转

    const url = new URL(link.href, window.location.origin);

    // 获取当前搜索关键词（可选）
    const keywords = document.querySelector("input[name='keywords']")?.value || "";
    if (keywords) url.searchParams.set("keywords", keywords);

    loadPage(url); // 调用你的 AJAX 加载函数
});



function submitFormCounts(){
	const form = document.getElementById("recCounts");
	const url = new URL(window.location.href);
	const keywords = document.querySelector("input[name='keywords']")?.value || "";
        if (keywords) url.searchParams.set("keywords", keywords);

	const recCountsSelect = document.querySelector("select[name='recCounts']");
	if (recCountsSelect && recCountsSelect.value) {
		url.searchParams.set("recCounts", recCountsSelect.value);
	}

	loadPage(url);
}

function checkboxall(a){
	var ck=document.getElementsByName("ids[]");
	for (var i = 0; i < ck.length; i++) {
		if(a.checked){
			ck[i].checked=true;
		}else{
			ck[i].checked=false;
		}
	}
}

function checkboxAllName(a){
	var ck=document.getElementsByName("names[]");
	for (var i = 0; i < ck.length; i++) {
		if(a.checked){
			ck[i].checked=true;
		}else{
			ck[i].checked=false;
		}
	}
}

function showlist(name){
	$("#srclist").val("正在加载中...");
	$.ajax({
		url: "/admin/channels",
		type: "POST",
		data: { category: name, getchannels: "" },
		success: function(data) {
			$("#srclist").val(data);
		}
	});

	$("#typename").val(name);
	$("#typename0").val(name);
	$("#categoryname").val(name);
}

function categorycheck(name){
	$.ajax({
		url: "/admin/channels",
		type: "POST",
		data: { category: name, forbiddenchannels: "" },
		success: function(data) {
			lightyear.notify(data.msg, data.type, 3000); // success, warning, danger, info
		}
	});
}

function tdBtnPOST(btn) {

	var action = window.location.pathname; // 默认当前路径

	var params = new URLSearchParams();

	
	if ($(btn).is(":checkbox")) {
		params.append(btn.name, btn.checked ? 1 : 0);
	}else{
		params.append(btn.name, btn.value);
	}

	// 使用 fetch AJAX 提交
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); // 先拿到响应内容
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); // 或 response.json()
	})
	.then(data => {
		lightyear.notify(data.msg, data.type, 3000);
		if (data.type === "success") {
			if ($('.modal-backdrop').length > 0) {
				document.querySelector('.modal-backdrop').remove();
				document.body.classList.remove('modal-open');
				document.body.style.overflow = ''; // 恢复滚动条
			}
			loadPage(action);
		}
	})
	.catch(err => {
		// console.error("提交失败:", err);
		lightyear.notify("提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
	});
}


function mealsGetCategory(btn) {

	var action = window.location.pathname;

	const params = new URLSearchParams();

	if (btn.name && btn.value) {
		params.append(btn.name, btn.value);
	}else if (btn.name) {
		params.append(btn.name, "");
	}else{
		lightyear.notify("表单提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
		return ;
	}

	if (btn.name === "editmeal") {
		var $tr = $(btn).closest("tr"); // 获取当前行的 jQuery 对象
		var mealid = $tr.find(".meal-id").data("value");
		var mealname = $tr.find(".meal-name").data("value");
		
		$("#mealId").val(mealid);
		$("#mealName").val(mealname);
	}

	// 使用 fetch AJAX 提交
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); // 先拿到响应内容
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); // 或 response.json()
	})
	.then(res => {
		
		if (res.type === "success") {
			// if ($('.modal-backdrop').length > 0) {
			// 	document.querySelector('.modal-backdrop').remove();
			// }
			if (res.data && res.data.length > 0) {
				var html = "";
				res.data.forEach(function(item) {
					html += '<label class="lyear-checkbox checkbox-inline">';
					html += '<input type="checkbox" name="names[]" value="' + item.name + '"' 
							+ (item.checked ? ' checked="checked"' : '') + '>';
					html += '<span>' + item.name + '</span>';
					html += '</label>';
				});
				$(".form-inline.meal-checkbox").html(html);
			}else{
				$(".form-inline.meal-checkbox").html("<span>无数据</span>");
			}
			// loadPage(action);
		}else{
			lightyear.notify(res.msg, res.type, 3000);
		}
	})
}

function epgsGetChannel(btn) {

	var action = window.location.pathname;

	const params = new URLSearchParams();

	if (btn.name && btn.value) {
		params.append(btn.name, btn.value);
	}else if (btn.name) {
		params.append(btn.name, "");
	}else{
		lightyear.notify("表单提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
		return ;
	}

	if (btn.name === "editepg") {
		var $tr = $(btn).closest("tr"); // 获取当前行的 jQuery 对象
		var epgid = $tr.find(".epg-id").data("value");
		var epgname = $tr.find(".epg-name").data("value");
		var epgremarks = $tr.find(".epg-remarks").data("value");

		var prefix = epgname.split("-")[0];
		var name = epgname.split("-")[1];
		
		$("#editepgselect").val(prefix);
		$("#epgId").val(epgid);
		$("#epgName").val(name);
		$("#epgRemarks").val(epgremarks);
	}

	// 使用 fetch AJAX 提交
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); // 先拿到响应内容
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); // 或 response.json()
	})
	.then(res => {
		// lightyear.notify(res.msg, res.type, 3000);
		if (res.type === "success") {
			// if ($('.modal-backdrop').length > 0) {
			// 	document.querySelector('.modal-backdrop').remove();
			// }
			if (res.data && res.data.length > 0) {
				var html = "";
				res.data.forEach(function(item) {
					html += '<label class="lyear-checkbox checkbox-inline">';
					html += '<div style=" float: left;background: ' + (item.checked ? '#7fff00;' : '#E7E7E7;') +' margin-right: 3px; margin-bottom: 3px; padding: 2px 5px;">';
					html += '<input type="checkbox" name="names[]" value="' + item.name + '"' 
							+ (item.checked ? ' checked="checked"' : '') + '>';
					html += '<span>' + item.name + '</span></div>';
					html += '</label>';
				});
				$(".form-inline.epg-checkbox").html(html);
			}else{
				$(".form-inline.epg-checkbox").html("<span>无数据</span>");
			}
			// loadPage(action);
		}else{
			lightyear.notify(res.msg, res.type, 3000);
		}
	})
}

function uploadIcon(event) {
	const file = event.target.files[0];
	if (!file) return;

	const formData = new FormData($('#appform')[0]);
	$.ajax({
		url: '/admin/client/uploadIcon',  // ✅ 上传接口
		type: 'POST',
		data: formData,
		contentType: false,
		processData: false,
		success: function(res) {
			$('#iconInput').val(''); // ✅ 清空文件输入框
			lightyear.notify(res.msg, res.type, 3000); // ✅ 上传成功后，显示提示信息
			if (res.code === 1) {
				$('#iconImg').attr('src', res.data.url);
				$('#iconContainer').show();
			}
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000);
			$('#iconContainer').hide();
			$('#iconInput').val('');
		}
	});
};

function deleteIcon() {
	src = $("#iconImg").attr("src"); // 获取图片的src属性值

	if (src.trim() === "") {
		$('#iconContainer').hide();
		$('#iconInput').val('');
		return;
	}

	const params = new URLSearchParams();
	params.append("deleteIcon", "");

	$.ajax({
		url: '/admin/client',   // ✅ 删除接口
		type: 'POST',
		data: params.toString(),
		contentType: 'application/x-www-form-urlencoded',
		success: function(res) {
			lightyear.notify(res.msg, res.type, 3000); 
			$('#iconContainer').hide();
			$('#iconImg').attr('src', '');
			$('#iconInput').val('');
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000); // ✅ 上传失败后，显示提示信息
		}
	});
};

function uploadBj(event) {
	const file = event.target.files[0];
	if (!file) return;

	const formData = new FormData($('#bj_form')[0]);
	$.ajax({
		url: '/admin/client/uploadBj',  // ✅ 上传接口
		type: 'POST',
		data: formData,
		contentType: false,
		processData: false,
		success: function(res) {
			$('#bjInput').val(''); // ✅ 清空文件输入框
			lightyear.notify(res.msg, res.type, 3000); // ✅ 上传成功后，显示提示信息
			if (res.code === 1) {
				imgName = res.data.name;
					$('.reload-bj').append(`
						<div class="form-group" id="bj_`+imgName+`" style="position:relative; margin-right: 5px;">
							<img id="`+imgName+`" src="/images/`+imgName+`.png" alt="预览" style="height:38px; border:1px solid #ccc; border-radius:4px; cursor:pointer;">
							<!-- 删除按钮 -->
							<span class="delete-btn" onclick="deleteBj('`+imgName+`')" data-name="`+imgName+`" id="`+imgName+`" style="
								position:absolute; top:-8px; right:-8px;
								background:#f00; color:#fff; border-radius:50%;
								font-size:14px; line-height:16px; width:16px; height:16px;
								text-align:center; cursor:pointer;")">×</span>
						</div>
						`);
			}
		},
		error: function(res) {
			$('#bjInput').val('');
			lightyear.notify(res.msg, res.type, 3000);
		}
	});
};

function deleteBj(name) {
	const params = new URLSearchParams();
	params.append("deleteBj", name);

	$.ajax({
		url: '/admin/client',   // ✅ 删除接口
		type: 'POST',
		data: params.toString(),
		contentType: 'application/x-www-form-urlencoded',
		success: function(res) {
			lightyear.notify(res.msg, res.type, 3000);
			if (res.code === 1) {
				$('#bj_'+name).remove();
				$('#bjInput').val('');
			}
			
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000); // ✅ 上传失败后，显示提示信息
			$('#bjInput').val(''); // ✅ 清空输入框
		}
	});
}

function moviesGet(btn) {

	if (btn.name === "editmovie") {
		var $tr = $(btn).closest("tr"); // 获取当前行的 jQuery 对象
		var movieid = $tr.find(".movie-id").data("value");
		var moviename = $tr.find(".movie-name").data("value");
		var movieapi = $tr.find(".movie-api").data("value");
		
		$("#movieId").val(movieid);
		$("#movieName").val(moviename);
		$("#mealApi").val(movieapi);
	}
}

function confirmAndSubmit(btn ,msg) {
    $.confirm({
        title: '操作确认',
        content: msg,
        type: 'green',
        buttons: {
            confirm: {
                text: '确认',
                btnClass: 'btn-success',
                action: function () {
					submitFormPOST(btn);
                }
            },
            cancel: {
                text: '取消',
                btnClass: 'btn-danger'
            }
        }
    });
}

function getCategory(btn) {
	var action = window.location.pathname;

	const params = new URLSearchParams();

	if (btn.name && btn.value) {
		params.append(btn.name, btn.value);
	}else if (btn.name) {
		params.append(btn.name, "");
	}else{
		lightyear.notify("表单提交失败", "danger", 3000);
		// document.querySelector('.modal-backdrop').remove();
		// document.body.classList.remove('modal-open');
		// document.body.style.overflow = ''; // 恢复滚动条
		return ;
	}

	if (btn.name === "editCategory") {
		var $tr = $(btn).closest("tr"); // 获取当前行的 jQuery 对象
		var cid = $tr.find(".c-id").data("value");
		var cname = $tr.find(".c-name").data("value");
		var curl = $tr.find(".c-url").data("value");
		var ca = $tr.find(".c-a").data("value");

		$("#cId").val(cid);
		$("#listname").val(cname);
		$("#listurl").val(curl);
		$("#autocategory").prop("checked", ca === 1);
	}
}