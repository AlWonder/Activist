package main

import (
	_ "activist_api/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "gopkg.in/gographics/imagick.v2/imagick"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/plugins/cors"
	"time"
	"log"
)
	// Database settings
const (
    db_name = "activist"
    db_user = "root"
    db_pw = "fluttershy"
    force = false
    verbose = false
)

func init() {
	log.Println("Main loaded")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", db_user + ":" + db_pw + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8&parseTime=True")
	orm.DefaultTimeLoc, _ = time.LoadLocation("Etc/GMT-9")
}

func main() {
	//imagick.Initialize()
	//defer imagick.Terminate()
	log.Println("l")
	beego.SetStaticPath("/api/storage", "static/usrfiles")
	beego.SetStaticPath("/api/storage/cover", "static/usrfiles/event/")
	beego.SetStaticPath("/api/storage/cover/sm", "static/usrfiles/event/")
	beego.SetStaticPath("/api/storage/avatar", "static/usrfiles/user/avatar/")
	beego.SetStaticPath("/api/storage/avatar/sm", "static/usrfiles/user/avatar/")
	beego.SetStaticPath("/api/storage/docs/tpl", "static/usrfiles/user/tpls/")
	beego.SetStaticPath("/api/storage/docs/form", "static/usrfiles/user/forms/")

	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "DELETE", "POST", "PUT", "PATCH", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Access-Control-Allow-Origin", "Content-Type", "enctype", "Authorization", "Content-Range", "Content-Disposition", "Content-Description", "Accept"},
        ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin"},
        AllowCredentials: true,
    }))

	beego.Run()
}
