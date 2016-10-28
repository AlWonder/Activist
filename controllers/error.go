package controllers

import (
    "github.com/astaxie/beego"
)

type ErrorController struct {
    beego.Controller
}

func (c *ErrorController) Error404() {
    c.TplName = "errors/404.tpl"
}