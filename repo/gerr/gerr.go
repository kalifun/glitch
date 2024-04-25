package gerr

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ErrorCode struct {
	Code    string
	Key     string
	Message ErrorMessage
}

type ErrorMessage struct {
	CN string
	EN string
}

var errorMap = make(map[string]ErrorCode)

type GError struct {
	key    string
	code   string
	format string
	args   []interface{}
}

func New(key, code, format string) *GError {
	return &GError{
		key:    key,
		code:   code,
		format: format,
	}
}

func (g *GError) Error() string {
	msg := fmt.Sprintf(g.format, g.args...)
	return fmt.Sprintf("Code:%s, Message: %s", g.code, msg)
}

func (g *GError) Args(args ...interface{}) *GError {
	newE := new(GError)
	newE.code = g.code
	newE.key = g.key
	newE.args = g.args
	newE.format = g.format
	newE.args = append(newE.args, args...)
	return newE
}

func (g *GError) ToGinH(lan string) gin.H {
	msg := "unknown"
	if err, ok := errorMap[g.key]; ok {
		if g.code == "" {
			g.code = err.Code
		}
		switch lan {
		case "en":
			msg = err.Message.EN
		case "cn":
			msg = err.Message.CN
		default:
			msg = err.Message.EN
		}
	} else {
		g.code = "unknow"
	}
	return gin.H{
		"error": map[string]string{
			"code":    g.code,
			"message": fmt.Sprintf(msg, g.args...),
		},
	}
}

func CreateError(data ErrorCode) {
	errorMap[data.Key] = data
}

func DescribeErrorMap() map[string]ErrorCode {
	return errorMap
}
