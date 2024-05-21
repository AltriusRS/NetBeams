package environment

import (
	"github.com/joho/godotenv"
)

func LoadEnvFile(file string) error {
	return godotenv.Load(file)
}
