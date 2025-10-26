package initalize

import (
	_ "embed"
	"fmt"
	"forge/infra/cache"
	"forge/infra/configs"
	"forge/infra/database"
	"forge/pkg/log"
	"forge/pkg/loop"
	"forge/util"
)

func Init() {
	path := initPath()
	introduce()
	configs.MustInit(path)
	log.InitLog(path, configs.Config())
	database.MustInitDatabase(configs.Config())
	cache.MustInitCache(configs.Config())
	loop.MustInitLoop()
}
func initPath() string {
	return util.GetRootPath("")
}

// dont like? see https://patorjk.com/software/taag/#p=display&f=Merlin1&t=PLUTO
//
//go:embed logo.txt
var logo string

func introduce() {
	fmt.Println(logo)
}
