package initalize

import (
	_ "embed"
	"fmt"
	"forge/biz/userservice"
	"forge/infra/cache"
	"forge/infra/configs"
	"forge/infra/coze"
	"forge/infra/database"
	"forge/infra/storage"
	"forge/interface/handler"
	"forge/pkg/log"
	"forge/pkg/loop"
	"forge/util"
)

func Init() {
	// load env
	path := initPath()
	introduce()
	log.InitLog(path, configs.Config())
	configs.MustInit(path)
	log.InitLog(path, configs.Config())
	database.MustInitDatabase(configs.Config())
	cache.MustInitCache(configs.Config())
	loop.MustInitLoop()
	coze.InitCozeService()
	storage.InitUserStorage()
	us := userservice.NewUserServiceImpl(storage.GetUserPersistence(), coze.GetCozeService())
	handler.MustInitHandler(us)

}
func initPath() string {
	return util.GetRootPath("")
}

// dont like? see https://patorjk.com/software/taag/#p=display&f=Merlin1&t=PLUTO
//
//go:embed .logo
var logo string

func introduce() {
	fmt.Println(logo)
}
