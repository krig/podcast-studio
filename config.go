package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"os"
	"io"
	"log"
	"strconv"
	"encoding/json"
)

type Config struct {
	data map[string]string
}

func NewConfig(filename string) *Config {
	cfg := &Config{}
	cfg.data = make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for {
		var value [2]string
		if err := dec.Decode(&value); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		cfg.data[value[0]] = value[1]
	}
	return cfg
}

func (cfg *Config) String(key string) string {
	return cfg.data[key]
}

func (cfg *Config) Color(key string) sdl.Color {
	i, err := strconv.ParseInt(cfg.data[key], 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	return hexcolor(uint32(i))
}

func (cfg *Config) Texture(rend *sdl.Renderer, key string) *sdl.Texture {
	surface := sdl.Load(cfg.data[key])
	if surface == nil {
		log.Fatal("Failed to load " + cfg.data[key])
	}
	defer surface.Free()
	return rend.CreateTextureFromSurface(surface)
}

func (cfg *Config) Int(key string) int {
	i, err := strconv.ParseInt(cfg.data[key], 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(i)
}