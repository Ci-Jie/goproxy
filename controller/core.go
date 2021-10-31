package controller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var clients = make(map[string]*gitlab.Client)

const publicRepo string = "https://proxy.golang.org"

type goproxy struct {
	Gitlabs []gitlabs `mapstructure:"gitlabs"`
	Storage storage   `mapstructure:"storage"`
}

type gitlabs struct {
	Domain   string `mapstructure:"domain"`
	Endpoint string `mapstructur:"endpoint"`
	Token    string `mapstructure:"token"`
}

type storage struct {
	Provider string `mapstructure:"provider"`
}

const (
	// DomainKey ...
	DomainKey string = "domain"
	// GroupKey ...
	GroupKey string = "group"
	// SubGroupKey ...
	SubGroupKey string = "subgroup"
	// ProjectKey ...
	ProjectKey string = "project"
	// PackageKey ...
	PackageKey string = "package"
	// VersionKey ...
	VersionKey string = "version"
)

// Init ...
func Init() {
	goproxy := &goproxy{}
	if err := viper.Unmarshal(goproxy); err != nil {
		log.Error(err)
	}
	for _, item := range goproxy.Gitlabs {
		var err error
		clients[item.Domain], err = gitlab.NewClient(item.Token, gitlab.WithBaseURL(item.Endpoint))
		if err != nil {
			panic(err)
		}
	}
}

// Handler ...
func Handler(c *fiber.Ctx) (err error) {
	path := strings.Split(strings.ReplaceAll(string(c.Context().URI().Path()), "/@v", ""), "/")
	c.Locals(DomainKey, path[1])
	c.Locals(PackageKey, strings.Join(path[2:len(path)-1], "/"))
	c.Locals(VersionKey, path[len(path)-1])
	switch {
	case strings.EqualFold(c.Locals(VersionKey).(string), "list"):
		return List(c)
	case strings.Contains(c.Locals(VersionKey).(string), ".info"):
		c.Locals(VersionKey, strings.ReplaceAll(c.Locals(VersionKey).(string), ".info", ""))
		return Version(c)
	case strings.Contains(c.Locals(VersionKey).(string), ".mod"):
		c.Locals(VersionKey, strings.ReplaceAll(c.Locals(VersionKey).(string), ".mod", ""))
		return Mod(c)
	case strings.Contains(c.Locals(VersionKey).(string), ".zip"):
		c.Locals(VersionKey, strings.ReplaceAll(c.Locals(VersionKey).(string), ".zip", ""))
		return Archive(c)
	}
	return nil
}
