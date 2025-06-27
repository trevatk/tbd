package setup

import "os"

func envLookup(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultValue
	}
	return v
}
