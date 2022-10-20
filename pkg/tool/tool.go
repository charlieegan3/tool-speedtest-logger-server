package tool

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/tool-speedtest-logger-server/pkg/tool/handlers"

	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"
)

//go:embed migrations
var speedtestLoggerMigrations embed.FS

// SpeedtestLogger is a tool to take results from a speedtest cli
// and store them in a database
type SpeedtestLogger struct {
	config             *gabs.Container
	db                 *sql.DB
	username, password string
}

func (s *SpeedtestLogger) Name() string {
	return "speedtest-logger-server"
}

func (s *SpeedtestLogger) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config:   true,
		Database: true,
		HTTP:     true,
	}
}

func (s *SpeedtestLogger) SetConfig(config map[string]any) error {
	var path string
	var ok bool
	cfg := gabs.Wrap(config)

	path = "username"
	s.username, ok = cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}

	s.password, ok = cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}

	return nil
}
func (s *SpeedtestLogger) DatabaseMigrations() (*embed.FS, string, error) {
	return &speedtestLoggerMigrations, "migrations", nil
}
func (s *SpeedtestLogger) DatabaseSet(db *sql.DB) {
	s.db = db
}

func (s *SpeedtestLogger) HTTPPath() string { return "speedtest-logger" }
func (s *SpeedtestLogger) HTTPAttach(router *mux.Router) error {
	router.HandleFunc(
		"/report",
		handlers.BuildReportHandler(s.db, s.username, s.password),
	).Methods("POST")

	return nil
}

func (s *SpeedtestLogger) Jobs() ([]apis.Job, error)                              { return []apis.Job{}, nil }
func (s *SpeedtestLogger) ExternalJobsFuncSet(f func(job apis.ExternalJob) error) {}
