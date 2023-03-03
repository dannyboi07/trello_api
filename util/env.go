package util

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadEnv() {
	f, err := ioutil.ReadFile(".env")
	if err != nil {
		Log.Fatalln("Failed to open .env, check the file's name or it's presence")
	}

	var newLineSplit []string = strings.Split(string(f), "\n")

	for i := 0; i < len(newLineSplit); i++ {
		var env = strings.SplitN(newLineSplit[i], "=", 2)
		if len(env) == 2 {
			os.Setenv(env[0], env[1])
		}
	}
}
