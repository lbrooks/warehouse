package warehouse

import "os"

func EnvOrDefault(variableName, defaultValue string) string {
	if val, ok := os.LookupEnv(variableName); ok && val != "" {
		return val
	}
	return defaultValue
}
