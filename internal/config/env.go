package config

import (
	"github.com/joho/godotenv"
	"github.com/om-baji/write-go/internal/utils"
)

func CheckVars() {
	err := godotenv.Load()

	if err != nil {
		utils.HandlExp(err)
	}
}
