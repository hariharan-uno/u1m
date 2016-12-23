package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
)

// this service's configuration
type specification struct {
	DB       string        `envconfig:"db"`
	Interval time.Duration `envconfig:"interval" default:"1h"`
	ZipURL   string        `envconfig:"zip_url" default:"http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"`
}

func load(db *sqlx.DB, zipURL string) error {
	logrus.Debugln("Downloading top-1m.csv.zip")
	resp, err := http.Get(zipURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	zipFile := bytes.NewReader(respBytes)
	contents, err := zip.NewReader(zipFile, resp.ContentLength)
	if err != nil {
		return err
	}

	logrus.Debugln("Unzipping")
	var rc io.ReadCloser
	for _, f := range contents.File {
		if f.Name != "top-1m.csv" {
			continue
		}
		rc, err = f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		break
	}
	if rc == nil {
		return fmt.Errorf("Did not find top-1m.csv in zip file")
	}

	if _, err = db.Exec("CREATE TEMPORARY TABLE temp_current LIKE current"); err != nil {
		return err
	}
	logrus.Debugln("Reading in values")
	var valueString []string
	var valueArgs []interface{}
	count := 0
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		rank := parts[0]
		domain := parts[1]
		valueString = append(valueString, "(?, ?)")
		valueArgs = append(valueArgs, rank)
		valueArgs = append(valueArgs, domain)
		count++
		if count < 10000 {
			continue
		}
		count = 0
		logrus.Debugln("Starting bulk insert")
		query := fmt.Sprintf("INSERT INTO temp_current (rank, domain) VALUES %s", strings.Join(valueString, ","))
		if _, err = db.Exec(query, valueArgs...); err != nil {
			logrus.WithField("mysql", "insert").Errorln(err)
		}
		valueString = nil
		valueArgs = nil
	}
	if _, err = db.Exec("INSERT current SELECT rank, domain FROM temp_current ON DUPLICATE KEY UPDATE current.domain=temp_current.domain"); err != nil {
		return err
	}
	if _, err = db.Exec("DROP TABLE temp_current"); err != nil {
		return err
	}

	logrus.Infoln("Finished loading")
	return nil
}

func main() {
	var err error

	var s specification
	err = envconfig.Process("APP", &s)
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Debugln("Connecting to database")
	db, err := sqlx.Connect("mysql", s.DB)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer db.Close()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infoln("Starting loader")

	for {
		if err = load(db, s.ZipURL); err != nil {
			logrus.Errorln(err)
		}
		time.Sleep(s.Interval)
	}
}
