package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path               string `yaml:"path" default:"./"`
	SegmentSize        int64  `yaml:"segment_size" default:"128"`
	SyncMode           string `yaml:"sync_mode" default:"immediate"`
	CheckpointInterval int64  `yaml:"checkpoint_interval" default:"10000000"`
	Compression        bool   `yaml:"compression" default:"false"`
	GrpcPort           int64  `yaml:"grpc_port" default:"50051"`
}

func LoadConfig(path string) *Config {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// println(config)

	return &config
}
