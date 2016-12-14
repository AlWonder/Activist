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
    beego.Router("/events/participants/:id([0-9]+)", &controllers.MainController{}, "get:ViewParticipants")
    beego.Router("/events/delete/:id([0-9]+)", &controllers.MainController{}, "get:DeleteEvent")
    beego.Router("/events/deny/:id([0-9]+)", &controllers.MainController{}, "get:DenyEvent")
    beego.Router("/profile", &controllers.MainController{}, "get:Profile")
    beego.Router("/profile/changepwd", &controllers.MainController{}, "get,post:NewPassword")
    beego.Router("/tags/find", &controllers.MainController{}, "get,post:SearchTags")

    // AngularJS directives routes
    beego.Router("/ng/profile/info", &controllers.MainController{}, "get:ProfileInfo")

    // JSON data
    beego.Router("/json/events/accepted", &controllers.MainController{}, "get:AcceptedEvents")
}
