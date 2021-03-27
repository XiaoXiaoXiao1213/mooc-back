package base

type ResCode int

const (
	ResCodeOk ResCode = 200
	ResError  ResCode = -1
)

type Code struct {
	Val int
	Msg string
}

type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
