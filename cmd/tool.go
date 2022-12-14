package main

import (
	"context"
	"log"

	"github.com/charlieegan3/toolbelt/pkg/database"
	"github.com/charlieegan3/toolbelt/pkg/tool"
	"github.com/spf13/viper"

	speedtestLoggerServerTool "github.com/charlieegan3/tool-speedtest-logger-server/pkg/tool"
)

func main() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	params := viper.GetStringMapString("database.params")
	connectionString := viper.GetString("database.connectionString")
	db, err := database.Init(connectionString, params, params["dbname"], false)
	if err != nil {
		log.Fatalf("failed to init DB: %s", err)
	}

	tb := tool.NewBelt()
	tb.SetConfig(viper.GetStringMap("tools"))
	tb.SetDatabase(db)

	err = tb.AddTool(&speedtestLoggerServerTool.SpeedtestLogger{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	tb.RunServer(context.Background(), "0.0.0.0", "3000")
}
