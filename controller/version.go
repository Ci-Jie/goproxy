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
	var isV0orV1 bool = true
	var v string
	splitedPID := strings.Split(pid, "/")
	if len(splitedPID) == 3 && strings.HasPrefix(splitedPID[2], "v") {
		pid = fmt.Sprintf("%s/%s", splitedPID[0], splitedPID[1])
		isV0orV1 = false
		v = splitedPID[2]
	}
	version := c.Locals(VersionKey).(string)
	vid := strings.TrimPrefix(strings.Split(strings.Split(version, "-")[0], ".")[0], "v")
	if isV0orV1 && vid != "0" && vid != "1" {
		c.SendStatus(http.StatusNotFound)
		return nil
	}
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
			commit, commitResp, err := clients[c.Locals(DomainKey).(string)].Commits.GetCommit(pid, finalVersion)
			// fmt.Println("commitResp.Request.URL:", commitResp.Request.URL)
			if err != nil {
				if commitResp.StatusCode == http.StatusNotFound {
					c.SendStatus(http.StatusNotFound)
					return nil
				} else {
					log.Warn(err)
					c.SendStatus(http.StatusInternalServerError)
				}
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
		if !isV0orV1 {
			version = fmt.Sprintf("%s/%s", v, version)
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
