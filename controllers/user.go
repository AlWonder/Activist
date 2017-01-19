package controllers

import (
	"log"
	pk "bee/activist/utilities/pbkdf2"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"bee/activist/models"
	"strconv"
	"errors"
	"encoding/hex"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"time"

)

var privateKey = []byte("pisos")

type Error struct {
	UserMessage    string      `json:"userMessage"`
	Code           float64     `json:"code"`
}

type LoginResponse struct {
	IdToken        string      `json:"idToken"`
	Errors         []Error     `json:"errors"`
}

type GetUserInfoResponse struct {
	User           UserInfo `json:"user"`
	Errors         []Error     `json:"errors"`
}

type UserInfo struct {
	Email          *string	`json:"email"`
	Group          *int64	  `json:"group"`
	FirstName      *string	`json:"firstName"`
	SecondName     *string	`json:"secondName"`
	LastName       *string	`json:"lastName"`
	Gender         *int64	  `json:"gender"`
}

func (c *MainController) Login() {
	request := make(map[string]interface{})
	var response LoginResponse

	// Checking for a correct JSON request. If not, throw an error to a client
  if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.appendLoginError(&response, "Request error", 400)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Checking if username and password fields are correct
	var email, password string
	if field, ok := request["username"].(string); ok {
		email = field
	} else {
		c.appendLoginError(&response, "Bad username field", 400)
	}
	if field, ok := request["password"].(string); ok {
		password = field
	} else {
		c.appendLoginError(&response, "Bad password field", 400)
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(email, "email")
	valid.Required(password, "password")
	valid.Email(email, "email")
	valid.MaxSize(email, 30, "email")
	valid.MaxSize(password, 30, "email")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendLoginError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having validation errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Getting a user from db
	user := c.getUser(email)
	if user == nil {
		c.appendLoginError(&response, "Пользователь с таким email не найден", 400)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Comparing passwords from a client and the database
	var x pk.PasswordHash
	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)
	var err error

	if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
		log.Println("ERROR:", err)
	}
	if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
		log.Println("ERROR:", err)
	}

	if !pk.MatchPassword(password, &x) {
		c.appendLoginError(&response, "Неверный email/пароль", 400)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Generating a token and sending it to a client
	token := c.generateToken(user.Email)
	response.IdToken = token

	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) SignUp() {
	user := make(map[string]interface{})
  json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	var response LoginResponse

	if userInterface, ok := user["user"].(map[string]interface{}); !ok {
		c.appendLoginError(&response, "Request error", 400)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = response
		c.ServeJSON()
		return
	} else {
		user = userInterface
	}

	// Validate input fields
	valid := validation.Validation{}
	valid.Email(user["email"], "email")
	valid.Required(user["email"], "email")
	valid.Required(user["password"], "password")
	valid.Required(user["group"], "group")
	valid.Required(user["firstName"], "first_name")
	valid.Required(user["secondName"], "second_name")
	valid.Required(user["lastName"], "last_name")
	valid.Required(user["gender"], "gender")
	valid.MaxSize(user["email"], 30, "email")
	valid.MaxSize(user["password"], 30, "email")
	valid.MaxSize(user["firstName"], 25, "first_name")
	valid.MaxSize(user["secondName"], 25, "second_name")
	valid.MaxSize(user["lastName"], 25, "last_name")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendLoginError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	var newUser models.User

	// Checking if fields have correct type
	c.checkStringField(&newUser.Email, user["email"], &response, "Email")
	c.checkStringField(&newUser.Password, user["password"], &response, "Password")
	c.checkStringField(&newUser.FirstName, user["firstName"], &response, "First name")
	c.checkStringField(&newUser.SecondName, user["secondName"], &response, "Second name")
	c.checkStringField(&newUser.LastName, user["lastName"], &response, "Last name")
	c.checkIntField(&newUser.Gender, user["gender"], &response, "Gender")
	c.checkIntField(&newUser.Group, user["group"], &response, "Group")

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// If it's alright, hashing the password and creating a new user
	h := pk.HashPassword(user["password"].(string))
	o := orm.NewOrm()

	newUser.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)
	log.Println("hex: " + hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt))
	if _, err := o.Insert(&newUser); err != nil {
		c.appendLoginError(&response, "Не удалось зарегистрироваться. Возможно, пользователь с таким именем уже существует", 400)
		log.Println(err)
	} else {
		// Generating a token and sending it to a client
		token := c.generateToken(newUser.Email)
		response.IdToken = token
	}

	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) GetUserInfo() {
	var response GetUserInfoResponse
	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.appendGetUserInfoError(&response, "Invalid token. Access denied", 401)
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = response
		c.ServeJSON()
		return
	} else {
		log.Println("Heya, " + payload["sub"].(string))
		user := c.getUser(payload["sub"].(string))
		log.Println(user)
		response.User = UserInfo {
                  Email:       &user.Email,
                  Group:       &user.Group,
                  FirstName:   &user.FirstName,
                  SecondName:  &user.SecondName,
                  LastName:    &user.LastName,
                  Gender:      &user.Gender,
		}
		c.Data["json"] = response
		c.ServeJSON()
	}
}

