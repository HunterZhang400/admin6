package common

// WoodeBoxAPIResponse 通用API响应结构
type RangeDataAPIResponse struct {
	Code      int            `json:"code"`
	Msg       string         `json:"msg"`
	Data      map[string]any `json:"data"`
	Timestamp int64          `json:"timestamp"`
}

type TotalDataAPIResponse struct {
	Code      int        `json:"code"`
	Msg       string     `json:"msg"`
	Data      *TotalData `json:"data"`
	Timestamp int64      `json:"timestamp"`
}

type TotalData struct {
	CashTimes         int     `json:"cash_times"`
	CashMoney         int     `json:"cash_money"`
	MyRechargeTimes   int     `json:"my_recharge_times"`
	MyRechargeMoney   int     `json:"my_recharge_money"`
	TotalAvailableBet float64 `json:"total_available_bet"`
}

// UserTotalData 用户总数据
type UserTotalData struct {
	CashTimes         int     `json:"cash_times"`          // 提现次数
	CashMoney         float64 `json:"cash_money"`          // 累计提现金额
	MyRechargeTimes   int     `json:"my_recharge_times"`   // 充值次数
	MyRechargeMoney   float64 `json:"my_recharge_money"`   // 累计充值金额
	TotalAvailableBet float64 `json:"total_available_bet"` // 累计有效投注
}

// UserDataByTimeRange 时间范围内用户数据
type UserDataByTimeRange struct {
	Money         float64 `json:"money"`          // 账户当前余额
	Nickname      string  `json:"nickname"`       // 账户昵称
	ValidBet      float64 `json:"vaild_bet"`      // 时间范围内有效投注
	DeltaBetWin   float64 `json:"delta_bet_win"`  // 时间范围内游戏损益
	RechargeTimes int     `json:"recharge_times"` // 时间范围内充值次数
	Recharge      float64 `json:"recharge"`       // 时间范围内充值金额
	CashTimes     int     `json:"cash_times"`     // 时间范围内提现次数
	Cash          float64 `json:"cash"`           // 时间范围内提现金额
}

// GetUserTotalDataRequest 获取用户总数据请求参数
type GetUserTotalDataRequest struct {
	Token  string `form:"token" binding:"required"`
	UserID string `form:"appId" binding:"required"` // appId字段实际存储的是userID
}

// GetUserDataByTimeRangeRequest 获取时间范围内用户数据请求参数
type GetUserDataByTimeRangeRequest struct {
	Token     string `form:"token" binding:"required"`
	UserID    string `form:"appId" binding:"required"` // appId字段实际存储的是userID
	StartTime int64  `form:"startTime" binding:"required"`
	EndTime   int64  `form:"endTime" binding:"required"`
}
