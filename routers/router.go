package routers

import (
	"bee/activist/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{}, "get:ToHome")
    beego.Router("/home", &controllers.MainController{})
    beego.Router("/register", &controllers.MainController{}, "get,post:Register")
    beego.Router("/login/:back", &controllers.MainController{}, "get,post:Login")
    beego.Router("/logout", &controllers.MainController{}, "get:Logout")
    beego.Router("/events/new", &controllers.MainController{}, "get,post:NewEvent")
    beego.Router("/events/edit/:id([0-9]+)", &controllers.MainController{}, "get,post:EditEvent")
    beego.Router("/events/view/:id([0-9]+)", &controllers.MainController{}, "get:ViewEvent")
    beego.Router("/events/join/:id([0-9]+)", &controllers.MainController{}, "get:JoinEvent")
    beego.Router("/events/delete/:id([0-9]+)", &controllers.MainController{}, "get:DeleteEvent")
    beego.Router("/events/deny/:id([0-9]+)", &controllers.MainController{}, "get:DenyEvent")
    beego.Router("/profile", &controllers.MainController{}, "get:Profile")
}
