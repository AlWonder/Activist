package routers

import (
	"activist_api/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/*", &controllers.MainController{})

	// JSON data
	beego.Router("/tags", &controllers.MainController{}, "get:QueryTags")
	beego.Router("/tags/:tag", &controllers.MainController{}, "get:QueryEventsByTag")
	beego.Router("/events", &controllers.MainController{}, "get:QueryEvents;post:AddEvent;put:EditEvent")
	beego.Router("/events/:id([0-9]+)", &controllers.MainController{}, "get:GetEvent;delete:DeleteEvent")
	beego.Router("/login", &controllers.MainController{}, "post:Login")
	beego.Router("/signup", &controllers.MainController{}, "post:SignUp")
	beego.Router("/users", &controllers.MainController{}, "get:GetUserInfo")
	beego.Router("/users/:id([0-9]+)/events", &controllers.MainController{}, "get:QueryUserEvents")
}
