function updateYear() {
	$('#year').text(new Date().getFullYear());
}
updateYear();
function replaceContent(html) {
	var $temp = $('<div>').html(html);
	var $newMain = $temp.find('main');
	if ($newMain.length) $('main').replaceWith($newMain);
	updateYear();
}
function getPath(url) {
    try {
        const u = new URL(url, window.location.origin); 
        return u.pathname; 
    } catch (e) {
        return url.split(/[?#]/)[0];
    }
}
function updateMenuActive(url) {
	url = getPath(url);
	$('.nav-drawer a').removeClass('active');
	$('.nav-drawer a').closest('li').removeClass('active');
	$('.nav-item-has-subnav').removeClass('open');
	$('.nav-drawer .nav-item').removeClass('active');
	var $link = $('.nav-drawer a[href="' + url + '"]');
	$link.parent("li").addClass('active');
	$link.closest(".nav-item-has-subnav").addClass("open");
	$link.closest('.nav-item').addClass('active'); 
	$('.selectpicker').selectpicker('refresh');
}
function loadPage(url, pushHistory = true) {
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
                    if (loadedCount === images.length) {
                        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
                    }
                }
            }
			if (pushHistory) history.pushState(null, '', url);
			updateMenuActive(url);
		},
		error: function(err) {
			console.error('请求失败:', err);
			return false; 
		}
	});
}
function submitFormPOST(btn) {
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		return;
	}
	lightyear.loading('show');
	var action = form.action || window.location.pathname; 
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		if (key === "iconfile" || key === "bjfile" || key === "paylistfile") {
			return;
		}
		params.append(key, value);
	});
	params.append(btn.name, "");
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(data => {
		lightyear.loading('hide');
		lightyear.notify(data.msg, data.type, 3000);
		if (data.type === "success") {
			var sub = $(btn).closest('.modal');
			sub.modal('hide');
			loadPage(window.location.href);
		}
	})
	.catch(err => {
		lightyear.notify("提交失败", "danger", 3000);
	});
}
function submitFormGET(btn) {
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		return;
	}
	var action = form.action || window.location.pathname;
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		params.append(key, value);
	});
	if (btn.name) {
		params.append(btn.name, "");
	}
	const url = action + (action.includes("?") ? "&" : "?") + params.toString();
	loadPage(url);
}
$('.nav-drawer .nav-item > a, .nav-drawer .nav-subnav a').click(function(e) {
	e.preventDefault();
	var url = $(this).attr('href');
	loadPage(url);
});
var path = window.location.pathname;
$('.nav-drawer a').each(function() {
	var href = $(this).attr('href');
	if (href !== "javascript:;" && path.startsWith(href)) {
		updateMenuActive(href);
	}
});
window.addEventListener('popstate', function() {
	loadPage(location.pathname, false);
});
$('#logout').click(function() {
	$.get('/admin/logout', function() {
		location.reload();
	});
});
document.addEventListener("click", function(e) {
	let  link = e.target.closest("a.aboutbottom");
    if (link) {
		e.preventDefault();
		const url = new URL(link.href, window.location.origin);
		loadPage(url);
		return; 
	}
    link = e.target.closest("a.loaduser");
    if (!link) return; 
    const mainContainer = document.querySelector("main.lyear-layout-content");
    if (!mainContainer || !mainContainer.contains(link)) return;
    e.preventDefault(); 
    const url = new URL(link.href, window.location.origin);
    const keywords = document.querySelector("input[name='keywords']")?.value || "";
    if (keywords) url.searchParams.set("keywords", keywords);
    loadPage(url); 
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
function clearCheck(){
	var ck=document.getElementsByName("names[]");
	for (var i = 0; i < ck.length; i++) {
		ck[i].checked=false;
	}
}
function tdBtnPOST(btn) {
	var action =  window.location.href; 
	var params = new URLSearchParams();
	if ($(btn).is(":checkbox")) {
		params.append(btn.name, btn.checked ? 1 : 0);
	}else{
		params.append(btn.name, btn.value);
	}
	lightyear.loading('show');
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(data => {
		lightyear.loading('hide');
		if (data.type === "success") {
			lightyear.notify(data.msg, data.type, 1000);
			if (btn.name.includes("del") && !btn.name.includes("dellist")) {
				const tr = btn.closest("tr");
				if (tr) tr.remove();
			}else if (btn.name.includes("Status")) {
				const tr = btn.closest("tr");
				console.log(tr);
				const statusTd = tr.querySelector("td.status-show");
				console.log(statusTd);
    			const font = statusTd.querySelector("font");
				console.log(font);
				if (font && font.textContent.includes("上线")) {
					font.textContent = "下线";
					font.color = "red";
					btn.classList.remove("btn-warning");
					btn.classList.add("btn-success");
					btn.textContent = "上线"; 
				} else {
					font.textContent = "上线";
					font.color = "#33a996";
					btn.classList.remove("btn-success");
					btn.classList.add("btn-warning");
					btn.textContent = "下线"; 
				}
			}else{
				var sub = $(btn).closest('.modal');
				sub.modal('hide');
				loadPage(action);
			}
		}else {
			lightyear.notify(data.msg, data.type, 3000);
		}
	})
	.catch(err => {
		lightyear.notify("提交失败"+err, "danger", 3000);
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
		return ;
	}
	if (btn.name === "editmeal") {
		var $tr = $(btn).closest("tr"); 
		var mealid = $tr.find(".meal-id").data("value");
		var mealname = $tr.find(".meal-name").data("value");
		$("#mealId").val(mealid);
		$("#mealName").val(mealname);
	}
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(res => {
		if (res.type === "success") {
			if (res.data && res.data.length > 0) {
				var html = "";
				res.data.forEach(function(item) {
					html += '<label class="lyear-checkbox checkbox-inline">';
					html += '<input type="checkbox" name="ids[]" value="' + item.id + '"' 
							+ (item.checked ? ' checked="checked"' : '') + '>';
					html += '<span>' + item.name + '</span>';
					html += '</label>';
				});
				$(".form-inline.meal-checkbox").html(html);
			}else{
				$(".form-inline.meal-checkbox").html("<span>无数据</span>");
			}
		}else{
			lightyear.notify(res.msg, res.type, 3000);
		}
	})
}
function epgEdit(btn) {
	var $tr = $(btn).closest("tr"); 
	var epgid = $tr.find(".epg-id").data("value");
	var epgname = $tr.find(".epg-name").data("value");
	var epgremarks = $tr.find(".epg-remarks").data("value");
	var index = epgname.indexOf("-");
	var prefix, name;
	if (index !== -1) {
		prefix = epgname.substring(0, index);
		name = epgname.substring(index + 1);
	} else {
		prefix = epgname;
		name = "";
	}
	$("#editepgselect").val(prefix);
	$("#epgId").val(epgid);
	$("#epgName").val(name);
	$("#epgRemarks").val(epgremarks);
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
		return ;
	}
	var $tr = $(btn).closest("tr"); 
	var epgid = $tr.find(".epg-id").data("value");
	var epgname = $tr.find(".epg-name").data("value");
	var index = epgname.indexOf("-");
	var prefix, name;
	if (index !== -1) {
		prefix = epgname.substring(0, index);
		name = epgname.substring(index + 1);
	} else {
		prefix = epgname;
		name = "";
	}
	$("#bdingepgselect").val(prefix);
	$("#epgId1").val(epgid);
	$("#epgName1").val(name);
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(res => {
		lightyear.notify(res.msg, res.type, 1000);
		if (res.type === "success") {
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
		}else{
			lightyear.notify(res.msg, res.type, 3000);
		}
	})
}
function uploadIcon(event) {
	const file = event.target.files[0];
	if (!file) return;
	const formData = new FormData();
	const fileInput = document.querySelector('input[name="iconfile"]'); 
	if (fileInput && fileInput.files.length > 0) {
	formData.append("iconfile", fileInput.files[0]);
	}else{
		lightyear.notify("❌ 请选择文件！", "danger", 1000);
		return;
	}
	$.ajax({
		url: '/admin/client/uploadIcon',  
		type: 'POST',
		data: formData,
		contentType: false,
		processData: false,
		success: function(res) {
			$('#iconInput').val(''); 
			lightyear.notify(res.msg, res.type, 1000); 
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
	src = $("#iconImg").attr("src"); 
	if (src.trim() === "") {
		$('#iconContainer').hide();
		$('#iconInput').val('');
		return;
	}
	const params = new URLSearchParams();
	params.append("deleteIcon", "");
	$.ajax({
		url: '/admin/client',   
		type: 'POST',
		data: params.toString(),
		contentType: 'application/x-www-form-urlencoded',
		success: function(res) {
			lightyear.notify(res.msg, res.type, 1000); 
			$('#iconContainer').hide();
			$('#iconImg').attr('src', '');
			$('#iconInput').val('');
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000); 
		}
	});
};
function uploadBj(event) {
	const file = event.target.files[0];
	if (!file) return;
	const formData = new FormData();
	const fileInput = document.querySelector('input[name="bjfile"]'); 
	if (fileInput && fileInput.files.length > 0) {
		formData.append("bjfile", fileInput.files[0]);
	} else {
		lightyear.notify("❌ 请选择文件！", "danger", 1000);
		return;
	}
	$.ajax({
		url: '/admin/client/uploadBj',  
		type: 'POST',
		data: formData,
		contentType: false,
		processData: false,
		success: function(res) {
			$('#bjInput').val(''); 
			lightyear.notify(res.msg, res.type, 1000); 
			if (res.code === 1) {
				const imgName = res.data.name;
				const html = `<div class="form-group" id="bj_${imgName}" style="position:relative; margin-right:5px;"img id="${imgName}" src="/images/${imgName}.png" alt="预览" style="height:38px; border:1px solid #ccc; border-radius:4px; cursor:pointer;"span class="delete-btn" onclick="deleteBj('${imgName}')" data-name="${imgName}" id="${imgName}_del" style="position:absolute; top:-8px; right:-8px; background:#f00; color:#fff; border-radius:50%; font-size:14px; line-height:16px; width:16px;  height:16px; text-align:center; cursor:pointer;">×</span/div>`;
				$('.reload-bj').append(html);
			}
		},
		error: function(err) {
			$('#bjInput').val('');
			lightyear.notify("❌ 上传失败！", "danger", 3000);
			console.error("Upload error:", err);
		}
	});
};
function deleteBj(name) {
	const params = new URLSearchParams();
	params.append("deleteBj", name);
	$.ajax({
		url: '/admin/client',   
		type: 'POST',
		data: params.toString(),
		contentType: 'application/x-www-form-urlencoded',
		success: function(res) {
			lightyear.notify(res.msg, res.type, 1000);
			if (res.code === 1) {
				$('#bj_'+name).remove();
				$('#bjInput').val('');
			}
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000); 
			$('#bjInput').val(''); 
		}
	});
}
function moviesGet(btn) {
	if (btn.name === "editmovie") {
		var $tr = $(btn).closest("tr"); 
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
function getCategoryList(btn) {
	var $tr = $(btn).closest("tr"); 
	var cid = $tr.find(".cl-id").data("value");
	var cname = $tr.find(".cl-name").data("value");
	var curl = $tr.find(".cl-url").data("value");
	var cua = $tr.find(".cl-ua").data("value");
	var ca = $tr.find(".cl-a").data("value");
	var cr = $tr.find(".cl-r").data("value");
	$("#clId").val(cid);
	$("#listname").val(cname);
	$("#listurl").val(curl);
	$("#listua").val(cua);
	$("#autocategory").prop("checked", ca === 1);
	$("#repeat").prop("checked", cr === 1);
}
function getCategory(btn) {
	var $tr = $(btn).closest("tr"); 
	var cid = $tr.find(".ca-id").data("value");
	var cname = $tr.find(".ca-name").data("value");
	var ctype = $tr.find(".ca-type").data("value");
	var cua = $tr.find(".ca-ua").data("value");
	var crules = $tr.find(".ca-rules").data("value");
	var cproxy = $tr.find(".ca-proxy").data("value");
	$("#caId").val(cid);
	$("#caname").val(cname);
	$("#caua").val(cua);
	$("#rules").val(crules);
	if (ctype === "auto") {
		document.getElementById('rules').disabled = false;
	}else{
		document.getElementById('rules').disabled = true;
	}
	$("#autoagg").prop("checked", ctype === "auto");
	$("#proxy").prop("checked", cproxy === 1);
}
function getChannels(id) {
    $("#showcaId").val(id);
    $.ajax({
        url: "/admin/channels",
        type: "POST",
        data: { caId: id, getchannels: "" },
        success: function(data) {
            lightyear.notify(data.msg, data.type, 3000);
            if (data.type === "success") {
                const tbody = document.getElementById("channellist_tbody");
                tbody.innerHTML = "";
                data.data.forEach(item => {
                    const tr = document.createElement("tr");
                    tr.align = "center";
                    const displayUrl = item.url.length > 20 
                        ? item.url.substring(0, 10) + "..." + item.url.slice(-10) 
                        : item.url;
                    tr.innerHTML =
  '<td style="display:none;" class="ch-id" data-value="' + item.id + '">' + item.id + '</td>' +
  '<td style="display:none;" class="ch-eid" data-value="' + item.e_id + '">' + item.e_id + '</td>' +
  '<td class="ch-name" data-value="' + item.name + '">' + item.name + '</td>' +
  '<td class="ch-url" data-value="' + item.url + '"><a href="' + item.url + '" target="_blank">' + displayUrl + '</a></td>' +
  '<td class="status-show">' + (item.status === 1 ? '<font color="#33a996">上线</font>' : '<font color="red">下线</font>') + '</td>' +
  '<td>' + (item.epg_name || '未绑定') + '</td>' +
  '<td>' + (item.logo === '' 
      ? '无' 
      : '<div id="logo_' + item.id + '" style="position:relative;"><img class="ch-logo" src="' + item.logo + '" alt="预览" style="background-color:black;height:38px;border:1px solid #ccc;border-radius:4px;cursor:pointer;"></div>') + '</td>' +
  '<td>' +
    '<button type="button" onclick="tdBtnPOST(this)" name="channelsStatus" value="' + item.id + '" class="btn btn-xs ' + (item.status === 1 ? 'btn-warning">下线' : 'btn-success">上线') + '</button>&nbsp;' +
    '<button class="btn btn-xs btn-info" type="button" value="' + item.id + '" data-toggle="modal" onclick="editChannel(this)" data-target="#editchannel">编辑</button>&nbsp;' +
    '<button class="btn btn-xs btn-danger" type="button" onclick="tdBtnPOST(this)" name="dellist" value="' + item.id + '">删除</button>' +
  '</td>';
                    tbody.appendChild(tr);
                    $('.selectpicker').selectpicker('refresh'); 
                });
            }
        },
        error: function() {
            lightyear.notify("请求失败", "danger", 3000);
        }
    });
}
function getChannels_pindao(id) {
    $("#showcaId").val(id);
    $.ajax({
        url: "/admin/fenlei",
        type: "POST",
        data: { caId: id, getchannels_pindao: "" },
        success: function(data) {
            lightyear.notify(data.msg, data.type, 3000);
            if (data.type === "success") {
                const tbody = document.getElementById("channellist_tbody_pindao");
                tbody.innerHTML = "";
                console.log('普通日志', data);
                data.data.forEach(item => {
                    const tr = document.createElement("tr");
                    tr.align = "center";
                    const displayUrl = item.url.length > 20
                        ? item.url.substring(0, 10) + "..." + item.url.slice(-10)
                        : item.url;
                    tr.innerHTML =
                        '<td style="display:none;" class="ch-id" data-value="' + item.id + '">' + item.id + '</td>' +
                        // '<td style="display:none;" class="rss-list" data-value="' + item.rss_list + '">' + item.rss_list + '</td>' +
                        '<td class="ch-name" data-value="' + item.name + '">' + item.name + '</td>' +
                        '<td class="status-show">' + (item.status === 1 ? '<font color="#33a996">上线</font>' : '<font color="red">下线</font>') + '</td>' +
                        '<td>' + (item.logo === ''
                            ? '无'
                            : '<div id="logo_' + item.id + '" style="position:relative;"><img class="ch-logo" src="' + item.logo + '" alt="预览" style="background-color:black;height:38px;border:1px solid #ccc;border-radius:4px;cursor:pointer;"></div>') + '</td>' +
                        '<td>' +
                        '<button type="button" onclick="tdBtnPOST(this)" name="channelsStatus" value="' + item.id + '" class="btn btn-xs ' + (item.status === 1 ? 'btn-warning">下线' : 'btn-success">上线') + '</button>&nbsp;' +
                        '<button class="btn btn-xs btn-info" type="button" value="' + item.id + '" data-toggle="modal" onclick="editChannel_pindao(this)" data-target="#editChannel_pindao">修改</button>&nbsp;' +
                        '<button class="btn btn-xs btn-danger" type="button" onclick="tdBtnPOST(this)" name="dellist" value="' + item.id + '">删除</button>' +
                        '</td>';
                    // 将 epg_rules 存储在 tr 的 data 属性中
                    $(tr).data('rss-rules', item.rss_rules);
                    $(tr).data('epg-rules', item.epg_rules);
                    $(tr).data('rss-list', item.rss_list);
                    $(tr).data('epg-list', item.epg_list);
                    tbody.appendChild(tr);
                    $('.selectpicker').selectpicker('refresh');
                });
            }
        },
        error: function() {
            lightyear.notify("请求失败", "danger", 3000);
        }
    });
}

function getChannelsTxt(btn){
	var $tr = $(btn).closest("tr"); 
	var cid = $tr.find(".ca-id").data("value");
	var cname = $tr.find(".ca-name").data("value");
	$("#showtxtcaId").val(cid);
	$("#showtxtCaname").val(cname);
	$.ajax({
		url: "/admin/channels",
		type: "POST",
		data: { caId: cid, getchannels: "" },
		success: function(data) {
			lightyear.notify(data.msg, data.type, 3000);
			if (data.type === "success") {
				const tbody = document.getElementById("channellist_tbody");
				var result = "";
				data.data.forEach(item => {
					result += item.status + "|" + item.name + "," + item.url + "\n";
				});
				$("#srclist").val(result);
			}
		},
		error: function() {
			lightyear.notify("请求失败", "danger", 3000);
		}
	});
}
function editChannel(btn) {
	var $tr = $(btn).closest("tr"); 
	var chid = $tr.find(".ch-id").data("value");
	var cname = $tr.find(".ch-name").data("value");
	var curl = $tr.find(".ch-url").data("value");
	var ceid = $tr.find(".ch-eid").data("value");
	var logoSrc = $(btn).closest("tr").find("td div img.ch-logo").attr("src") || "";
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		return;
	}
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		params.append(key, value);
	});
	$("#editcaId").val(params.get("caId"));
	$("#chId").val(chid);
	$("#chname").val(cname);
	$("#chURL").val(curl);
	$("#e_id").val(ceid);
	$("#e_id").selectpicker('refresh');
	if (!logoSrc) {
    	$("#logoContainerEdit").hide();
	} else {
		$("#logoContainerEdit").show().find("img").attr("src", logoSrc);
	}
}
function editChannel_pindao(btn) {
    var $tr = $(btn).closest("tr");
    var chid = $tr.find(".ch-id").data("value");
    var cname = $tr.find(".ch-name").data("value");
    var rss_rules = $tr.data("rss-rules");
    var epg_rules = $tr.data("epg-rules");
    var rss_list = $tr.data("rss-list");
    var epg_list = $tr.data("epg-list");
    // console.log("rss_rules",rss_rules)
    // console.log("epg_rules",epg_rules)
    // console.log("rss_list",rss_list)
    var curl = $tr.find(".ch-url").data("value");
    // var ceid = $tr.find(".ch-eid").data("value");
    var logoSrc = $(btn).closest("tr").find("td div img.ch-logo").attr("src") || "";
    var form = btn.closest("form");
    if (!form) {
        lightyear.notify("表单提交失败", "danger", 3000);
        return;
    }
    const params = new URLSearchParams();
    new FormData(form).forEach((value, key) => {
        console.log(key,value)
        params.append(key, value);
    });
    console.log("editChannel_pindao")
    $("#editcaId").val(params.get("caId"));
    $("#chId").val(chid);
    $("#chname").val(cname);
    $("#chURL").val(curl);
    // $("#e_id").val(ceid);
    // $("#e_id").selectpicker('refresh');
    $("#rss_rules").val(rss_rules);
    $("#epg_rules").val(epg_rules);

    // 处理 rss_list（多选框）
    if (rss_list) {
        // 将字符串转换为数组（如 "-1,0,1,2" -> ["-1", "0", "1", "2"]）
        var rssArray = rss_list.split(",");
        // 设置多选框选中状态
        $("#rss_list_id").val(rssArray);
    }

    if (epg_list) {
        // 将字符串转换为数组（如 "-1,0,1,2" -> ["-1", "0", "1", "2"]）
        var epgArray = epg_list.split(",");
        // 设置多选框选中状态
        $("#epg_list_id").val(epgArray);
    }

    // 刷新选择器
    $("#rss_list_id").selectpicker("refresh");
    $("#epg_list_id").selectpicker("refresh");

    // 添加事件监听器，当选择"全部"或"无"时取消其他选项
    $('#rss_list_id, #epg_list_id').on('changed.bs.select', function(e, clickedIndex, isSelected, previousValue) {
        var select = $(this);
        var selectedValues = select.val() || [];

        // 获取当前点击的值
        var clickedValue = $(select.find('option')[clickedIndex]).val();

        // 如果选择了"全部"(-1)或"无"(0)
        if (clickedValue === "-1" || clickedValue === "0") {
            // 取消所有其他选项，只保留当前点击的"全部"或"无"
            var newValues = selectedValues.filter(val => val === clickedValue);
            select.val(newValues);
        }
        // 如果选择了其他选项(非"全部"和"无")
        else if (clickedValue !== "-1" && clickedValue !== "0") {
            // 取消"全部"和"无"选项
            var newValues = selectedValues.filter(val => val !== "-1" && val !== "0");
            select.val(newValues);
        }

        select.selectpicker('refresh');
    });

    if (!logoSrc) {
        $("#logoContainerEdit").hide();
    } else {
        $("#logoContainerEdit").show().find("img").attr("src", logoSrc);
    }
}

function saveChannelsOne(btn){
	var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		return;
	}
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		if (key === "logofile") {
			return;
		}
		params.append(key, value);
	});
	params.append("saveChannelsOne", "");
	var editchmodal = $(btn).closest('.modal');
	$.ajax({
		url: "/admin/channels",
		type: "POST",
		data: params.toString(),
		success: function(data) {
			lightyear.notify(data.msg, data.type, 3000);
			if (data.type === "success") {
				editchmodal.modal('hide');
				getChannels(params.get("caId"));
			}
		},
		error: function() {
			lightyear.notify("请求失败", "danger", 3000);
		}
	});
}

