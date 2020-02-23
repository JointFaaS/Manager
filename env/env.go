package env

import "errors"

// Env defines the supported runtime env
type Env string

const (
	// PYTHON27 python 2.7
	PYTHON27 Env = "python27"

	// PYTHON3 python 3
	PYTHON3  Env = "python3"

	// JAVA8 java8
	JAVA8	 Env = "java8"
)

// ConvertStrToEnv validate the input str and return a valid env
func ConvertStrToEnv(str string) (Env, error) {
	if str == "python3" {
		return PYTHON3, nil
	} else {
		return "", errors.New("Not support env")
	}
}