func (c *MainController) getUser(email string) *models.User {
	o := orm.NewOrm()
	user := models.User{Email: email}
	err := o.Read(&user, "email")

	if err == orm.ErrNoRows {
    	log.Println("No result found.")
    	return nil
	} else if err == orm.ErrMissPK {
	    log.Println("No primary key found.")
	    return nil
	}
	return &user
}

func (c *MainController) generateToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "http://localhost:8080",
    "sub": username,
		"iat": time.Now().Unix(),
    "exp": time.Now().Unix() + 36000,
	})

	tokenString, _ := token.SignedString(privateKey)

	return tokenString
}

func (c *MainController) validateToken() (jwt.MapClaims, error){
	tokenString := c.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		log.Println("Token not found")
		return nil, errors.New("Couldn't find Authorization header")
	}

	token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
	    return privateKey, nil
	})

	if token.Valid {
			log.Println("Heya!")
	    return token.Claims.(jwt.MapClaims), nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
	    if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	        log.Println("That's not even a token")
					return nil, err
	    } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	        // Token is either expired or not active yet
	        log.Println("Timing is everything")
					return nil, err
	    } else {
	        log.Println("Couldn't handle this token:", err)
					return nil, err
	    }
	} else {
	    log.Println("Couldn't handle this token:", err)
			return nil, err
	}
}

// Checks that json field is string and appends an error into response if it isn't
func (c *MainController) checkStringField(userProperty *string, field interface{}, response *LoginResponse, fieldName string) {
	if checkedField, ok := field.(string); ok {
		*userProperty = checkedField
		log.Println(checkedField)
	} else {
		c.appendLoginError(response, "Datatype error in " + fieldName, 400)
	}
}

func (c *MainController) checkIntField(userProperty *int64, field interface{}, response *LoginResponse, fieldName string) {
	if stringField, ok := field.(string); ok {
		if checkedField, err := strconv.ParseInt(stringField, 10, 64); err != nil {
			c.appendLoginError(response, "Datatype error in " + fieldName, 400)
		} else {
			*userProperty = checkedField
		}
	} else {
		c.appendLoginError(response, "Datatype error in " + fieldName, 400)
	}
}

// Appends an error into the response body
func (c *MainController) appendLoginError(response *LoginResponse, message string, code float64) {
	response.Errors = append(response.Errors, Error {
		UserMessage: message,
		Code: code,
	})
}

// The same as previous but for get user response
// Where are my generics Google?!
func (c *MainController) appendGetUserInfoError(response *GetUserInfoResponse, message string, code float64) {
	response.Errors = append(response.Errors, Error {
		UserMessage: message,
		Code: code,
	})
}

func (c *MainController) NewPassword() {
	m := c.getSessionInfo()
    if m == nil {
        c.Redirect("/home", 302)
    }

	flash := beego.NewFlash()
	c.activeContent("changepwd", "Изменить", []string{}, []string{})
	if c.Ctx.Input.Method() != "POST" {
		return
	}

	userId := m["id"].(int64)
	oldPassword := c.Input().Get("old_password")
	newPassword := c.Input().Get("new_password")
	if ok := c.changePassword(userId, oldPassword, newPassword); ok == true {
		c.Redirect("/profile", 302)
	} else {
		log.Println("Nope")
		flash.Error("Неверный старый пароль.")
		flash.Store(&c.Controller)
	}
}


func (c *MainController) changePassword(userId int64, oldPassword, newPassword string) bool {
	h := pk.HashPassword(oldPassword)

	o := orm.NewOrm()

	user := models.User{Id: userId}

	err := o.Read(&user);
	if err == orm.ErrNoRows {
    	log.Println("No result found.")
    	return false
	} else if err == orm.ErrMissPK {
	    log.Println("No primary key found.")
	    return false
	}

	var x pk.PasswordHash
	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
		log.Println("ERROR:", err)
	}
	if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
		log.Println("ERROR:", err)
	}

	if !pk.MatchPassword(oldPassword, &x) {
		log.Println("Passwords don't match")
		return false
	}

	hn := pk.HashPassword(newPassword)

	user.Password = hex.EncodeToString(hn.Hash) + hex.EncodeToString(hn.Salt)
	log.Println("hex: " + hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt))

	if _, err = o.Update(&user); err != nil {
        log.Println("changePassword, data update: ", err)
		return false
    }
    return true
}
