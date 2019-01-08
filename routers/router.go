package routers

import (
	"github.com/astaxie/beego/context"
	"mybeego/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/user/*",beego.BeforeExec,fliterFunc)
    //注册
	beego.Router("/register", &controllers.UserController{},"get:ShowReg;post:HandleReg")
    //激活
	beego.Router("/active", &controllers.UserController{},"get:ActiveUser")
    //登录
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	//跳转首页
	beego.Router("/", &controllers.GoodsController{},"get:ShowIndex")

	//商品详情展示
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowGoodsDetail")
	//商品列表页
	beego.Router("/goodsList",&controllers.GoodsController{},"get:ShowList")
	//商品搜索
	beego.Router("/goodsSearch",&controllers.GoodsController{},"post:HandleSearch")

	//退出登录
	beego.Router("/user/logout",&controllers.UserController{},"get:Logout")
	//用户中心信息页
	beego.Router("/user/userCenterInfo",&controllers.UserController{},"get:ShowUserCenterInfo")
	//用户中心订单页
	beego.Router("/user/userCenterOrder",&controllers.UserController{},"get:ShowUserCenterOrder")
	//用户中心地址页
	beego.Router("/user/userCenterSite",&controllers.UserController{},"get:ShowUserCenterSite;post:HandleUserCenterSite")
}

var fliterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302,"/login")
		return
	}
}


