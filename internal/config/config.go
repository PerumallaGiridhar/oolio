package config

import (
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout int
	ReadTimeout       int
	WriteTimeout      int
	IdleTimeout       int
}

type Config struct {
	Server     ServerConfig
	PromoFiles []string
}

func getEnvWithDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvIntWithDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			Addr:         getEnvWithDefault("ADDR", ":8080"),
			ReadTimeout:  getEnvIntWithDefault("READ_TIMEOUT", 5),
			WriteTimeout: getEnvIntWithDefault("WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvIntWithDefault("IDLE_TIMEOUT", 60),
		},
		PromoFiles: splitCSV(getEnvWithDefault("PROMO_FIELS", "/Users/giridhar/Downloads/safe_extract/couponbase1,/Users/giridhar/Downloads/safe_extract/couponbase2,/Users/giridhar/Downloads/safe_extract/couponbase3")),
	}
}
