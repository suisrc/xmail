package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	g18n "github.com/suisrc/gin-i18n"
)

// 定义错误
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405
var (
	ErrNone                = &ErrorNone{}
	Err200Success          = &ErrorModel{Status: 200, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "", Other: "请求成功"}}
	Err200Redirect         = ErrorOfS(200, ShowPage, "R_REDIRECT_PAGE", "页面重定向")
	Err400BadParam         = ErrorOfS(400, ShowWarn, "P_SERVE_BAD-PARAM", "请求参数错误")
	Err400BadRequest       = ErrorOfS(400, ShowWarn, "S_SERVE_BAD-REQUEST", "请求发生错误")
	Err404NotFound         = ErrorOfS(404, ShowWarn, "S_SERVE_NOT-FOUND", "发出的请求针对的是不存在的记录，服务器没有进行操作")
	Err405MethodNotAllowed = ErrorOfS(405, ShowWarn, "S_SERVE_METHOD-NOT-ALLOWED", "请求的方法不允许")
	Err406NotAcceptable    = ErrorOfS(406, ShowWarn, "S_SERVE_NOT-ACCEPTABLE", "请求的格式不可得")
	Err429TooManyRequests  = ErrorOfS(429, ShowWarn, "S_SERVE_TOO-MANY-REQUESTS", "请求次数过多")
	Err500InternalServer   = ErrorOfS(500, ShowWarn, "S_SERVE_INTERNAL-SERVER", "服务器发生错误")

	Err401Unauthorized = ErrorOfS(401, ShowWarn, "S_SERVE_UNAUTHORIZED", "用户没有权限(令牌、用户名、密码错误)")
	Err403Forbidden    = ErrorOfS(403, ShowWarn, "S_SERVE_FORBIDDEN", "用户未得到授权，访问是被禁止的")
	Err403ErrorToken   = ErrorOfS(403, ShowWarn, "S_SERVE_ERR-TOKNE", "请求令牌异常")

	ErrNoneNilToken = ErrorOfS(401, ShowWarn, "P_UNAUTHORIZED_NONENIL-TOKEN", "用户没有权限(请求令牌为空)") // ErrNoneNilToken 没有令牌
	ErrExpiredToken = ErrorOfS(401, ShowWarn, "P_UNAUTHORIZED_EXPIRED-TOKEN", "用户没有权限(请求令牌过期)") // ErrExpiredToken 过期令牌
	ErrInvalidToken = ErrorOfS(401, ShowWarn, "S_UNAUTHORIZED_INVALID-TOKEN", "用户没有权限(请求令牌无效)") // ErrInvalidToken 无效令牌
	ErrSignoutToken = ErrorOfS(401, ShowWarn, "S_UNAUTHORIZED_SIGNOUT-TOKEN", "用户没有权限(请求令牌注销)") // ErrSignoutToken 登出令牌
)

// D -> object
type D interface{}

// H -> map
type H map[string]interface{}

//  错误类型	错误码约定	举例
//  参数异常	P_XX_XX		P_CAMPAIGN_NameNotNull: 运营活动名不能为空
//  业务异常	B_XX_XX		B_CAMPAIGN_NameAlreadyExist: 运营活动名已存在
//  系统异常	S_XX_XX		S_DATABASE_ERROR: 数据库错误
//  重定向		R_XX_XX		R_REDIRECT_PAGE: 页面重定向

const (
	ShowNone   = 0 // ShowNone 静音
	ShowWarn   = 1 // ShowWarn 消息警告
	ShowError  = 2 // ShowError 消息错误
	ShowNotify = 4 // ShowNotify 通知；
	ShowPage   = 9 // ShowPage 页
)

// ErrorModel 异常模型
type ErrorModel struct {
	Data         interface{}
	Status       int
	ShowType     int
	ErrorMessage *i18n.Message
	ErrorArgs    H
}

func (a *ErrorModel) Error() string {
	return fmt.Sprintf("[%d,%d]%s:%s", a.Status, a.ShowType, a.ErrorMessage.ID, a.ErrorMessage.Other)
}

func ErrorOf(code, message string) *ErrorModel {
	return &ErrorModel{Status: 200, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: code, Other: message}}
}

func ErrorOfS(status, stype int, code, message string) *ErrorModel {
	return &ErrorModel{Status: status, ShowType: stype, ErrorMessage: &i18n.Message{ID: code, Other: message}}
}

//===========================================================================================

func CanRefreshToken(err error) bool {
	if em, ok := err.(*ErrorModel); ok {
		return strings.HasPrefix(em.ErrorMessage.ID, "P_UNAUTHORIZED_")
	}
	return false
}

//===========================================================================================

// Error 异常的请求结果体
type Error struct {
	Success      bool        `json:"success"`                // 请求成功, false
	Data         interface{} `json:"data,omitempty"`         // 响应数据
	ErrorCode    string      `json:"errorCode,omitempty"`    // 错误代码
	ErrorMessage string      `json:"errorMessage,omitempty"` // 向用户显示消息
	ShowType     int         `json:"showType,omitempty"`     // 错误显示类型：0静音； 1条消息警告； 2消息错误； 4通知； 9页
	TraceID      string      `json:"traceId"`                // 请求ID
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.ErrorMessage)
}

//===========================================================================================

// Success 正常请求结构体
// ErrorCode, ErrorMessage, ShowType 为空
type Success struct { // Error
	Success bool        `json:"success"`           // 请求成功, false
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Total   int64       `json:"total,omitempty"`   // 总数
	TraceID string      `json:"traceId,omitempty"` // 请求ID
}

func (e *Success) Error() string {
	return "success"
}

//===========================================================================================

// ErrorRedirect 重定向
type ErrorRedirect struct {
	Status   int    // http.StatusSeeOther
	State    string // 状态, 用户还原现场
	Location string
}

