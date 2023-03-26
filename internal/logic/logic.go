package logic

import (
	"github.com/gin-gonic/gin"
	"lottery_server/internal/svc"
	"lottery_server/model"
	"math/rand"
	"time"
)

type Logic struct {
	svcCtx *svc.ServiceContext
}

func NewLogic(svcCtx *svc.ServiceContext) *Logic {
	return &Logic{svcCtx: svcCtx}
}

// http://localhost:8080/get_code?phone=12345678901
func (l *Logic) GetCode(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": "6666",
	})
}

type LoginReq struct {
	Phone string `form:"phone"`
	Code  string `form:"code"`
}

// http://localhost:8080/login?phone=12345678901&code=6666
func (l *Logic) Login(ctx *gin.Context) {
	req := LoginReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"token": req.Phone,
		"id":    req.Phone,
	})
}

// http://localhost:8080/prize_list
func (l *Logic) PrizeList(ctx *gin.Context) {
	var prizes []*model.Prize
	l.svcCtx.Sqlite().Model(&model.Prize{}).Find(&prizes)
	var prizeMap = make(map[int]*model.Prize)
	for _, prize := range prizes {
		prizeMap[prize.Id] = prize
	}
	var drawRecords []*model.DrawRecord
	l.svcCtx.Sqlite().Model(&model.DrawRecord{}).Find(&drawRecords)
	var resp []gin.H
	for _, drawRecord := range drawRecords {
		prize := prizeMap[drawRecord.PrizeId]
		var userCount int64
		l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where("draw_record_id = ?", drawRecord.Id).Count(&userCount)
		resp = append(resp, gin.H{
			"draw_record_id":  drawRecord.Id,
			"prize_name":      prize.Name,
			"is_end":          drawRecord.DrawTime > 0,
			"prize_desc":      prize.Desc,
			"prize_icon":      prize.Icon,
			"user_count":      userCount,
			"remaining_count": l.svcCtx.Config.DrawCount - userCount,
		})
	}
	ctx.JSON(200, gin.H{
		"prizes": resp,
	})
}

// http://localhost:8080/draw?user_id=12345678901&draw_record_id=1
// http://localhost:8080/draw?user_id=12345678902&draw_record_id=1
// http://localhost:8080/draw?user_id=12345678903&draw_record_id=1
type DrawReq struct {
	UserId       string `form:"user_id"`
	DrawRecordId int    `form:"draw_record_id"`
}

func (l *Logic) Draw(ctx *gin.Context) {
	req := DrawReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})
		return
	}
	// 查询当前进行的抽奖
	drawRecord := &model.DrawRecord{}
	err = l.svcCtx.Sqlite().Model(&model.DrawRecord{}).Where("id = ?", req.DrawRecordId).First(drawRecord).Error
	if err != nil {
		// 查询失败
		ctx.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		// 是否已开奖
		if drawRecord.DrawTime > 0 {
			ctx.JSON(500, gin.H{
				"err": "已开奖",
			})
		}
		// 有进行中的抽奖
		// 加入抽奖队列
		l.svcCtx.Sqlite().Create(&model.UserDrawRecord{
			UserId:       req.UserId,
			DrawRecordId: drawRecord.Id,
			Status:       0,
			CreateTime:   time.Now().Unix(),
			EndTime:      0,
		})
		// 查询人数是否够了  够了就开始抽奖
		var count int64
		l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where("draw_record_id = ?", drawRecord.Id).Count(&count)
		if count >= l.svcCtx.Config.DrawCount {
			// 开始抽奖
			// 随机获取一个用户
			var userDrawRecords []*model.UserDrawRecord
			l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where("draw_record_id = ?", drawRecord.Id).Find(&userDrawRecords)
			// 随机一个
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(userDrawRecords), func(i, j int) {
				userDrawRecords[i], userDrawRecords[j] = userDrawRecords[j], userDrawRecords[i]
			})
			userDrawRecord := userDrawRecords[0]
			// 更新抽奖记录
			l.svcCtx.Sqlite().Model(&model.DrawRecord{}).Where("id = ?", drawRecord.Id).Updates(map[string]interface{}{
				"draw_time": time.Now().Unix(),
				"user_id":   userDrawRecord.UserId,
			})
			// 更新用户抽奖记录
			l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where(
				"user_id = ? AND draw_record_id = ?", userDrawRecord.UserId, userDrawRecord.DrawRecordId).Updates(map[string]interface{}{
				"status":   1,
				"end_time": time.Now().Unix(),
			})
			l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where(
				"user_id != ? AND draw_record_id = ?", userDrawRecord.UserId, userDrawRecord.DrawRecordId).Updates(map[string]interface{}{
				"status":   2,
				"end_time": time.Now().Unix(),
			})
		}
		ctx.JSON(200, gin.H{
			"draw_record_id": drawRecord.Id,
		})
	}
}

