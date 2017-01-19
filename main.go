package main

import (
	_ "bee/activist/routers"
	//"bee/activist/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/plugins/cors"
	//"github.com/auth0/go-jwt-middleware"
	"bee/activist/controllers"
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
	orm.DefaultTimeLoc = time.UTC
	now := time.Now()
	log.Println(now)
	log.Println(now.Location())
	log.Println(now.Zone())

}

func main() {
	/*jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      secret := os.Getenv("AUTH0_CLIENT_SECRET")
      if secret == "" {
        return nil, errors.New("AUTH0_CLIENT_SECRET is not set")
      }
      return secret, nil
    },
  })*/

	beego.ErrorController(&controllers.ErrorController{})

	beego.SetStaticPath("/img","static/images")
	beego.SetStaticPath("/css","static/styles")
	beego.SetStaticPath("/js","static/javascript")
	beego.SetStaticPath("/fonts","static/fonts")
	beego.SetStaticPath("/tpl","static/templates")

	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "DELETE", "POST", "PUT", "PATCH", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Access-Control-Allow-Origin", "Content-Type", "Authorization"},
        ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin"},
        AllowCredentials: true,
    }))

	beego.Run()
}
