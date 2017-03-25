package controllers

import (
	"activist_api/models"
	pk "activist_api/utilities/pbkdf2"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var privateKey = []byte("B7tfu34bfkderf43bfkj4bfkjerf")

func (c *MainController) Login() {
	var request models.LoginRequest
	var response models.LoginResponse

	// Checking for a correct JSON request. If not, throw an error to a client
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	valid := validation.Validation{}
	valid.Required(request.Email, "email")
	valid.Required(request.Password, "password")
	valid.Email(request.Email, "email")
	valid.MaxSize(request.Email, 30, "email")
	valid.MaxSize(request.Password, 30, "email")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			//c.appendLoginError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
		c.sendError("Неверно введённые данные", 1)
		return
	}

	// Getting a user from db
	user := models.GetUserByEmail(request.Email)
	if user == nil {
		c.sendError("Неверный email / пароль", 14)
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

	if !pk.MatchPassword(request.Password, &x) {
		c.sendError("Неверный email / пароль", 14)
		return
	}

	// Generating a token and sending it to a client
	token := c.generateToken(user.Id)
	response.Ok = true
	response.IdToken = token

	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) SignUp() {
	var request models.SignUpRequest
	var response models.LoginResponse

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	log.Println(request)

	// Validate input fields
	valid := validation.Validation{}
	valid.Email(request.User.Email, "email")
	valid.Required(request.User.Email, "email")
	valid.Required(request.User.Password, "password")
	valid.Required(request.User.Group, "group")
	valid.Required(request.User.FirstName, "first_name")
	valid.Required(request.User.SecondName, "second_name")
	valid.Required(request.User.LastName, "last_name")
	valid.Required(request.User.Gender, "gender")
	valid.MaxSize(request.User.Email, 30, "email")
	valid.MaxSize(request.User.Password, 30, "email")
	valid.MaxSize(request.User.FirstName, 25, "first_name")
	valid.MaxSize(request.User.SecondName, 25, "second_name")
	valid.MaxSize(request.User.LastName, 25, "last_name")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			models.AppendError(&response.Errors, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// If it's alright, hashing the password and creating a new user
	h := pk.HashPassword(request.User.Password)
	o := orm.NewOrm()

	request.User.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)

	// Uploading an avatar
	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
    if file != nil {
        // get the filename
        fileName := header.Filename
        // save to server
        err := c.SaveToFile("file", "/static/user/" + fileName)
				log.Println(err)
    }

	if _, err := o.Insert(&request.User); err != nil {
		c.sendError("Couldn't sign up. The user probably already exists", 14)
		log.Println(err)
	} else {
		// Generating a token and sending it to a client
		token := c.generateToken(request.User.Id)
		response.IdToken = token
	}

	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) generateToken(userId int64) string {
	// Filling the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "http://localhost:8080",
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 36000,
	})

	tokenString, _ := token.SignedString(privateKey)

	return tokenString
}

// Checks authorization header
// If it's valid, returns a payload
func (c *MainController) validateToken() (jwt.MapClaims, error) {
	// Get a token from request header
	tokenString := c.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		log.Println("Token not found")
		return nil, errors.New("Couldn't find Authorization header")
	}

	token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if token == nil {
		return nil, errors.New("Token is null")
	}

	if token.Valid {
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