type CheckDrawReq struct {
	DrawRecordId int `form:"draw_record_id"`
}

// http://localhost:8080/check_draw?draw_record_id=1
func (l *Logic) CheckDraw(ctx *gin.Context) {
	req := CheckDrawReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})
		return
	}
	// 查询当前进行的抽奖
	drawRecord := &model.DrawRecord{}
	l.svcCtx.Sqlite().Model(&model.DrawRecord{}).Where("id = ?", req.DrawRecordId).First(drawRecord)
	// 查询参与的人
	var userDrawRecords []*model.UserDrawRecord
	l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where("draw_record_id = ?", drawRecord.Id).Find(&userDrawRecords)
	// 奖品信息
	prize := &model.Prize{}
	l.svcCtx.Sqlite().Model(&model.Prize{}).Where("id = ?", drawRecord.PrizeId).First(prize)
	ctx.JSON(200, gin.H{
		"draw_record":       drawRecord,
		"user_draw_records": userDrawRecords,
		"prize":             prize,
	})
}

type DrawRecordReq struct {
	UserId string `form:"user_id"`
}

// http://localhost:8080/draw_record?user_id=12345678901
// DrawRecord 查询我自己的抽奖记录
func (l *Logic) DrawRecord(ctx *gin.Context) {
	req := DrawRecordReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})
		return
	}
	/*
		SELECT ud.UserId, ud.DrawRecordId, ud.Status, ud.CreateTime, ud.EndTime,
		dr.PrizeId, dr.DrawTime, dr.StartTime,
		p.Id, p.Name, p.Desc, p.Count, p.Icon
		FROM UserDrawRecord ud
		LEFT JOIN DrawRecord dr ON ud.DrawRecordId = dr.Id
		LEFT JOIN Prize p ON dr.PrizeId = p.Id
		WHERE ud.UserId = '用户ID'
	*/
	/*
			"draw_record_id":  drawRecord.Id,
		"prize_name":      prize.Name,
		"is_end":          drawRecord.DrawTime > 0,
		"prize_desc":      prize.Desc,
		"prize_icon":      prize.Icon,
		"user_count":      userCount,
		"remaining_count": l.svcCtx.Config.DrawCount - userCount,
	*/
	type Result struct {
		DrawRecordId   int    `json:"draw_record_id"`
		PrizeName      string `json:"prize_name"`
		IsEnd          bool   `json:"is_end"`
		Status         int    `json:"status"`
		DrawTime       int64  `json:"draw_time"`
		PrizeDesc      string `json:"prize_desc"`
		PrizeIcon      string `json:"prize_icon"`
		UserCount      int    `json:"user_count"`
		RemainingCount int    `json:"remaining_count"`
	}
	var results []*Result
	err = l.svcCtx.Sqlite().Table("user_draw_record").
		Select(""+
			"user_draw_record.draw_record_id AS draw_record_id, "+
			"user_draw_record.status AS status, "+
			"prize.name AS prize_name, "+
			"draw_record.draw_time AS draw_time, "+
			"prize.desc AS prize_desc, "+
			"prize.icon AS prize_icon").
		Joins("LEFT JOIN draw_record ON user_draw_record.draw_record_id = draw_record.id").
		Joins("LEFT JOIN prize ON draw_record.prize_id = prize.id").
		Where("user_draw_record.user_id = ?", req.UserId).
		Find(&results).Error
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询失败",
		})
		return
	}
	for _, result := range results {
		var userDrawRecords []*model.UserDrawRecord
		l.svcCtx.Sqlite().Model(&model.UserDrawRecord{}).Where("draw_record_id = ?", result.DrawRecordId).Find(&userDrawRecords)
		result.UserCount = len(userDrawRecords)
		result.RemainingCount = int(l.svcCtx.Config.DrawCount) - result.UserCount
		result.IsEnd = result.DrawTime > 0
	}
	ctx.JSON(200, gin.H{
		"draw_records": results,
	})
}
