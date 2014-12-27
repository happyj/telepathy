// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	. "telepathy/Utils"
	"telepathy/controllers"
	"telepathy/models"
)

func init() {
	//ns := beego.NewNamespace("/v1",
	//	beego.NSNamespace("/object",
	//		beego.NSInclude(
	//			&controllers.ObjectController{},
	//		),
	//	),
	//	beego.NSNamespace("/user",
	//		beego.NSInclude(
	//			&controllers.UserController{},
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)

	Run()
}

func Run() {
	//初始化
	initialize()

	fmt.Println("Starting....")

	fmt.Println("Start ok")
}
func initialize() {
	//判断初始化参数
	initArgs()

	models.Connect()

	router()
	beego.AddFuncMap("stringsToJson", StringsToJson)
}
func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
	}
}
func router() {
	beego.Router("/telepathy/user/login", &controllers.UserController{}, "*:Register")
	beego.Router("/telepathy/user/logout", &controllers.UserController{}, "*:RegisterEnd")
	beego.Router("/telepathy/user/changepwd", &controllers.UserController{}, "*:RegisterEnd")
	beego.Router("/telepathy/user/register", &controllers.UserController{}, "*:Register")
	beego.Router("/telepathy/user/registerend", &controllers.UserController{}, "*:RegisterEnd")
	beego.Router("/telepathy/user/updateprofile", &controllers.UserController{}, "*:UpdateProfile")
	beego.Router("/telepathy/user/deluser", &controllers.UserController{}, "*:DelUser")

}
