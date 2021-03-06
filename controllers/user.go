package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"mybeego/models"
	"regexp"
	"strconv"
)

type UserController struct {
	beego.Controller
}

/*****************************Register*****************************/

/**
显示注册页面
 */
func (this *UserController) ShowReg() {
	this.TplName="register.html"
}

/**
处理注册数据
 */
func (this *UserController)HandleReg()  {
	//获取数据
	userName:= this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")

	//数据校验
	if userName == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整,请重新注册"
		this.TplName = "register.html"
		return
	}
	if pwd != cpwd {
		this.Data["errmsg"] = "两次输入密码不一致,请重新输入"
		this.TplName = "register.html"
		return
	}
	reg, _ := regexp.Compile("^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z0-9]{2,6}$")
	res := reg.FindString(email)

	if res == ""{
		this.Data["errmsg"] = "邮箱格式不正确,请重新输入"
		this.TplName = "register.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	user.Email = email

	_, err := o.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "注册失败,请重新输入"
		this.TplName = "register.html"
		return
	}

	//发送邮件
	emailConfig := `{"username":"563364657@qq.com","password":"cgapyzgkkczubdea","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "563364657@qq.com"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	//注意这里我们发送给用户的是激活请求地址
	emailConn.Text = "http://192.168.23.27:8080/active?id="+strconv.Itoa(user.Id)

	emailConn.Send()

	//返回视图
	this.Ctx.WriteString("注册成功,请去相应邮箱激活用户")
}

/**
邮箱激活
 */
func (this *UserController) ActiveUser()  {
	//获取数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}

	//更新操作
	o:= orm.NewOrm()
	var user models.User
	user.Id = id
	err = o.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}

	user.Active = true
	o.Update(&user)

	//返回视图
	this.Redirect("/login",302)
}

func (this *UserController) ShowLogin()  {

	userName := this.Ctx.GetCookie("userName")

	//解码
	temp, _ := base64.StdEncoding.DecodeString(userName)
	if string(temp) == "" {
		this.Data["userName"] =""
		this.Data["checked"] = ""
	}else {
		this.Data["userName"] = string(temp)
		this.Data["checked"] = "checked"
	}

	this.TplName="login.html"
}

/**************************Login***************************/
//处理登录业务
func (this *UserController)HandleLogin()  {

	//获取数据
	userName := this.GetString("username")
	pwd := this.GetString("pwd")

	//校验数据
	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "登录数据不完整,请重新输入!"
		this.TplName = "login.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName

	err := o.Read(&user,"Name")
	if err != nil {
		this.Data["errmsg"] = "用户名或密码错误,请重新输入"
		this.TplName = "login.html"
		return
	}

	if user.PassWord != pwd {
		this.Data["errmsg"] = "用户名或密码错误,请重新输入"
		this.TplName = "login.html"
		return
	}

	if user.Active != true {
		this.Data["errmsg"] = "用户未激活,请前往邮箱激活"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")

	if remember == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",temp,7*24*60*60)
	}else {
		this.Ctx.SetCookie("userName","",-1)
	}

	this.SetSession("userName",userName)

	//返回视图
	//this.Ctx.WriteString("登录成功")
	this.Redirect("/",302)
}

//退出登录
func (this *UserController) Logout()  {
	this.DelSession("userName")
	this.Redirect("/login",302)
}

//展示用户中心信息页
func (this *UserController)ShowUserCenterInfo()  {
	userName := GetUserName(&this.Controller)

	this.Data["userName"] = userName

	o := orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name",userName).Filter("Isdefault",true).One(&addr)

	if addr.Id == 0 {
		this.Data["addr"]=""
	}else {
		this.Data["addr"] = addr
	}
	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_info.html"
}

//用户中心订单页
func (this * UserController)ShowUserCenterOrder()  {
	GetUserName(&this.Controller)
	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_order.html"
}

//用户中心地址页
func (this *UserController)ShowUserCenterSite()  {
	userName := GetUserName(&this.Controller)
	this.Data["userName"] = userName
	//获取地址信息
	o := orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name",userName).Filter("Isdefault",true).One(&addr)

	//传递给视图
	this.Data["addr"] = addr
	this.Layout="userCenterLayout.html"
	this.TplName="user_center_site.html"
}

//处理用户中心地址数据
func (this *UserController)HandleUserCenterSite()  {
	//获取数据
	receiver := this.GetString("receiver")
	addr := this.GetString("addr")
	zipCode := this.GetString("zipCode")
	phone := this.GetString("phone")
	//校验数据
	if receiver == "" || addr == "" || zipCode == "" || phone == "" {
		beego.Info("添加数据不完整")
		this.Redirect("/user/userCenterSite",302)
		return
	}

	//处理数据
	//插入操作
	o := orm.NewOrm()
	var addrUser models.Address
	addrUser.Isdefault=true
	err := o.Read(&addrUser,"Isdefault")
	//添加默认地址之前需要把原来的默认地址更新成非默认地址
	if err == nil {
		addrUser.Isdefault=false
		o.Update(&addrUser)
	}

	//更新默认地址时，给原来的地址对象的ID赋值了，这时候用原来的地址对象插入，意思时用原来的ID做插入操作,会报错
	//关联
	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")

	var addUserNew models.Address

	addUserNew.Receiver = receiver
	addUserNew.Zipcode = zipCode
	addUserNew.Addr = addr
	addUserNew.Phone = phone
	addUserNew.Isdefault = true
	addUserNew.User = &user
	o.Insert(&addUserNew)

	//返回视图
	this.Redirect("/user/userCenterSite",302)
}

