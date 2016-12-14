package main

import (
	_ "bee/activist/routers"
	//"bee/activist/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"bee/activist/controllers"
	"time"
)

// Database settings
const (
    db_name = "activist"
    db_user = "root"
    db_pw = "1111"
    force = false
    verbose = false
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", db_user + ":" + db_pw + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8&parseTime=True")
	orm.DefaultTimeLoc = time.UTC
	now := time.Now()
	log.Println(now)
	log.Println(now.Location())
	log.Println(now.Zone())
}

func main() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.SetStaticPath("/img","static/images")
	beego.SetStaticPath("/css","static/styles")
	beego.SetStaticPath("/js","static/javascript")
	beego.SetStaticPath("/fonts","static/fonts")
	beego.SetStaticPath("/tpl","static/templates")

	beego.Run()
}
