package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// List ...
func List(c *fiber.Ctx) (err error) {
	pid, err := url.QueryUnescape(c.Locals(PackageKey).(string))
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusBadRequest)
		return
	}
	var output string
	if _, ok := clients[c.Locals(DomainKey).(string)]; !ok {
		resp, _ := http.Get(fmt.Sprintf("%s%s", publicRepo, string(c.Request().URI().Path())))
		bOutput, _ := ioutil.ReadAll(resp.Body)
		output = string(bOutput)
	} else {
		tags, _, err := clients[c.Locals(DomainKey).(string)].Tags.ListTags(pid, &gitlab.ListTagsOptions{})
		if err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
		result := make([]string, len(tags))
		for index, tag := range tags {
			result[index] = tag.Name
		}
		output = strings.Join(result, "\n")
	}
	c.SendString(output)
	return
}
