/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"log"
	"os"
)

func GetEnvVarStringRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing %s environment variable", key)
	}

	return value
}
