package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	s "github.com/Ci-Jie/goproxy/storage"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// Version ...
func Version(c *fiber.Ctx) (err error) {
	pid, err := url.QueryUnescape(c.Locals(PackageKey).(string))
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusBadRequest)
		return
	}
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
	var fileName string = fmt.Sprintf("%s.info", version)
	type output struct {
		Version string `json:"Version"`
		Time    string `json:"Time"`
	}
	o := &output{
		Version: version,
	}
	var bOutput []byte
	exist, err := s.Use().Check(pid, version, fileName)
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	if exist {
		log.Infof("Get %s from backend stroage", fileName)
		bOutput, err = s.Use().Get(pid, version, fileName)
		if err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
	} else {
		log.Infof("%s doesn't exist in backend storage.", fileName)
		if _, ok := clients[c.Locals(DomainKey).(string)]; !ok {
			log.Infof("Downloading %s from %s", fileName, publicRepo)
			resp, _ := http.Get(fmt.Sprintf("%s%s", publicRepo, string(c.Request().URI().Path())))
			bOutput, _ = ioutil.ReadAll(resp.Body)
		} else {
			log.Infof("Downloading %s from %s", fileName, c.Locals(DomainKey).(string))
			commit, _, err := clients[c.Locals(DomainKey).(string)].Commits.GetCommit(pid, finalVersion)
			if err != nil {
				log.Warn(err)
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
			o.Time = commit.CommittedDate.Format(time.RFC3339)
			bOutput, err = json.Marshal(o)
			if err != nil {
				log.Error(err)
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
		}
		log.Infof("Saving %s into backend storage.", fileName)
		if err := s.Use().Create(pid, version, fileName, bOutput); err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
	}
	c.Send(bOutput)
	return
}
