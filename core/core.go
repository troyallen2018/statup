package core

import (
	"github.com/GeertJohan/go.rice"
	"github.com/hunterlong/statup/plugin"
	"github.com/hunterlong/statup/types"
	"time"
)

type PluginJSON types.PluginJSON
type PluginRepos types.PluginRepos

type Core struct {
	Name           string     `db:"name" json:"name"`
	Description    string     `db:"description" json:"name"`
	Config         string     `db:"config" json:"-"`
	ApiKey         string     `db:"api_key" json:"-"`
	ApiSecret      string     `db:"api_secret" json:"-"`
	Style          string     `db:"style" json:"-"`
	Footer         string     `db:"footer" json:"-"`
	Domain         string     `db:"domain" json:"domain,omitempty"`
	Version        string     `db:"version" json:"version,omitempty"`
	Services       []*Service `json:"services,omitempty"`
	Plugins        []plugin.Info
	Repos          []PluginJSON
	AllPlugins     []plugin.PluginActions
	Communications []*types.Communication
	OfflineAssets  bool
	DbConnection   string
	started        time.Time
}

var (
	Configs     *Config
	CoreApp     *Core
	SqlBox      *rice.Box
	CssBox      *rice.Box
	ScssBox     *rice.Box
	JsBox       *rice.Box
	TmplBox     *rice.Box
	EmailBox    *rice.Box
	SetupMode   bool
	UsingAssets bool
)

func init() {
	CoreApp = NewCore()
}

func NewCore() *Core {
	CoreApp = new(Core)
	CoreApp.started = time.Now()
	return CoreApp
}

func InitApp() {
	SelectCore()
	SelectAllCommunications()
	InsertDefaultComms()
	LoadDefaultCommunications()
	SelectAllServices()
	CheckServices()
	go DatabaseMaintence()
}

func (c *Core) Update() (*Core, error) {
	res := DbSession.Collection("core").Find().Limit(1)
	err := res.Update(c)
	return c, err
}

func (c Core) UsingAssets() bool {
	return UsingAssets
}

func (c Core) SassVars() string {
	if !UsingAssets {
		return ""
	}
	return OpenAsset("scss/variables.scss")
}

func (c Core) BaseSASS() string {
	if !UsingAssets {
		return ""
	}
	return OpenAsset("scss/base.scss")
}

func (c Core) MobileSASS() string {
	if !UsingAssets {
		return ""
	}
	return OpenAsset("scss/mobile.scss")
}

func (c Core) AllOnline() bool {
	for _, s := range CoreApp.Services {
		if !s.Online {
			return false
		}
	}
	return true
}

func SelectCore() (*Core, error) {
	var c *Core
	err := DbSession.Collection("core").Find().One(&c)
	if err != nil {
		return nil, err
	}
	CoreApp = c
	CoreApp.DbConnection = Configs.Connection
	CoreApp.Version = VERSION
	CoreApp.Services, _ = SelectAllServices()
	//store = sessions.NewCookieStore([]byte(core.ApiSecret))
	return CoreApp, err
}
