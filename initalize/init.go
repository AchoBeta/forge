package initalize

import (
	_ "embed"
	"fmt"
	"forge/biz/cosservice"
	"forge/biz/mindmapservice"
	"forge/biz/userservice"
	"forge/infra/cache"
	"forge/infra/configs"
	"forge/infra/cos"
	"forge/infra/coze"
	"forge/infra/database"
	"forge/infra/storage"
	"forge/interface/handler"
	"forge/interface/router"
	"forge/pkg/log"

	// "forge/pkg/loop"
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
	// TODO: cozeloop配置好后启用
	// loop.MustInitLoop()
	coze.InitCozeService()
	storage.InitUserStorage()
	storage.InitMindMapStorage()

	// snowflake - 从配置文件读取节点ID
	snowflakeConfig := configs.Config().GetSnowflakeConfig()
	if err := util.InitSnowflake(snowflakeConfig.NodeID); err != nil {
		// 初始化失败，直接 panic 提示原因
		panic(fmt.Sprintf("init snowflake failed: %v", err))
	}

	// 从配置文件读取JWT配置并创建JWTUtil
	jwtConfig := configs.Config().GetJWTConfig()
	jwtUtil := util.NewJWTUtil(jwtConfig.SecretKey, jwtConfig.ExpireHours)

	// 依赖注入：创建COS服务实例
	cosConfig := configs.Config().GetCOSConfig()
	cosService := cos.NewCOSService(cosConfig)

	us := userservice.NewUserServiceImpl(storage.GetUserPersistence(), coze.GetCozeService(), jwtUtil)
	mms := mindmapservice.NewMindMapServiceImpl(storage.GetMindMapPersistence())
	cs := cosservice.NewCOSServiceImpl(cosService, cosConfig)
	handler.MustInitHandler(us, mms, cs)

	// 初始化JWT鉴权中间件
	router.InitJWTAuth(us)

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
