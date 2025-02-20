package yaml

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jeremywohl/flatten"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func load(data []byte) error {
	config := make(map[string]interface{})
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	flat, _ := flatten.Flatten(config, "", flatten.UnderscoreStyle)
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}
	for k, v := range flat {
		upperK := strings.ToUpper(k)
		if !currentEnv[upperK] {
			_ = os.Setenv(upperK, fmt.Sprint(v))
		}
	}
	return nil
}

func LoadReader(reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return load(data)
}

func loadFile(file string) {
	data, _ := ioutil.ReadFile(file)
	err := load(data)
	if err != nil {
		panic(err)
	}
}

func Load(env string) {
	wd, _ := os.Getwd()
	matches, _ := filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-%s-local.%s", env, "y*ml")))
	for _, item := range matches {
		loadFile(item)
	}
	if "test" != env {
		matches, _ = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-local.%s", "y*ml")))
		for _, item := range matches {
			loadFile(item)
		}
	}
	matches, _ = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-%s.%s", env, "y*ml")))
	for _, item := range matches {
		loadFile(item)
	}
	matches, _ = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app.%s", "y*ml")))
	for _, item := range matches {
		loadFile(item)
	}
}

func LoadReaderAsMap(reader io.Reader) (map[string]interface{}, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return loadAsMap(data)
}

func LoadFileAsMap(file string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return loadAsMap(data)
}

func loadAsMap(data []byte) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return flatten.Flatten(config, "", flatten.DotStyle)
}
