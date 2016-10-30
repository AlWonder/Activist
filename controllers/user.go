package controllers

import (
	"log"
	pk "bee/activist/utilities/pbkdf2"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"bee/activist/models"
	"strconv"
	"encoding/hex"
	"strings"
	"time"
)

func (c *MainController) Login() {
	back := strings.Replace(c.Ctx.Input.Param(":back"), ">", "/", -1)
	log.Println("back is", back)
	c.activeContent("login")
	if c.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		var x pk.PasswordHash
		email := c.Input().Get("email")
		password := c.Input().Get("password")

		valid := validation.Validation{}
		valid.Required(email, "email")
		valid.Required(password, "password")
		valid.Email(email, "email")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			c.Data["Errors"] = errormap
			return
		}
		//log.Println("Authorization is", email, ":", password)

		user := c.getUser(email)
		if user != nil {

			/*log.Println("id: ", user.Id)
			log.Println("login: ", user.Email)
			log.Println("passwd: ", user.Password)
			log.Println("usr_group: ", user.UserGroup)
			log.Println("Password to scan:", user.Password)*/
			
			x.Hash = make([]byte, 32)
			x.Salt = make([]byte, 16)
			var err error

			if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
				log.Println("ERROR:", err)
			}
			if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
				log.Println("ERROR:", err)
			}
			//log.Println("decoded password is", x)
		} else {
			flash.Error("No such user/email")
			flash.Store(&c.Controller)
			return
		}

		if !pk.MatchPassword(password, &x) {
			flash.Error("Wrong password")
			flash.Store(&c.Controller)
			return
		}

		m := make(map[string]interface{})
		m["id"] = user.Id
		m["email"] = user.Email
		m["group"] = user.UserGroup
		m["timestamp"] = time.Now()
		m["first_name"] = user.FirstName
		m["second_name"] = user.SecondName
		m["last_name"] = user.LastName
		m["gender"] = user.Gender 
		c.SetSession("activist", m)
		c.Redirect("/"+back, 302)
		
		flash.Notice("Welcome, " + c.Input().Get("email"))
		c.Redirect("/"+back, 302)
	}
}

func (this *MainController) Logout() {
	this.DelSession("activist")
	this.Redirect("/home", 302)
}

func (c *MainController) getUser(email string) *models.User {
	o := orm.NewOrm()
	user := models.User{Email: email}
	err := o.Read(&user, "email")

	//err := o.Raw("SELECT * FROM users WHERE login = ?", email).QueryRow(&user)

	if err == orm.ErrNoRows {
    	log.Println("No result found.")
    	return nil
	} else if err == orm.ErrMissPK {
	    log.Println("No primary key found.")
	    return nil
	}
	return &user
}

func (c *MainController) Register() {
	c.activeContent("register")
	if c.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		email := c.Input().Get("email")
		password := c.Input().Get("password")
		password2 := c.Input().Get("password2")
		group, err := strconv.ParseInt(c.Input().Get("group"), 10 , 64)
		if err != nil {
			flash.Error("Wrong group")
			flash.Store(&c.Controller)
			return
		}
		firstName := c.Input().Get("first_name")
		secondName := c.Input().Get("second_name")
		lastName := c.Input().Get("last_name")
		gender, err := strconv.ParseInt(c.Input().Get("gender"), 10 , 64)
		if err != nil {
			flash.Error("Wrong gender")
			flash.Store(&c.Controller)
			return
		}

		valid := validation.Validation{}
		valid.Email(email, "email")
		valid.Required(email, "email")
		valid.Required(password, "password")
		valid.Required(password2, "password2")
		valid.Required(group, "group")
		valid.Required(firstName, "first_name")
		valid.Required(secondName, "second_name")
		valid.Required(lastName, "last_name")
		valid.Required(gender, "gender")
		valid.MaxSize(email, 30, "email")
		valid.MaxSize(firstName, 25, "first_name")
		valid.MaxSize(secondName, 25, "second_name")
		valid.MaxSize(lastName, 25, "last_name")

		if valid.HasErrors() {
			errormap := []string{}
			log.Println("There are some valid errors")
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			c.Data["Errors"] = errormap
			return
		}

		if password != password2 {
			flash.Error("Passwords don't match")
			flash.Store(&c.Controller)
			return
		}

		h := pk.HashPassword(password)

		o := orm.NewOrm()
	    

	    user := models.User{Email: email, UserGroup: group, FirstName: firstName, SecondName: secondName, LastName: lastName, Gender: gender}

		// Convert password hash to string
		user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)
		log.Println("hex: " + hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt))

	    log.Println(o.Insert(&user))
	    c.Redirect("/home", 302)
	}
}
