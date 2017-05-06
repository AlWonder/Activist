package routers

import (
	"activist_api/controllers"
	"github.com/astaxie/beego"
	"log"
)

func init() {
	beego.Router("/*", &controllers.MainController{})

	// JSON data
	beego.Router("/api/index", &controllers.MainController{}, "get:IndexPage")
	beego.Router("/api/tags", &controllers.TagController{}, "get:QueryTags")
	beego.Router("/api/tags/:tag", &controllers.EventController{}, "get:QueryEventsByTag")
	beego.Router("/api/events", &controllers.EventController{}, "get:QueryEvents;post:AddEvent;put:EditEvent")
	beego.Router("/api/events/:id([0-9]+)", &controllers.EventController{}, "get:GetEvent;delete:DeleteEvent")
	beego.Router("/api/events/joined", &controllers.EventController{}, "get:QueryJoinedEvents")
	beego.Router("/api/events/:id([0-9]+)/join", &controllers.EventController{}, "post:JoinEvent;delete:DenyEvent")
	beego.Router("/api/events/:id([0-9]+)/joined", &controllers.UserController{}, "get:GetJoinedUsers")
	beego.Router("/api/events/:id([0-9]+)/cover", &controllers.FileController{}, "post:AddCover;put:EditCover")
	beego.Router("/api/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/api/signup", &controllers.AuthController{}, "post:SignUp")
	beego.Router("/api/users", &controllers.UserController{}, "get:GetUserInfo")
	beego.Router("/api/users/avatar", &controllers.FileController{}, "post:AddAvatar")
	beego.Router("/api/users/:id([0-9]+)/events", &controllers.EventController{}, "get:QueryUserEvents")
	beego.Router("/api/tpl", &controllers.FileController{}, "post:AddFormTemplate")
	beego.Router("/api/tpl/:id([0-9]+)", &controllers.FormController{}, "get:QueryUserFormTemplates")
	beego.Router("/api/form", &controllers.FormController{}, "get:QueryUserForms")
	beego.Router("/api/form/:id([0-9]+)", &controllers.FileController{}, "post:AddVolunteerForm;delete:DeleteVolunteerForm")
	beego.Router("/api/xaccel/tpl/:id([0-9]+)/:name", &controllers.MainController{}, "get:XAccelTemplate")
	beego.Router("/api/xaccel/form/:id([0-9]+)/:name", &controllers.MainController{}, "get:XAccelForm")
	beego.Router("/api/xaccel/generate/tpl/:id", &controllers.MainController{}, "get:GenerateTemplateToken")
	beego.Router("/api/xaccel/generate/form/:id", &controllers.MainController{}, "get:GenerateFormToken")
	//beego.Router("/upload", &controllers.MainController{}, "post:UploadFile")

	log.Println("Routers loaded")
}
