package generator

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
