package app

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
)

type appConfig struct {
	Logger struct {
		LogLvl string // debug, info, error
	}

	MainnetConfig *liteclient.GlobalConfig

	Wallet struct {
		Seed          []string
		AddressParse  string
		JettonAddress string
		UuidGoogle          string
	}
}

var CFG *appConfig = &appConfig{}

func InitConfig() error {
	godotenv.Load(".env")

	CFG.Logger.LogLvl = os.Getenv("LOG_LVL")

	jsonConfig, err := os.Open("mainnet-config.json")
	if err != nil {
		return err
	}

	if err := json.NewDecoder(jsonConfig).Decode(&CFG.MainnetConfig); err != nil {
		return err
	}
	defer jsonConfig.Close()

	CFG.Wallet.Seed = strings.Split(os.Getenv("SEED"), " ")

	CFG.Wallet.AddressParse = os.Getenv("ADDRESS_PARSE")

	CFG.Wallet.JettonAddress = os.Getenv("JETTON_ADDRESS")

	CFG.Wallet.UuidGoogle = uuid.New().String()
	logrus.Info("AddressParse: ", CFG.Wallet.AddressParse)

	logrus.Info("JettonAddress: ", CFG.Wallet.JettonAddress)

	logrus.Info("Uuid: ", CFG.Wallet.UuidGoogle)


	return nil
}
