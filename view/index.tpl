<!doctype html>
<html>
  	<head>
    	<title>{{ .title }}</title>
    	<meta http-equiv="content-type" content="text/html; charset=utf-8">
		
	</head>
<style>
table,tr,td{
	border:1px #000000 solid;
}
</style>
<body>


==================以下是public接口(不需要登录)====================
<br />


<p>1、注册接口<br />POST<br />
/public/reg.do</p>
<table>
<form method="post" action="/public/reg.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>用户名:必填(必须字母开头,长度是6-15位的字母和数字组成)</td>
</tr>
<tr>
<td>pwd</td>
<td><input type="text" name="pwd"></td>
<td>密码:必填(长度是6-15位的字母或数字组成)</td>
</tr>
<tr>
<td>reg_code</td>
<td><input type="text" name="reg_code"></td>
<td>注册邀请码:可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>2、判断用户名是否可用<br />POST<br />
/public/account_allow.do</p>
<table>
<form method="post" action="/public/account_allow.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>账号:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>3、判断密码是否可用<br />POST<br />
/public/pwd_allow.do</p>
<table>
<form method="post" action="/public/pwd_allow.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pwd</td>
<td><input name="pwd" type="text"></td>
<td>密码:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>4、登录<br />POST<br />
/public/login.do</p>
<table>
<form method="post" action="/public/login.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>用户名:必填</td>
</tr>
<tr>
<td>pwd</td>
<td><input type="text" name="pwd"></td>
<td>密码:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>5、退出<br />POST<br />
/public/logout.do</p>
<table>
<form method="post" action="/public/logout.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>6、启动页<br />POST<br />
/public/loading.do</p>
<table>
<form method="post" action="/public/loading.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>7、轮播图<br />POST<br />
/public/banner.do</p>
<table>
<form method="post" action="/public/banner.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>8、弹窗公告<br />POST<br />
/public/pop.do</p>
<table>
<form method="post" action="/public/pop.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>9、视频标签类型<br />对应设计图的导航栏<br />POST<br />
/public/tag_type.do</p>
<table>
<form method="post" action="/public/tag_type.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_top</td>
<td><input type="text" name="is_top"></td>
<td>是否顶部:可空(1=顶部显示的,0=非顶部显示的,不填=所有标签类型)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>10、热门推荐视频(根据观看次数排序)<br />POST<br />
/public/hot_video.do</p>
<table>
<form method="post" action="/public/hot_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_hot</td>
<td><input type="text" name="is_hot"></td>
<td>是否热门:可空(1=热门,0=非热门,不填=所有标签类型)</td>
</tr>
<tr>
<td>is_recomm</td>
<td><input type="text" name="is_recomm"></td>
<td>是否推荐:可空(1=推荐,0=非推荐,不填=所有标签类型)</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>11、通过标签类型查询视频(根据观看次数排序)<br />设计图:点击导航栏得到下面的视频列表,视频播放地址等信息需要查询视频详情获取/user/video_info.do<br />POST<br />
/public/type_video.do</p>
<table>
<form method="post" action="/public/type_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>type_id</td>
<td><input type="text" name="type_id"></td>
<td>标签类型ID:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>12、通过标签查询视频(根据观看次数排序)<br />设计图没有对应页面,可以先不用<br />POST<br />
/public/tag_video.do</p>
<table>
<form method="post" action="/public/tag_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>tag_id</td>
<td><input type="text" name="tag_id"></td>
<td>标签ID:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<!-- 弃用
<p>14、女优列表(根据视频数量排序)<br />POST<br />
/public/actor_list.do</p>
<table>
<form method="post" action="/public/actor_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_hot</td>
<td><input type="text" name="is_hot"></td>
<td>是否热门:可空(1=热门,0=非热门,不填=所有女优)</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
-->

<!-- 弃用
<p>15、游戏平台类型列表(根据ID排序)<br />POST<br />
/public/plat_type.do</p>
<table>
<form method="post" action="/public/plat_type.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
-->