func (e *ErrorRedirect) Error() string {
	return "Redirect: " + e.Location
}

func (e *ErrorRedirect) ErrCode() (code string, message string) {
	if e.State == "" {
		return "redirect", e.Location
	}
	return "redirect", "[state." + e.State + "]" + e.Location
}

//===========================================================================================

// ErrorData data内容
type ErrorData struct {
	Code        int
	ContentType string
	Data        []byte
}

func (e *ErrorData) Error() string {
	return "DATA: " + string(e.Data)
}

// ErrorHTML html内容
type ErrorHTML struct {
	Status int
	Name   string
	Obj    interface{}
}

func (e *ErrorHTML) Error() string {
	return "HTML: " + e.Name
}

// ErrorNone 返回值已经被处理,无返回值
type ErrorNone struct {
}

func (e *ErrorNone) Error() string {
	return "none"
}

//===========================================================================================

// GetTraceID ...
func GetTraceID(ctx *gin.Context) string {
	if tid, ok := ctx.Get("x_request_id"); ok {
		return tid.(string)
	}
	// 优先从请求头中获取请求ID
	tid := ctx.GetHeader("X-Request-Id")
	// log.Println(traceID)
	if tid == "" {
		// 没有自建
		v, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		tid = v.String()
	}
	ctx.Set("x_request_id", tid)
	return tid
}

//===========================================================================================

// NewError 包装响应错误
func NewError(ctx *gin.Context, showType int, emsg *i18n.Message, args map[string]interface{}) *Error {
	res := &Error{
		Success:      false,
		ErrorCode:    emsg.ID,
		ErrorMessage: g18n.FormatMessage(ctx, emsg, args),
		ShowType:     showType,
		TraceID:      GetTraceID(ctx),
	}
	return res
}

// New0Error 包装响应错误, 没有参数
func New0Error(ctx *gin.Context, showType int, emsg *i18n.Message) *Error {
	return NewError(ctx, showType, emsg, nil)
}

// NewSuccess 包装响应结果
func NewSuccess(ctx *gin.Context, data interface{}) *Success {
	res := &Success{
		Success: true,
		Data:    data,
		TraceID: GetTraceID(ctx),
	}
	return res
}

// NewErrorWithData 包装响应错误
func NewErrorWithData(ctx *gin.Context, showType int, emsg *i18n.Message, args map[string]interface{}, data interface{}) *Error {
	res := &Error{
		Success:      false,
		Data:         data,
		ErrorCode:    emsg.ID,
		ErrorMessage: g18n.FormatMessage(ctx, emsg, args),
		ShowType:     showType,
		TraceID:      GetTraceID(ctx),
	}
	return res
}

// New0ErrorWithData 包装响应错误, 没有参数
func New0ErrorWithData(ctx *gin.Context, showType int, emsg *i18n.Message, data interface{}) *Error {
	return NewErrorWithData(ctx, showType, emsg, nil, data)
}

// NewWrapError 包装响应错误
func NewWrapError(ctx *gin.Context, em *ErrorModel) *Error {
	if em.ErrorMessage.ID == "" {
		res := &Error{
			Success: true,
			Data:    em.Data,
			TraceID: GetTraceID(ctx),
		}
		return res
	}
	res := &Error{
		Success:      false,
		Data:         em.Data,
		ErrorCode:    em.ErrorMessage.ID,
		ErrorMessage: g18n.FormatMessage(ctx, em.ErrorMessage, em.ErrorArgs),
		ShowType:     em.ShowType,
		TraceID:      GetTraceID(ctx),
	}
	return res
}

// NewWrapError 包装响应错误
func NewWrapErrorWithData(ctx *gin.Context, em *ErrorModel, data interface{}) *Error {
	if em.ErrorMessage.ID == "" {
		res := &Error{
			Success: true,
			Data:    data,
			TraceID: GetTraceID(ctx),
		}
		return res
	}
	res := &Error{
		Success:      false,
		Data:         data,
		ErrorCode:    em.ErrorMessage.ID,
		ErrorMessage: g18n.FormatMessage(ctx, em.ErrorMessage, em.ErrorArgs),
		ShowType:     em.ShowType,
		TraceID:      GetTraceID(ctx),
	}
	return res
}

// NewWrap400BadParams 无法解析异常
func NewWrap400BadParams(ctx *gin.Context, err error) *ErrorModel {
	return &ErrorModel{
		Status:       Err400BadParam.Status,
		ShowType:     Err400BadParam.ShowType,
		ErrorMessage: &i18n.Message{ID: "P_SERVE_BAD-PARAMS", Other: "请求参数错误:{{.error}}"},
		ErrorArgs:    H{"error": err.Error()},
	}
}

func ErrToHtml(c *gin.Context, err error) {
	// 将json异常转为html异常
	tmpl := gin.H{"title": "错误"}
	switch err := err.(type) {
	case *Error:
		tmpl["errorCode"] = err.ErrorCode
		tmpl["errorMessage"] = err.ErrorMessage
		tmpl["traceId"] = err.TraceID
	case *ErrorModel:
		er2 := NewWrapError(c, err)
		tmpl["errorCode"] = er2.ErrorCode
		tmpl["errorMessage"] = er2.ErrorMessage
		tmpl["traceId"] = er2.TraceID
	default:
		// log.Errorf(ctx, log.ErrorWithStack(err))
		er2 := NewWrapError(c, Err500InternalServer)
		tmpl["errorCode"] = er2.ErrorCode
		tmpl["errorMessage"] = er2.ErrorMessage
		tmpl["traceId"] = er2.TraceID
	}
	// 输出html文本
	c.HTML(http.StatusOK, "error1.html", tmpl)
	c.Abort()
}
