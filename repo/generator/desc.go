package generator

import (
	"fmt"
	"strings"

	"github.com/kalifun/glitch/utils"
)

type ErrorDesc struct {
	Error []struct {
		Key     string `yaml:"key"`
		Code    string `yaml:"code"`
		Message struct {
			Cn string `yaml:"cn"`
			En string `yaml:"en"`
		} `yaml:"message"`
	} `yaml:"error"`
}

func (e ErrorDesc) ToString() string {
	var vars []string
	var inits []string
	var declare []string
	for i, v := range e.Error {
		if i == 0 {
			inits = append(inits, "func init() {\n")
			declare = append(declare, "var (\n")
		}

		low := utils.FirstLower(v.Key)
		s := fmt.Sprintf(DeclareErr, low, v.Key, v.Code, v.Message.Cn, v.Message.En)
		vars = append(vars, s+"\n")
		inits = append(inits, fmt.Sprintf("gerr.CreateError(%s)\n", low))
		up := utils.FirstUpper(v.Key)
		declare = append(declare, fmt.Sprintf(`%s = gerr.New("%s", "%s", "%s")`, up, v.Key, v.Code, v.Message.En)+"\n")

		if i == len(e.Error)-1 {
			inits = append(inits, "\n}")
			declare = append(declare, "\n)")
		}
	}
	return strings.Join(vars, "") + "\n" + strings.Join(inits, "") + "\n" + strings.Join(declare, "")
}