<!-- 弃用
<p>16、根据平台类型查询平台列表(根据sort排序)<br />POST<br />
/public/plat_list.do</p>
<table>
<form method="post" action="/public/plat_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>type_code</td>
<td><input type="text" name="type_code"></td>
<td>平台类型编码:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
-->


<p>17、关键字搜索(根据sort排序)<br />POST<br />
/public/kw_search.do</p>
<table>
<form method="post" action="/public/kw_search.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>key_word</td>
<td><input type="text" name="key_word"></td>
<td>平台类型编码:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<!-- 弃用
<p>18、活动类型列表(根据sort排序)<br />POST<br />
/public/active_type.do</p>
<table>
<form method="post" action="/public/active_type.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
-->


<!-- 弃用
<p>19、活动列表(根据sort排序)<br />POST<br />
/public/active_list.do</p>
<table>
<form method="post" action="/public/active_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>type_id</td>
<td><input type="text" name="type_id"></td>
<td>活动类型ID:可空</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
-->


<p>20、客服联系方式<br />POST<br />
/public/cs_info.do</p>
<table>
<form method="post" action="/public/cs_info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>21、获取app的版本信息<br />POST<br />
/public/app_ver.do</p>
<table>
<form method="post" action="/public/app_ver.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>22、app的每次启动请求接口(包括第一次安装)<br />POST<br />
/public/open_app.do</p>
<table>
<form method="post" action="/public/open_app.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input type="text" name="account"></td>
<td>用户名:可空</td>
</tr>
<tr>
<td>phone_info</td>
<td><input type="text" name="phone_info"></td>
<td>手机的型号信息:必填</td>
</tr>
<tr>
<td>phone_os</td>
<td><input type="text" name="phone_os"></td>
<td>手机操作系统信息:必填</td>
</tr>
<tr>
<td>uid</td>
<td><input type="text" name="uid"></td>
<td>手机唯一编码信息:必填</td>
</tr>
<tr>
<td>phone_number</td>
<td><input type="text" name="phone_number"></td>
<td>手机的电话号码:可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>23、app的每次关闭请求接口<br />POST<br />
/public/close_app.do</p>
<table>
<form method="post" action="/public/close_app.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>uid</td>
<td><input type="text" name="uid"></td>
<td>手机唯一编码信息:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>24、标签列表<br />POST<br />
/public/tag_list.do</p>
<table>
<form method="post" action="/public/tag_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>25、女优资料<br />POST<br />/public/actor_data.do</p>
<table>
<form method="post" action="/public/actor_data.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>26、女优资料详情;通过女优查询视频(根据观看次数排序)<br />POST<br />
/public/actor_video.do</p>
<table>
<form method="post" action="/public/actor_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>actor_id</td>
<td><input type="text" name="actor_id"></td>
<td>女优ID:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>27、购买积分的比例<br />POST<br />
/public/rate.do</p>
<table>
<form method="post" action="/public/rate.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>28、获取茄子主推视频(根据sort排序)<br />POST<br />
/public/qz_recomm.do</p>
<table>
<form method="post" action="/public/qz_recomm.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>29、茄子主推短视频的观看次数+1<br />POST<br />/public/add_view.do</p>
<table>
<form method="post" action="/public/add_view.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>id</td>
<td><input type="text" name="id"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>30、影片的观看次数+1,以及获取影片的总观看次数<br />POST<br />/public/add_see.do</p>
<table>
<form method="post" action="/public/add_see.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>31、首页-标题栏,每个标题栏下面取4个视频<br />POST<br />/public/home_view.do</p>
<table>
  <form method="post" action="/public/home_view.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<br /><br /><br /><br />
==================以下是user接口====================
<br />


