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
	splitedPID := strings.Split(pid, "/")
	if len(splitedPID) == 3 && strings.HasPrefix(splitedPID[2], "v") {
		pid = fmt.Sprintf("%s/%s", splitedPID[0], splitedPID[1])
	}
	var output string
	if _, ok := clients[c.Locals(DomainKey).(string)]; !ok {
		resp, _ := http.Get(fmt.Sprintf("%s%s", publicRepo, string(c.Request().URI().Path())))
		bOutput, _ := ioutil.ReadAll(resp.Body)
		output = string(bOutput)
	} else {
		tags, tagsResp, err := clients[c.Locals(DomainKey).(string)].Tags.ListTags(pid, &gitlab.ListTagsOptions{})
		fmt.Println("tagsResp.Request.URL:", tagsResp.Request.URL)
		if err != nil {
			if tagsResp.StatusCode == http.StatusNotFound {
				c.SendStatus(http.StatusNotFound)
				return nil
			} else {
				c.SendStatus(http.StatusInternalServerError)
				log.Error(err)
			}
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
