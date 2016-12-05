package controllers

import (
	"github.com/astaxie/beego/orm"
	"bee/activist/models"
	"log"
	"strings"
)

func (c *MainController) findTags(tag string) *[]models.Tag {
	var tags []models.Tag

	o := orm.NewOrm()
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	like := "%" + tag + "%"


	_, err := o.Raw(`SELECT name
					 FROM tags
					 WHERE name LIKE ?`, like).QueryRows(&tags)
	if err != nil {
		log.Println("findTags: ", err)
		return nil
	}
	log.Println(tags)
	return &tags
}