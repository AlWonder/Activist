package controllers

import (
  "github.com/astaxie/beego"
	"fmt"
  "bee/activist/models"
  "encoding/json"
	/*"strings"
	"net/http"
	"io/ioutil"*/
)

type SignupController struct {
  beego.Controller
}

func (c *SignupController) SignUp() {
  fmt.Println("Heya")
  /*
  url := "https://alwonder.eu.auth0.com/dbconnections/signup"

	payload := strings.NewReader("{\"client_id\": \"OJpLcs68G0OymE5F6VDoWHgFfWJJ8V3C\",\"email\": \"$('#signup-email').val()\",\"password\": \"$('#signup-password').val()\",\"user_metadata\": {\"name\": \"john\",\"color\": \"red\"}}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
  */

  u := models.User{}
  var dat map[string]interface{}
  json.Unmarshal(c.Ctx.Input.RequestBody, &dat)
  fmt.Println(dat)
    json.Unmarshal(c.Ctx.Input.RequestBody, &u)
  fmt.Println(u)
  c.Data["json"] = map[string]string{}
  c.ServeJSON()
}
