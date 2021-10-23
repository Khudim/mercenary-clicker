package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	TemplateMatcherUrl string      `yaml:"matcherUrl"`
	Templates          []*Template `yaml:"templates"`
}

type Template struct {
	id     string
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Debug  bool    `json:"debug"`
}

func fromProperties() AppConfig {
	var appConfig AppConfig

	if file, err := ioutil.ReadFile("./props.yaml"); err == nil {
		if err := yaml.Unmarshal(file, &appConfig); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	return appConfig
}