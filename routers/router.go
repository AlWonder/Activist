package routers

import (
	"bee/activist/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/*", &controllers.MainController{})

    // JSON data
    beego.Router("/events", &controllers.MainController{}, "get:QueryEvents")
		beego.Router("/events/:id([0-9]+", &controllers.MainController{}, "get:GetEvent")
		beego.Router("/login", &controllers.MainController{}, "post:Login")
		beego.Router("/signup", &controllers.MainController{}, "post:SignUp")
}
