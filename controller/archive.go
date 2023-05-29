package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	s "github.com/Ci-Jie/goproxy/storage"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// Archive ...
func Archive(c *fiber.Ctx) (err error) {
	const fileName string = "source.zip"
	var format string = "zip"
	pid, err := url.QueryUnescape(c.Locals(PackageKey).(string))
	if err != nil {
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
	parts := strings.Split(version, "-")
	var finalVersion string
	switch len(parts) {
	case 1, 2:
		finalVersion = version
	case 3:
		finalVersion = parts[2]
	default:
		c.SendString("Invalid Version")
		c.Status(http.StatusInternalServerError)
		return
	}
	var buff *bytes.Buffer
	exist, err := s.Use().Check(pid, version, fileName)
	if err != nil {
		log.Error(err)
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	if exist {
		log.Infof("Get %s from backend storage", fileName)
		content, err := s.Use().Get(pid, version, fileName)
		if err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
		buff = bytes.NewBuffer(content)
	} else {
		if _, ok := clients[c.Locals(DomainKey).(string)]; !ok {
			log.Infof("Downloading %s from %s", fileName, publicRepo)
			resp, err := http.Get(fmt.Sprintf("%s%s", publicRepo, string(c.Request().URI().Path())))
			if err != nil {
				log.Error(err)
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
			content, _ := ioutil.ReadAll(resp.Body)
			buff = bytes.NewBuffer(content)
		} else {
			buff = bytes.NewBuffer([]byte{})
			writer := zip.NewWriter(buff)
			log.Infof("Downloading %s from %s", fileName, c.Locals(DomainKey).(string))
			gitResp, _, err := clients[c.Locals(DomainKey).(string)].Repositories.Archive(pid, &gitlab.ArchiveOptions{
				Format: &format,
				SHA:    &finalVersion,
			})
			if err != nil {
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
			reader, err := zip.NewReader(bytes.NewReader(gitResp), int64(len(gitResp)))
			if err != nil {
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
			for _, item := range reader.File {
				parts := strings.Split(item.Name, "/")
				if len(parts) == 0 {
					continue
				}
				directory := fmt.Sprintf("%s@%s", pid, version)
				var zfile io.Writer
				var err error
				if isV0orV1 {
					zfile, err = writer.Create(fmt.Sprintf("pegasus-cloud.com/%s", strings.Replace(item.Name, parts[0], directory, 1)))
				} else {
					path := fmt.Sprintf("pegasus-cloud.com/%s", strings.Replace(item.Name, parts[0], directory, 1))
					path = strings.ReplaceAll(path, "@", fmt.Sprintf("/%s@", v))
					zfile, err = writer.Create(path)
					fmt.Println("path:", path)
				}
				if err != nil {
					c.SendStatus(http.StatusInternalServerError)
					return err
				}
				closer, err := item.Open()
				defer closer.Close()
				if err != nil {
					c.SendStatus(http.StatusInternalServerError)
					return err
				}
				data, err := ioutil.ReadAll(closer)
				if err != nil {
					c.SendStatus(http.StatusInternalServerError)
					return err
				}
				if _, err := zfile.Write(data); err != nil {
					c.SendStatus(http.StatusInternalServerError)
					return err
				}
			}
			if err := writer.Close(); err != nil {
				c.SendStatus(http.StatusInternalServerError)
				return err
			}
		}
		if !isV0orV1 {
			version = fmt.Sprintf("%s/%s", v, version)
		}
		log.Infof("Saving %s into backend stroage.", fileName)
		if err := s.Use().Create(pid, version, fileName, buff.Bytes()); err != nil {
			log.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return err
		}
	}
	c.Response().Header.Set("Content-Type", "application/zip")
	c.Response().Header.Set("Content-Length", strconv.FormatInt(int64(buff.Len()), 10))
	c.SendString(buff.String())
	c.Status(http.StatusOK)
	return
}
