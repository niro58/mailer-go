package environment

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

type Env struct {
	ApiAuth      string `env:"API_AUTH"`
	ClientsPath  string `env:"CLIENTS_PATH"`
	TemplatesDir string `env:"TEMPLATES_DIR"`
	Mode         string `env:"MODE"`
}

var Environment *Env

func NewEnv() *Env {
	var env Env

	_ = godotenv.Load(path.Join(getRoot(), "/.env"), "")

	envVars := reflect.ValueOf(Env{})
	for i := 0; i < envVars.NumField(); i++ {
		field := envVars.Type().Field(i)
		fieldEnvName := field.Tag.Get("env")
		if fieldEnvName == "" {
			continue
		}

		v, ok := os.LookupEnv(fieldEnvName)

		if !ok {
			panic("Environment variable not found: " + fieldEnvName)
		}

		el := reflect.ValueOf(&env).Elem().FieldByName(field.Name)
		if el.Type().Kind() == reflect.String {
			el.SetString(v)
		} else if el.Type().Kind() == reflect.Bool {
			el.SetBool(v == "true")
		}
	}
	return &env
}