<p>1、查询用户信息<br />POST<br />
/user/user_info.do</p>
<table>
<form method="post" action="/user/user_info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>2、视频详情页信息<br />POST<br />/user/video_info.do</p>
<table>
<form method="post" action="/user/video_info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>3、视频详情页的推荐视频<br />POST<br />/user/recomm_video.do</p>
<table>
<form method="post" action="/user/recomm_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>4、视频详情页的评论<br />POST<br />/user/video_comm.do</p>
<table>
<form method="post" action="/user/video_comm.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr>
<td>p_id</td>
<td><input type="text" name="p_id"></td>
<td>上级评论的ID:可空,默认0</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>5、评论视频<br />POST<br />/user/comm.do</p>
<table>
<form method="post" action="/user/comm.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr>
<td>p_id</td>
<td><input type="text" name="p_id"></td>
<td>上级评论的ID:可空,默认0</td>
</tr>
<tr>
<td>content</td>
<td><input type="text" name="content"></td>
<td>评论内容:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>6、点赞/取消点赞视频,以及获取点赞总次数<br />POST<br />/user/like_video.do</p>
<table>
<form method="post" action="/user/like_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>7、点赞/取消点赞评论,以及获取点赞总次数<br />POST<br />/user/like_comm.do</p>
<table>
<form method="post" action="/user/like_comm.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>comm_id</td>
<td><input type="text" name="comm_id"></td>
<td>评论ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>8、收藏/取消收藏视频<br />POST<br />/user/star_video.do</p>
<table>
<form method="post" action="/user/star_video.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>vid</td>
<td><input type="text" name="vid"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>9、查看自己的评论<br />POST<br />/user/user_comm.do</p>
<table>
<form method="post" action="/user/user_comm.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>10、获取余额&积分<br />POST<br />/user/user_money.do</p>
<table>
<form method="post" action="/user/user_money.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>11、上传头像<br />POST<br />/user/user_avatar.do</p>
<table>
<form method="post" action="/user/user_avatar.do" enctype="multipart/form-data">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>avatar</td>
<td><input type="file" name="avatar"></td>
<td>头像:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>12、修改昵称<br />POST<br />/user/user_nickname.do</p>
<table>
<form method="post" action="/user/user_nickname.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>nick_name</td>
<td><input type="text" name="nick_name"></td>
<td>昵称:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>13、修改年龄<br />POST<br />/user/user_age.do</p>
<table>
<form method="post" action="/user/user_age.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>age</td>
<td><input type="text" name="age"></td>
<td>年龄:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>14、修改性别<br />POST<br />/user/user_sex.do</p>
<table>
<form method="post" action="/user/user_sex.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>sex</td>
<td><input type="text" name="sex"></td>
<td>性别:必填(性别:0=男,1=女)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>15、修改个性签名<br />POST<br />/user/user_persign.do</p>
<table>
<form method="post" action="/user/user_persign.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>per_sign</td>
<td><input type="text" name="per_sign"></td>
<td>个性签名:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>16、修改密码<br />POST<br />/user/user_editpwd.do</p>
<table>
<form method="post" action="/user/user_editpwd.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>

</tr>
<tr>
<td>old_pwd</td>
<td><input type="text" name="old_pwd"></td>
<td>旧密码:必填</td>
</tr>

<tr>
<td>pwd</td>
<td><input type="text" name="pwd"></td>
<td>新密码:必填</td>
</tr>

