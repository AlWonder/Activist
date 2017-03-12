package main

import (
	_ "activist_api/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/plugins/cors"
	"time"
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
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", db_user + ":" + db_pw + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8&parseTime=True")
	orm.DefaultTimeLoc, _ = time.LoadLocation("Etc/GMT-9")
}

func main() {
	beego.SetStaticPath("/storage", "static/usrfiles")
	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "DELETE", "POST", "PUT", "PATCH", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Access-Control-Allow-Origin", "Content-Type", "enctype", "Authorization", "Content-Range", "Content-Disposition", "Content-Description", "Accept"},
        ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin"},
        AllowCredentials: true,
    }))

	beego.Run()
}