function saveChannelsOne_fenlei(btn){
    var form = btn.closest("form");
    if (!form) {
        lightyear.notify("表单提交失败", "danger", 3000);
        return;
    }
    const params = new URLSearchParams();
    new FormData(form).forEach((value, key) => {
        if (key === "logofile") {
            return;
        }
        params.append(key, value);
    });

    // 处理rss_list_id - 如果未选择任何选项，则设置为"0"
    var rssListValues = $("#rss_list_id").val();
    if (!rssListValues || rssListValues.length === 0) {
        params.set("rss_list_id[]", "0");
    }
    // 处理epg_list_id - 如果未选择任何选项，则设置为"0"
    var epgListValues = $("#epg_list_id").val();
    if (!epgListValues || epgListValues.length === 0) {
        params.set("epg_list_id[]", "0");
    }

    params.append("saveChannelsOne_fenlei", "");
    var editchmodal = $(btn).closest('.modal');
    console.log('params',params)
    $.ajax({
        url: "/admin/fenlei",
        type: "POST",
        data: params.toString(),
        success: function(data) {
            lightyear.notify(data.msg, data.type, 3000);
            if (data.type === "success") {
                editchmodal.modal('hide');
                getChannels_pindao(params.get("caId"));
            }
        },
        error: function() {
            lightyear.notify("请求失败", "danger", 3000);
        }
    });
}
function getLogo(){
	var $select = $("#e_id");
    var logoSrc = $select.find("option:selected").data("value") || "";
    if (logoSrc) {
        $("#logoContainerEdit").show().find("img").attr("src", logoSrc);
    } else {
        $("#logoContainerEdit").hide().find("img").attr("src", "");
    }
}
function rssPOST(btn) {
	var $tr = $(btn).closest("tr"); 
	var mealid = $tr.find(".meal-id").data("value");
	const params = new URLSearchParams();
	params.append("id", mealid);
	fetch("/admin/getRssUrl", {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(res => {
		if (res.type === "success") {
			lightyear.notify(res.msg, res.type, 1000);
			res.data.forEach(function (item, index) {
				if (item.type === 'txt'){
					$("#getnewkey").val(mealid);
					$("#rsstxt").val(item.url);
				}else if (item.type === 'm3u8'){
					$("#rssm3u").val(item.url);
				}else if (item.type === 'epg'){
				    $("#rssepg").val(item.url);
				}
			});
		}else{
			lightyear.notify(res.msg, res.type, 3000);
		}
	})
}
function CopyRss(textarea) {
    textarea.select();
    document.execCommand("copy"); 
    lightyear.notify("✅ 已复制到剪贴板！", "success", 1000);
};
function uploadPayList(event) {
	var action = window.location.pathname;
	const file = event.target.files[0];
	if (!file) return;
	const formData = new FormData();
	const fileInput = document.querySelector('input[name="paylistfile"]'); 
	if (fileInput && fileInput.files.length > 0) {
	formData.append("paylistfile", fileInput.files[0]);
	}else{
		lightyear.notify("❌ 请选择文件！", "danger", 1000);
		return;
	}
	$.ajax({
		url: '/admin/channels/uploadPayList',  
		type: 'POST',
		data: formData,
		contentType: false,
		processData: false,
		success: function(res) {
			$('#paylistfile').val(''); 
			lightyear.notify(res.msg, res.type, 1000); 
			loadPage(action); 
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000);
			$('#paylistfile').val('');
		}
	});
};
function uploadEpg(event) {
	var input = event.target;
    var file = input.files[0];
    if (file) {
        $('#epgfilename').text(file.name); 
    } else {
        $('#epgfilename').text('');
    }
};
function uploadLogo(input) {
    if (!input.files || input.files.length === 0) return;
    var file = input.files[0];
    var tr = $(input).closest("tr");
    var rowName = tr.data("name");
    var formData = new FormData();
    formData.append("uploadlogo", file);
    formData.append("epgname", rowName); 
    $.ajax({
        url: "/admin/channels/uploadLogo",
        type: "POST",
        data: formData,
        processData: false, 
        contentType: false, 
        success: function(data) {
            if (typeof data === "string") {
                try {
                    data = JSON.parse(data);
                } catch (err) {
                    console.error("解析 JSON 失败:", err);
                    lightyear.notify("上传失败", "danger", 3000);
                    return;
                }
            }
			if (data && typeof data.msg === "string" && data.msg.includes('/admin/login')) {
                window.location.href = "/admin/login";
                return;
            }
            if (data && data.type === "success") {
                lightyear.notify(data.msg, data.type, 3000);
				const fileInputs = document.querySelectorAll('input[type="file"]');
    			fileInputs.forEach(input => input.value = "");
                loadPage(window.location.href);
            } else if (data && data.msg) {
                lightyear.notify(data.msg, data.type || "danger", 3000);
				const fileInputs = document.querySelectorAll('input[type="file"]');
    			fileInputs.forEach(input => input.value = "");
            }
        },
        error: function(err) {
            console.error(err);
            lightyear.notify("上传失败", "danger", 3000);
        }
    });
}
function getEpgList(btn) {
	var $tr = $(btn).closest("tr"); 
	var cid = $tr.find(".e-id").data("value");
	var cname = $tr.find(".e-name").data("value");
	var curl = $tr.find(".e-url").data("value");
	var cua = $tr.find(".e-ua").data("value");
	$("#eid").val(cid);
	$("#epgfromname").val(cname);
	$("#epgfromurl").val(curl);
	$("#epgfromua").val(cua);
}
function deleteLogo(id) {
	const params = new URLSearchParams();
	params.append("deleteLogo", id);
	$.ajax({
		url: '/admin/epgsList',   
		type: 'POST',
		data: params.toString(),
		contentType: 'application/x-www-form-urlencoded',
		success: function(res) {
			lightyear.notify(res.msg, res.type, 1000);
			if (res.code === 1) {
				$('#logo_'+id).remove();
				$('#uploadlogo').val('');
			}
		},
		error: function(res) {
			lightyear.notify(res.msg, res.type, 3000); 
			$('#uploadlogo').val(''); 
		}
	});
}
function getnewkey(btn){
	$.confirm({
        title: '操作确认',
        content: "刷新KEY后之前的链接将不可用，确认刷新吗？",
        type: 'green',
        buttons: {
            confirm: {
                text: '确认',
                btnClass: 'btn-success',
                action: function () {
					const params = new URLSearchParams();
					params.append("getnewkey", btn.value);
					$.ajax({
						url: '/admin/getRssUrl',   
						type: 'POST',
						data: params.toString(),
						contentType: 'application/x-www-form-urlencoded',
						success: function(res) {
							lightyear.notify(res.msg, res.type, 1000);
							if (res.type === "success") {
								lightyear.notify(res.msg, res.type, 1000);
								res.data.forEach(function (item, index) {
									if (item.type === 'txt'){
										$("#rsstxt").val(item.url);
									}else if (item.type === 'm3u8'){
										$("#rssm3u").val(item.url);
									}else if (item.type === 'epg'){
										$("#rssepg").val(item.url);
									}
								});
							}else{
								lightyear.notify(res.msg, res.type, 3000);
							}
						},
						error: function(res) {
							lightyear.notify(res.msg, res.type, 3000); 
						}
					});
                }
            },
            cancel: {
                text: '取消',
                btnClass: 'btn-danger'
            }
        }
    });
}
function BuildApk(btn) {
var form = btn.closest("form");
	if (!form) {
		lightyear.notify("表单提交失败", "danger", 3000);
		return;
	}
	var action = form.action || window.location.pathname; 
	const params = new URLSearchParams();
	new FormData(form).forEach((value, key) => {
		params.append(key, value);
	});
	params.append(btn.name, "");
	if (!params.has("up_sets")) {
		params.append("up_sets", 0);
	}
	fetch(action, {
		method: "POST",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		body: params.toString()
	})
	.then(async response => {
		const text = await response.text(); 
		if (text.includes('/admin/login')) {
			window.location.href = "/admin/login";
			throw text;
		}
		return JSON.parse(text); 
	})
	.then(data => {
		lightyear.notify(data.msg, data.type, 3000);
		if (data.type === "success") {
			$('#submitappinfo').prop('disabled', true);
			$('.download-link').prop('disabled', true);
			getBuildStatus()
		}
	})
	.catch(err => {
		lightyear.notify("提交失败:"+ err, "danger", 3000);
	});
}
function getBuildStatus() {
    var timer = setInterval(function() {
        $.getJSON('/admin/client/buildStatus', function(resp) {
            $('#apksize').val(resp.data.size);
            if(resp.code === 1) {
				lightyear.notify(resp.msg, resp.type, 1000);
				$('.download-link').each(function() {
					$(this).attr('href', resp.data.url);
					$(this).attr('download', resp.data.name);
				});
				$('#app_version').val(resp.data.version);
                $('#submitappinfo').prop('disabled', false);
				$('.download-link').prop('disabled', false);
                clearInterval(timer);
            }
        }).fail(function() {
			lightyear.notify("请求失败，稍后重试...", 'danger', 1000);
        });
    }, 1000); 
}
function toggleLock(btn) {
    const input = document.getElementById("serverUrl");
    const icon = btn.querySelector("i");
    if (input.hasAttribute("readonly")) {
		$.confirm({
			title: '操作确认',
			content: '修改APK连接地址后，需要重新构建APK，且之前APK可能出现网络连接失败，确认修改吗？',
			type: 'green',
			buttons: {
				confirm: {
					text: '确认',
					btnClass: 'btn-danger',
					action: function () {
						input.removeAttribute("readonly");
						input.focus();
						input.value = window.location.origin;
						icon.className = "mdi mdi-lock-open"; 
						btn.classList.remove("btn-danger");
						btn.classList.add("btn-primary"); 
					}
				},
				cancel: {
					text: '取消',
					btnClass: 'btn-success'
				}
			}
		});
    } else {
        input.setAttribute("readonly", true);
        icon.className = "mdi mdi-lock"; 
		btn.classList.remove("btn-primary");
        btn.classList.add("btn-danger"); 
    }
}
function caMoveup() {
	const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
		lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
	const tr = firstChecked.closest("tr");
    if (!tr) {
		lightyear.notify("未找到对应行", 'danger', 1000);
        return;
	};
	const prevTr = tr.previousElementSibling;
    if (!prevTr) {
		lightyear.notify("已经是第一行了", 'danger', 1000);
		return;
	}; 
    const value = firstChecked.value;
	const params = new URLSearchParams();
	params.append("moveup", value);
    $.ajax({
        url: "/admin/channels",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
				tbody.insertBefore(tr, prevTr); 
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
			lightyear.notify("操作失败", 'danger', 1000);
		}
    });
}
function caMoveup_fenlei() {
    const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
        lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
    const tr = firstChecked.closest("tr");
    if (!tr) {
        lightyear.notify("未找到对应行", 'danger', 1000);
        return;
    };
    const prevTr = tr.previousElementSibling;
    if (!prevTr) {
        lightyear.notify("已经是第一行了", 'danger', 1000);
        return;
    };
    const value = firstChecked.value;
    const params = new URLSearchParams();
    params.append("moveup", value);
    $.ajax({
        url: "/admin/fenlei",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
                tbody.insertBefore(tr, prevTr);
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
            lightyear.notify("操作失败", 'danger', 1000);
        }
    });
}
function caMovedown() {
    const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
        lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
    const tr = firstChecked.closest("tr");
    if (!tr) {
        lightyear.notify("未找到对应行", 'danger', 1000);
        return;
    }
    const nextTr = tr.nextElementSibling;
    if (!nextTr) {
        lightyear.notify("已经是最后一行了", 'danger', 1000);
        return; 
    }
    const value = firstChecked.value;
    const params = new URLSearchParams();
    params.append("movedown", value);
    $.ajax({
        url: "/admin/channels",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
                tbody.insertBefore(tr, nextTr.nextElementSibling);
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
            lightyear.notify("操作失败", 'danger', 1000);
        }
    });
}
function caMovedown_fenlei() {
    const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
        lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
    const tr = firstChecked.closest("tr");
    if (!tr) {
        lightyear.notify("未找到对应行", 'danger', 1000);
        return;
    }
    const nextTr = tr.nextElementSibling;
    if (!nextTr) {
        lightyear.notify("已经是最后一行了", 'danger', 1000);
        return;
    }
    const value = firstChecked.value;
    const params = new URLSearchParams();
    params.append("movedown", value);
    $.ajax({
        url: "/admin/fenlei",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
                tbody.insertBefore(tr, nextTr.nextElementSibling);
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
            lightyear.notify("操作失败", 'danger', 1000);
        }
    });
}
function caMovetop() {
    const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
        lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
    const tr = firstChecked.closest("tr");
    if (!tr) {
        lightyear.notify("未找到对应行", 'danger', 1000);
        return;
    }
    if (tr === tbody.firstElementChild) {
        lightyear.notify("已经在最顶端了", 'danger', 1000);
        return;
    }
    const value = firstChecked.value;
    const params = new URLSearchParams();
    params.append("movetop", value);
    $.ajax({
        url: "/admin/channels",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
                tbody.insertBefore(tr, tbody.firstElementChild);
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
            lightyear.notify("操作失败", 'danger', 1000);
        }
    });
}
function caMovetop_fenlei() {
    const tbody = document.getElementById("categorylist_tbody");
    if (!tbody) return;
    const firstChecked = tbody.querySelector("input[type='checkbox']:checked");
    if (!firstChecked) {
        lightyear.notify("请先选择一个分类", 'danger', 1000);
        return;
    }
    const tr = firstChecked.closest("tr");
    if (!tr) {
        lightyear.notify("未找到对应行", 'danger', 1000);
        return;
    }
    if (tr === tbody.firstElementChild) {
        lightyear.notify("已经在最顶端了", 'danger', 1000);
        return;
    }
    const value = firstChecked.value;
    const params = new URLSearchParams();
    params.append("movetop", value);
    $.ajax({
        url: "/admin/fenlei",
        type: "POST",
        data: params.toString(),
        success: function (data) {
            if (data.type === "success") {
                lightyear.notify(data.msg, data.type, 1000);
                tbody.insertBefore(tr, tbody.firstElementChild);
            } else {
                lightyear.notify(data.msg, data.type, 1000);
            }
        },
        error: function () {
            lightyear.notify("操作失败", 'danger', 1000);
        }
    });
}