<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>17、获取茄子钱包充值、优惠、消费记录<br />POST<br />
/user/qz_wallet.do</p>
<table>
<form method="post" action="/user/qz_wallet.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_type</td>
<td><input type="text" name="order_type"></td>
<td>订单类型:必填(1=充值记录,3=优惠记录,5=消费记录)</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>18、获取kok币的购买、优惠、消费记录<br />POST<br />
/user/kok_gold.do</p>
<table>
<form method="post" action="/user/kok_gold.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_type</td>
<td><input type="text" name="order_type"></td>
<td>订单类型:必填(5=获取积分记录,3=优惠记录,6=消费积分记录)</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>19、查看自己的收藏<br />POST<br />/user/user_collect.do</p>
<table>
<form method="post" action="/user/user_collect.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>20、删除指定的收藏视频<br />POST<br />/user/del_collect.do</p>
<table>
  <form method="post" action="/user/del_collect.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    
    <tr>
    <td>video_id</td>
    <td><input type="text" name="video_id[]"></td>
    <td>视频ID 1</td>
    </tr>
    
    <tr>
      <td>video_id</td>
      <td><input type="text" name="video_id[]"></td>
      <td>视频ID 2</td>
    </tr>
     
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>21、清空收藏<br />POST<br />/user/clear_collect.do</p>
<table>
<form method="post" action="/user/clear_collect.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>22、查看观看历史<br />POST<br />/user/user_history.do</p>
<table>
<form method="post" action="/user/user_history.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>23、删除指定的观看历史<br />POST<br />/user/del_history.do</p>
<table>
  <form method="post" action="/user/del_history.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    
    <tr>
    <td>video_id</td>
    <td><input type="text" name="video_id[]"></td>
    <td>视频ID 1</td>
    </tr>
    
    <tr>
      <td>video_id</td>
      <td><input type="text" name="video_id[]"></td>
      <td>视频ID 2</td>
    </tr>
     
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>

<p>24、清空观看历史<br />POST<br />/user/clear_history.do</p>
<table>
<form method="post" action="/user/clear_history.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>25、查看未读消息数量<br />POST<br />/user/unread_num.do</p>
<table>
<form method="post" action="/user/unread_num.do">
<tr>

<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>

<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>26、消息中心消息<br />POST<br />/user/sys_msg.do</p>
<table>
<form method="post" action="/user/sys_msg.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>msg_type</td>
<td><input type="text" name="msg_type"></td>
<td>消息类型:必填(0=系统消息,1=评论回复,2=评论点赞)</td>
</tr>
<tr>
<td>page</td>
<td><input type="text" name="page"></td>
<td>当前页:可空(默认第一页)</td>
</tr>
<tr>
<td>page_size</td>
<td><input type="text" name="page_size"></td>
<td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>27、删除指定消息中心消息<br />POST<br />/user/del_msg.do</p>
<table>
  <form method="post" action="/user/del_msg.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    
    <tr>
    <td>msg_id</td>
    <td><input type="text" name="msg_id[]"></td>
    <td>消息ID,即sys_msg.do接口返回的id</td>
    </tr>
    
    <tr>
    <td>msg_id</td>
    <td><input type="text" name="msg_id[]"></td>
    <td>消息ID,即sys_msg.do接口返回的id</td>
    </tr>
     
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>28、把此消息标记为已读<br />POST<br />/user/view_msg.do</p>
<table>
<form method="post" action="/user/view_msg.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>msg_id</td>
<td><input type="text" name="msg_id"></td>
<td>消息ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>29、用户点赞/取消点赞茄子主推视频以及查看视频点赞总次数<br />POST<br />/user/like_theme.do</p>
<table>
<form method="post" action="/user/like_theme.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>id</td>
<td><input type="text" name="id"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>30、查看自己是否点赞茄子主推视频以及查看视频点赞总次数<br />POST<br />/user/see_like.do</p>
<table>
<form method="post" action="/user/see_like.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>id</td>
<td><input type="text" name="id"></td>
<td>视频ID:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>31、获取初次绑定手机验证码/获取解绑手机验证码<br />POST<br />/user/phone_verify.do</p>
<table>
<form method="post" action="/user/phone_verify.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>phone</td>
<td><input type="text" name="phone"></td>
<td>手机号:必填(初次绑定:手机号为需绑定的手机号,换绑:手机号为需解绑的手机号)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>32、校验手机验证码<br />POST<br />/user/phone_bind.do</p>
<table>
<form method="post" action="/user/phone_bind.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>code</td>
<td><input type="text" name="code"></td>
<td>验证码:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>33、获取换绑手机验证码<br />POST<br />/user/phone_chg.do</p>
<table>
<form method="post" action="/user/phone_chg.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>phone</td>
<td><input type="text" name="phone"></td>
<td>手机号:必填(手机号为需要绑定的手机号)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>



