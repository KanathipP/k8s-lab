package config

import (
	"os"
	"strings"
)

type Config struct {
	Namespace        string
	LabelSelector    string
	CORSAllowOrigins []string
}

func Load() Config {
	ns := os.Getenv("TARGET_NAMESPACE")
	ls := os.Getenv("LABEL_SELECTOR")
	corsEnv := os.Getenv("CORS_ALLOW_ORIGINS")
	var origins []string
	if corsEnv == "" {
		origins = []string{"*"}
	} else {
		for _, o := range strings.Split(corsEnv, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
	}

	return Config{
		Namespace:        ns,
		LabelSelector:    ls,
		CORSAllowOrigins: origins,
	}
}
