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
	DB       string        `envconfig:"db" default:"root:timberslide@tcp(localhost:3306)/u1m"`
	Interval time.Duration `envconfig:"interval" default:"24h"`
	ZipURL   string        `envconfig:"zip_url" default:"http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"`
}

func load(db *sqlx.DB, day string, zipURL string) error {
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

	logrus.Debugln("Reading in values")
	var valueString []string
	var valueArgs []interface{}
	count := 0
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		rank := parts[0]
		name := parts[1]
		valueString = append(valueString, "(?, ?, ?)")
		valueArgs = append(valueArgs, name)
		valueArgs = append(valueArgs, rank)
		valueArgs = append(valueArgs, day)
		count++
		if count < 10000 {
			continue
		}
		count = 0
		logrus.Debugln("Starting bulk insert")
		query := fmt.Sprintf("INSERT INTO ranking (name, rank, day) VALUES %s", strings.Join(valueString, ","))
		if _, err = db.Exec(query, valueArgs...); err != nil {
			logrus.WithField("mysql", "insert").Errorln(err)
		}
		valueString = nil
		valueArgs = nil
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
		day := time.Now().Format("2006-01-02")
		if err = load(db, day, s.ZipURL); err != nil {
			logrus.Errorln(err)
		}
		time.Sleep(s.Interval)
	}
}