<br /><br /><br /><br />
==================以下是task(任务)接口====================
<br />


<p>1、我的任务-签到记录<br />POST<br />/task/checkin_log.do</p>
<table>
<form method="post" action="/task/checkin_log.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>2、我的任务-点击签到<br />POST<br />/task/checkin_click.do</p>
<table>
<form method="post" action="/task/checkin_click.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>3、我的任务-任务列表<br />POST<br />/task/task_list.do</p>
<table>
<form method="post" action="/task/task_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>4、我的任务-完成任务<br />POST<br />/task/task_click.do</p>
<table>
<form method="post" action="/task/task_click.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>id</td>
<td><input type="text" name="id"></td>
<td>任务id</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<br /><br /><br /><br />
=======================财务相关的接口=======================
<br />


<p>1、查询支付方式列表的接口<br />POST<br />
/finance/pay_list.do</p>
<table>
<form method="post" action="/finance/pay_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>2、支付接口<br />POST<br />
/finance/pay.do</p>
<table>
<form method="post" action="/finance/pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_id</td>
<td><input type="text" name="pay_id"></td>
<td>pay_list接口中的ID:必填</td>
</tr>
<tr>
<td>amount</td>
<td><input type="text" name="amount"></td>
<td>额度:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>3、购买积分<br />POST<br />
/finance/score.do</p>
<table>
<form method="post" action="/finance/score.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>score</td>
<td><input type="text" name="score"></td>
<td>需要购买的积分:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>




<br /><br /><br /><br />
==================以下是blog(广场)接口====================
<br />

<p>051、获取博客列表<br />POST<br />/blog/blog_list.do</p>
<table>
  <form method="post" action="/blog/blog_list.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>is_top</td>
      <td><input type="text" name="is_top" value="" ></td>
      <td>是否置顶:0=否,1=是,不传(空)=所有</td>
    </tr>
    <tr>
      <td>is_good</td>
      <td><input type="text" name="is_good" value="" ></td>
      <td>是否精品:0=否,1=是,不传(空)=所有</td>
    </tr>
    <tr>
      <td>content</td>
      <td><input type="text" name="content" value="" ></td>
      <td>搜索帖子,非必填</td>
    </tr>
    <tr>
      <td>page</td>
      <td><input type="text" name="page"></td>
      <td>当前页:可空(默认第一页)</td>
    </tr>
    <tr>
      <td>page_size</td>
      <td><input type="text" name="page_size"></td>
      <td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>052、获取博客列表的全部 评论,回复; 评论=在博客下评论, 回复=对博客下的评论进行回复<br />POST<br />/blog/blog_comm.do</p>
<table>
  <form method="post" action="/blog/blog_comm.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>blog_id</td>
      <td><input type="text" name="blog_id" value="" ></td>
      <td>blog_list.do 返回的id</td>
    </tr>
    <tr>
      <td>comment_id</td>
      <td><input type="text" name="comment_id" value="" ></td>
      <td>当前id; 博客下评论里面的二级评论</td>
    </tr>
    <tr>
      <td>page</td>
      <td><input type="text" name="page"></td>
      <td>当前页:可空(默认第一页)</td>
    </tr>
    <tr>
      <td>page_size</td>
      <td><input type="text" name="page_size"></td>
      <td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>053、 评论博客,或者回复评论<br />POST<br />/blog/blog_comm_reply.do</p>
