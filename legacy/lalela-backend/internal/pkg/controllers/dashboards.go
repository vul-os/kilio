package controllers

import (
	"bytes"
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"database/sql"
	"fmt"
	_ "github.com/prestodb/presto-go-client/presto"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type DashsCon struct{}
func (t *DashsCon) DashboardsGet(r *http.Request, args *models.DashboardsEgRequest, reply *models.DashboardsEgResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}
	sqlquery := ""
	queryFile := "queries/" + args.Client + "/" + args.Dashboard + "/" + args.Query + ".sql"
	data := map[string]string{}
	for index, element := range args.Filters {
		data[index] = `'` + strings.Join(element[:], `','`) + `'`
	}
	fmt.Println(data)
	sqlquery, err = ParseTemplate(queryFile, data)
	if err != nil {
		log.Print(middleware.NewError(err))
	}
	fmt.Println(sqlquery)
	dsn := "https://imran:1mr4n@cognizance.page:8443?catalog=hive&schema=pns"
	db, _ := sql.Open("presto", dsn)
	fmt.Println(db)
	rows, err := db.Query(sqlquery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	dataPoints := services.DashboardsMap[args.Client][args.Dashboard][args.Query]
	DataPoints, err := dataPoints(rows)
	if err != nil {
		log.Fatal(err)
	}
	reply.DataPoints = DataPoints
	return nil
}
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}