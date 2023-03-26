package server

import (
	"github.com/gin-gonic/gin"
	"lottery_server/internal/logic"
	"lottery_server/internal/svc"
)

type HttpServer struct {
	svcCtx *svc.ServiceContext
	engine *gin.Engine
}

func NewHttpServer(svcCtx *svc.ServiceContext) *HttpServer {
	return &HttpServer{svcCtx: svcCtx}
}

func (s *HttpServer) Start() {
	s.engine = gin.Default()

	s.initRouter()

	s.engine.Run(s.svcCtx.Config.ListenOn)
}

func (s *HttpServer) initRouter() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	l := logic.NewLogic(s.svcCtx)
	s.engine.Use(cors)
	// 获取验证码
	s.engine.GET("/get_code", l.GetCode)
	// 登录 手机号+固定验证码
	s.engine.GET("/login", l.Login)
	// 抽奖池列表
	s.engine.GET("/prize_list", l.PrizeList)
	// 抽奖
	s.engine.GET("/draw", l.Draw)
	// 查询是否中奖
	s.engine.GET("/check_draw", l.CheckDraw)
	// 查询自己抽奖记录
	s.engine.GET("/draw_record", l.DrawRecord)
}

func cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Next()
}
