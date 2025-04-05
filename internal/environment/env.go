package env

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/joho/godotenv"
)

func getRoot() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../..")

}

type Environment struct {
	ApiAuth string `env:"API_AUTH"`
	GinMode string `env:"GIN_MODE"`
	Port    string `env:"PORT"`
}

var Env *Environment

func NewEnv() {
	var env Environment

	_ = godotenv.Load(path.Join(getRoot(), "/.env"))

	envVars := reflect.ValueOf(Environment{})
	for i := 0; i < envVars.NumField(); i++ {
		field := envVars.Type().Field(i)
		fieldEnvName := field.Tag.Get("env")
		if fieldEnvName == "" {
			continue
		}

		v, ok := os.LookupEnv(fieldEnvName)

		if !ok {
			panic("Env variable not found: " + fieldEnvName)
		}

		el := reflect.ValueOf(&env).Elem().FieldByName(field.Name)
		if el.Type().Kind() == reflect.String {
			el.SetString(v)
		} else if el.Type().Kind() == reflect.Bool {
			el.SetBool(v == "true")
		}
	}

	Env = &env
}
