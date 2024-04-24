package gerr

import "fmt"

type ErrorCode struct {
	Code string
	Key  string
	Msg  ErrorMessage
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

func CreateError(data ErrorCode) {
	errorMap[data.Key] = data
}

func DescribeErrorMap() map[string]ErrorCode {
	return errorMap
}
