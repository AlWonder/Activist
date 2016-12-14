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

func IsZero(t time.Time) (zero bool) {
	zero = t.IsZero()
	return
}

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
	beego.AddFuncMap("iszero",IsZero)

	beego.SetStaticPath("/images","static/images")
	beego.SetStaticPath("/css","static/css")
	beego.SetStaticPath("/js","static/js")
	beego.SetStaticPath("/fonts","static/fonts")


    beego.Run()
}

