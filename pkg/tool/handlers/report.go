package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/doug-martin/goqu/v9"
	"io"
	"net/http"
	"strconv"
	"time"
)

type result struct {
	ServerID      string `json:"id"`
	ServerName    string `json:"name"`
	ServerCountry string `json:"country"`
	Lat           string `json:"lat"`
	Lon           string `json:"lon"`
	Sponsor       string `json:"sponsor"`

	DlSpeed float64 `json:"dl_speed"`
	UlSpeed float64 `json:"ul_speed"`
	Latency int64   `json:"latency"`

	Client string `json:"client"`
}

func BuildReportHandler(db *sql.DB, username, password string) func(http.ResponseWriter, *http.Request) {
	goquDB := goqu.New("postgres", db)

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// check that the basic auth has been set
		requestUsername, requestPassword, ok := r.BasicAuth()
		if requestUsername != username || requestPassword != password && !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var res result
		err = json.Unmarshal(body, &res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		parsedLat, err := strconv.ParseFloat(res.Lat, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to parse lat"))
			return
		}
		parsedLon, err := strconv.ParseFloat(res.Lon, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to parse lon"))
			return
		}

		resultRecord := goqu.Record{
			"client":         res.Client,
			"server_id":      res.ServerID,
			"server_name":    res.ServerName,
			"server_country": res.ServerCountry,
			"sponsor":        res.Sponsor,
			"latitude":       parsedLat,
			"longitude":      parsedLon,

			"latency":  res.Latency,
			"dl_speed": res.DlSpeed,
			"ul_speed": res.UlSpeed,

			"created_at": time.Now(),
		}

		ins := goquDB.Insert("speedtest_logger_server.results").Rows(resultRecord).Executor()
		_, err = ins.Exec()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
