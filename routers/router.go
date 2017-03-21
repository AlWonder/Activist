package routers

import (
	"activist_api/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/*", &controllers.MainController{})

	// JSON data
	beego.Router("/index", &controllers.MainController{}, "get:IndexPage")
	beego.Router("/tags", &controllers.MainController{}, "get:QueryTags")
	beego.Router("/tags/:tag", &controllers.MainController{}, "get:QueryEventsByTag")
	beego.Router("/tags/:tag/status", &controllers.MainController{}, "get:GetTagStatus;post:AddTagStatus;delete:DeleteTagStatus")
	beego.Router("/events", &controllers.MainController{}, "get:QueryEvents;post:AddEvent;put:EditEvent")
	beego.Router("/events/:id([0-9]+)", &controllers.MainController{}, "get:GetEvent;delete:DeleteEvent")
	beego.Router("/events/:id([0-9]+)/join", &controllers.MainController{}, "post:JoinEvent;delete:DenyEvent")
	beego.Router("/events/:id([0-9]+)/joined", &controllers.MainController{}, "get:GetJoinedUsers")
	beego.Router("/events/:id([0-9]+)/cover", &controllers.MainController{}, "post:AddCover")
	beego.Router("/login", &controllers.MainController{}, "post:Login")
	beego.Router("/signup", &controllers.MainController{}, "post:SignUp")
	beego.Router("/users", &controllers.MainController{}, "get:GetUserInfo")
	beego.Router("/users/avatar", &controllers.MainController{}, "post:AddAvatar")
	beego.Router("/users/:id([0-9]+)/events", &controllers.MainController{}, "get:QueryUserEvents")
	beego.Router("/users/:id([0-9]+)/joined", &controllers.MainController{}, "get:QueryJoinedEvents")
	//beego.Router("/upload", &controllers.MainController{}, "post:UploadFile")
}
