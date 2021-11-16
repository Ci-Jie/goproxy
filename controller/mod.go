package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	s "github.com/Ci-Jie/goproxy/storage"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

const fileName = "go.mod"

// Mod ...
func Mod(c *fiber.Ctx) (err error) {
	const fileName string = "go.mod"
	pid, err := url.QueryUnescape(c.Locals(PackageKey).(string))
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusBadRequest)
		return
	}
	// version := c.Params("version")
	version := c.Locals(VersionKey).(string)
	parts := strings.Split(version, "-")
	var finalVersion string
	switch len(parts) {
	case 1, 2:
		finalVersion = version
	case 3:
		finalVersion = parts[2]
	default:
		c.SendString("Version is invalid")
		c.Status(http.StatusInternalServerError)
		return
	}
	var content []byte
	exist, err := s.Use().Check(pid, version, fileName)
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	if exist {
		log.Infof("Get %s from backend storage", fileName)
		content, err = s.Use().Get(pid, version, fileName)
		if err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
	} else {
		if _, ok := clients[c.Locals(DomainKey).(string)]; !ok {
			log.Infof("Downloading %s from %s", fileName, publicRepo)
			resp, _ := http.Get(fmt.Sprintf("%s%s", publicRepo, string(c.Request().URI().Path())))
			content, _ = ioutil.ReadAll(resp.Body)
		} else {
			log.Infof("Downloading %s from %s", fileName, c.Locals(DomainKey).(string))
			content, _, err = clients[c.Locals(DomainKey).(string)].RepositoryFiles.GetRawFile(pid, fileName, &gitlab.GetRawFileOptions{
				Ref: &finalVersion,
			})
			if err != nil {
				log.Error(err)
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
		}
		log.Infof("Saving %s into backend stroage.", fileName)
		if err := s.Use().Create(pid, version, fileName, content); err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
	}
	c.SendString(string(content))
	c.Status(http.StatusOK)
	return nil
}
