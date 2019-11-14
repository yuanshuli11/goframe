package app

import (
	"go-war/app/controller/army"
	"go-war/app/controller/auth"
	"go-war/app/controller/home"
	"go-war/app/controller/maps"
	"go-war/app/controller/resource"
	"go-war/app/controller/tech"
	"go-war/app/controller/user"
	"go-war/app/middleware"
	_ "go-war/internal/init"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func init() {
	engine = gin.New()
	engine.Use(middleware.RequestInfo())
}

func GetRouters() *gin.Engine {

	authRouter := engine.Group("/auth", middleware.VerifySign())
	{

		authRouter.POST("/login", auth.Login())
		authRouter.POST("/register", auth.Register())

	}
	homeRouter := engine.Group("/", middleware.VerifySign(), middleware.VerifyLogin())
	{
		homeRouter.GET("/user/getinfo", user.GetInfo())
		homeRouter.POST("/user/init", user.InitUser())
		homeRouter.GET("/home/getResource", home.GetResource())

		homeRouter.GET("/army/getArmyBuilding", army.GetArmyBuilding())
		homeRouter.GET("/army/getArmyList", army.GetArmyList())
		homeRouter.GET("/army/getDefenseList", army.GetDefenseList())
		homeRouter.GET("/army/getBuildngDetail", army.GetBuildingDetail())

		homeRouter.GET("/resource/getResourceBuilding", resource.GetResourceBuilding())
		homeRouter.GET("/tech/getUserTech", tech.GetUserTech())
		homeRouter.GET("/maps/getMap", maps.GetMap())
	}
	return engine
}
