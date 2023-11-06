package core

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ResSuccess 包装响应错误
// 禁止service层调用,请使用NewSuccess替换
func ResSuccess(ctx *gin.Context, v interface{}) error {
	if res, ok := v.(*Success); ok {
		ResJSON(ctx, http.StatusOK, res)
		return res
	} else {
		res := NewSuccess(ctx, v)
		ResJSON(ctx, http.StatusOK, res)
		//ctx.JSON(http.StatusOK, res)
		//ctx.Abort()
		return res
	}
}

// ResError 包装响应错误
// 禁止service层调用,请使用NewWarpError替换
func ResError(ctx *gin.Context, em *ErrorModel) error {
	res := NewWrapError(ctx, em)
	ResJSON(ctx, em.Status, res)
	return res
}

// ResJSON 响应JSON数据
// 禁止service层调用
func ResJSON(ctx *gin.Context, status int, v interface{}) {
	if ctx == nil {
		return
	}
	var data []byte
	switch v := v.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		if buf, err := json.Marshal(v); err != nil {
			panic(err)
		} else {
			data = buf
		}
	}

	if status == 0 {
		status = http.StatusOK
	}
	ctx.Data(status, "application/json; charset=utf-8", data)
	ctx.Abort()

	// ctx.JSON(status, v)
	// ctx.PureJSON(status, v)
}

// ResError1 处理了返回值
func ResError1(cx *gin.Context, err error, uem *ErrorModel) {
	if FixError0(cx, err) {
		return
	}
	logrus.Error(ErrorWithStack1(err))
	ResError(cx, uem) // 返回备用异常
}

// ResError2 处理了返回值
func ResError2(cx *gin.Context, err error, uem *ErrorModel) {
	logrus.Error(ErrorWithStack1(err))
	ResError(cx, uem) // 返回备用异常
}

// FixError0 上级应用已经处理了返回值
func FixError0(ctx *gin.Context, err error) bool {
	if err == nil {
		return true
	}
	switch err := err.(type) {
	case *Success, *Error:
		ResJSON(ctx, http.StatusOK, err)
		return true
	case *ErrorRedirect:
		status := err.Status
		if status <= 0 {
			status = http.StatusFound // 303 -> 302
		}
		ctx.Redirect(status, err.Location)
		ctx.Abort()
		return true
	case *ErrorData:
		status := err.Code
		if status <= 0 {
			status = http.StatusOK
		}
		ctx.Data(status, err.ContentType, err.Data)
		ctx.Abort()
		return true
	case *ErrorHTML:
		status := err.Status
		if status <= 0 {
			status = http.StatusOK
		}
		ctx.HTML(status, err.Name, err.Obj)
		ctx.Abort()
		return true
	case *ErrorNone:
		// do nothing
		return true
	case *ErrorModel:
		ResJSON(ctx, err.Status, NewWrapError(ctx, err))
		return true
	default:
		// e := err.Error()
		return false
	}
}

// FixError 修复返回的异常
func FixError(ctx *gin.Context, err error, uem *ErrorModel, fix func()) {
	if FixError0(ctx, err) {
		return
	}
	if fix != nil {
		fix()
	}
	if uem != nil {
		ResError(ctx, uem)
	}
}

// Fix500Logger 修复返回的异常
func Fix500Logger(cc *gin.Context, err error) {
	FixError(cc, err, Err500InternalServer, func() { ErrorWithStack1(err) })
}

// FixError1 上级应用已经处理了返回值
func FixError1(cx *gin.Context, err error) (string, string) {
	switch err := err.(type) {
	case *Success, *ErrorNone:
		return "", ""
	case *ErrorRedirect:
		return err.ErrCode()
	case *Error:
		return err.ErrorCode, err.ErrorMessage
	case *ErrorModel:
		er1 := NewWrapError(cx, err)
		return er1.ErrorCode, er1.ErrorMessage
	default:
		ErrorWithStack1(err)
		er1 := NewWrapError(cx, Err500InternalServer)
		return er1.ErrorCode, er1.ErrorMessage
	}
}

// FixError2 上级应用已经处理了返回值
func FixError2(ctx *gin.Context, err error) error {
	if err == nil {
		return nil
	}
	switch err := err.(type) {
	case *Success, *Error:
		ResJSON(ctx, http.StatusOK, err)
		return nil
	case *ErrorRedirect:
		status := err.Status
		if status <= 0 {
			status = http.StatusFound // 303 -> 302
		}
		ctx.Redirect(status, err.Location)
		ctx.Abort()
		return nil
	case *ErrorNone:
		// do nothing
		return nil
	case *ErrorModel:
		ResJSON(ctx, err.Status, NewWrapError(ctx, err))
		if err.Status >= 400 {
			return ErrNone // 异常已经被处理,这依然是一个系统异常
		}
		return nil
	default:
		return err
	}
}
