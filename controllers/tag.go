package controllers

import (
	"github.com/astaxie/beego/orm"
	"activist_api/models"
	"log"
	"strings"
)

func(c *MainController) QueryTags() {
	tag := c.Input().Get("query")
	tags := c.getTags(tag)
	c.Data["json"] = &tags
	c.ServeJSON()
}

func (c *MainController) getTags(tag string) *[]string {
	var tags []string

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
	return &tags
}

func (c *MainController) addTags(tags []string) []int64 {
	var tagIds []int64
	o := orm.NewOrm()
	for _, tag := range tags {
		newTag := models.Tag{
			Name: tag,
		}

		if created, id, err := o.ReadOrCreate(&newTag, "Name"); err == nil {
	    if created {
        log.Println("New Insert an object. Id:", id)
	    } else {
        log.Println("Get an object. Id:", id)
	    }
			tagIds = append(tagIds, id)
		} else {
			log.Println(err)
		}
	}

	return tagIds
}

func (c *MainController) deleteEventTags(eventId int64, tags []string) error {
	o := orm.NewOrm()
	for _, tag := range tags {
		log.Println("Deleting " + tag)
		if _, err := o.Delete(&models.Tag{Name: tag}); err != nil {
			return err
		}
	}
	return nil
}

func (c *MainController) addEventTags(eventId int64, tagIds []int64) bool {
	var ok bool
	o := orm.NewOrm()
	for _, tagId := range tagIds {
		tagEvent := models.EventTag {
			EventId: eventId,
			TagId: tagId,
		}

		if _, err := o.Insert(&tagEvent); err != nil {
			log.Println(err)
			ok = false
		}
	}
	ok = true
	return ok
}

func (c *MainController) getTagsByEventId(id int64) *[]string {
	var tags []string
	o := orm.NewOrm()
	if _, err := o.Raw(`SELECT name
		FROM tags
		INNER JOIN events_tags ON tags.id = events_tags.tag_id
		WHERE event_id = ?`, id).QueryRows(&tags); err != nil {
		log.Println(err)
		return nil
	}
	return &tags
}