<table>
  <form method="post" action="/blog/blog_comm_reply.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>blog_id</td>
      <td><input type="text" name="blog_id" value="" ></td>
      <td>blog_list.do 返回的id</td>
    </tr>
    <tr>
      <td>comment_id</td>
      <td><input type="text" name="comment_id" value="" ></td>
      <td>评论id,如果是在博客下面评论=不需要,如果实在评论下面回复=需要评论id</td>
    </tr>
    <tr>
      <td>content</td>
      <td><input type="text" name="content" value="" ></td>
      <td>评论内容</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>054、 博客 点赞/取消<br />POST<br />/blog/blog_like.do</p>
<table>
  <form method="post" action="/blog/blog_like.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>blog_id</td>
      <td><input type="text" name="blog_id" value="" ></td>
      <td>blog_list.do 返回的id</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>055、 评论 点赞/取消<br />POST<br />/blog/comm_like.do</p>
<table>
  <form method="post" action="/blog/comm_like.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>comm_id</td>
      <td><input type="text" name="comm_id" value="" ></td>
      <td>blog_comm.do 返回的id</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>

<!--
<p>056、 批量上传图片<br />POST<br />/blog/upload_img_more.do</p>
<table>
  <form method="post" action="/blog/upload_img_more.do" enctype="multipart/form-data">
    <tr>
    <td>参数名</td>
    <td>值</td>
    <td>描述</td>
    </tr>
    
    <tr>
    <td>img</td>
    <td><input type="file" name="img[]">图片1</td>
    <td>图片</td>
    </tr>
    
    <tr>
    <td>img</td>
    <td><input type="file" name="img[]">图片2</td>
    <td>图片</td>
    </tr>   
     
    <tr>
    <td>img</td>
    <td><input type="file" name="img[]">图片3</td>
    <td>图片</td>
    </tr>    
    
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>-->


<p>057、 发帖<br />POST<br />/blog/blog_publish.do</p>
<table>
  <form method="post" action="/blog/blog_publish.do" enctype="multipart/form-data">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>

    <tr>
      <td>content</td>
      <td><input type="text" name="content" value="" ></td>
      <td>帖子内容</td>
    </tr>
    
    <tr>
    <td>img</td>
    <td><input type="file" name="img[]">图片1</td>
    <td>图片</td>
    </tr>
    
    <tr>
      <td>img</td>
      <td><input type="file" name="img[]">图片2</td>
      <td>图片</td>
    </tr>
     
    <tr>
      <td>img</td>
      <td><input type="file" name="img[]">图片3</td>
      <td>图片</td>
    </tr>

    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>058 获取某个博客信息<br />POST<br />/blog/blog_one.do</p>
<table>
  <form method="post" action="/blog/blog_one.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>id</td>
      <td><input type="text" name="id" value="" ></td>
      <td>博客id</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>059 获取某个评论的信息;正常来说只获取p_id=0的评论id,评论里面的回复不需要刷新获取<br />POST<br />/blog/blog_comm_one.do</p>
<table>
  <form method="post" action="/blog/blog_comm_one.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>comment_id</td>
      <td><input type="text" name="comment_id" value="" ></td>
      <td>评论id; blog_comm.do返回的id</td>
    </tr>
    <tr>
      <td>page</td>
      <td><input type="text" name="page"></td>
      <td>当前页:可空(默认第一页)</td>
    </tr>
    <tr>
      <td>page_size</td>
      <td><input type="text" name="page_size"></td>
      <td>每页显示数量:可空(默认20条,最少1条,最多100条)</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>


<p>060、 删除博客<br />POST<br />/blog/del_blog.do</p>
<table>
  <form method="post" action="/blog/del_blog.do">
    <tr>
      <td>参数名</td>
      <td>值</td>
      <td>描述</td>
    </tr>
    <tr>
      <td>blog_id</td>
      <td><input type="text" name="blog_id" value="" ></td>
      <td>blog_list.do 返回的id</td>
    </tr>
    <tr><td><input type="submit" value="ok"></td></tr>
  </form>
</table>



</body>
</html>