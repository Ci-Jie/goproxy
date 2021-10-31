package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xanzy/go-gitlab"
)

// Archive2 ...
func Archive2(c *fiber.Ctx) (err error) {
	var format string = "zip"
	p, err := url.QueryUnescape(fmt.Sprintf("%s/%s", c.Params("group"), c.Params("project")))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return
	}
	finalVersion := c.Params("version")
	parts := strings.Split(finalVersion, "-")
	var version string
	switch len(parts) {
	case 1:
		version = finalVersion
	case 3:
		version = parts[2]
	default:
		c.SendString("Invalid Version")
		c.Status(http.StatusInternalServerError)
		return
	}
	content, _, err := clients[c.Params("domain")].Repositories.Archive(p, &gitlab.ArchiveOptions{
		Format: &format,
		SHA:    &version,
	})
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return
	}
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return
	}
	buffer := bytes.NewBuffer([]byte{})
	writer := zip.NewWriter(buffer)
	directory := filepath.Join("tmp", p, finalVersion)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0775)
	}
	file, err := os.Create(filepath.Join(directory, "source.zip"))
	defer file.Close()
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	for _, item := range reader.File {
		parts := strings.Split(item.Name, "/")
		if len(parts) == 0 {
			continue
		}
		directory := fmt.Sprintf("%s@%s", p, finalVersion)
		zfile, err := writer.Create(fmt.Sprintf("pegasus-cloud.com/%s", strings.Replace(item.Name, parts[0], directory, 1)))
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
	if _, err := file.Write(buffer.Bytes()); err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	if err := writer.Close(); err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	c.Response().Header.Set("Content-Type", "application/zip")
	c.Response().Header.Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	c.SendString(buffer.String())
	c.Status(http.StatusOK)
	return
}
