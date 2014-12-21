package controllers

import (
	"github.com/astaxie/beego"
)

type CommonController struct {
	beego.Controller
}

func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJson()
}

func init() {
}
