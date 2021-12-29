package dataStore

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/go-sql-driver/mysql"
)

const recordStatus int = 1

type DBConnectionInfo struct {
	UserName string
	Password string
	Host     string
	Database string
	Port     string
	Charset  string
}

/*
  レコード更新
*/
func UpdateDBData() {
	log.Println("レコード更新開始")

	// 実処理
	db, _ := dbConnect()
	stmtUpdate, err := db.Prepare("UPDATE demo_table SET status=?, updated_at=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtUpdate.Close()
	result, err := stmtUpdate.Exec(recordStatus, time.Now(), 1)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		log.Println("レコード更新数の取得に失敗しました")
		rowsAffect = 0
	}
	defer db.Close()

	log.Println(fmt.Sprintf("更新レコード数: %d", rowsAffect))
	log.Println("データ更新処理終了")
}

func dbConnect() (*sql.DB, error) {
	connectInfo := getDBConnectionInfo()
	connectStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s",
		connectInfo.UserName,
		connectInfo.Password,
		connectInfo.Host,
		connectInfo.Port,
		connectInfo.Database,
		connectInfo.Charset,
	)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func getDBConnectionInfo() *DBConnectionInfo {
	return &DBConnectionInfo{
		UserName: getDBSecretValueFromSSM("/lambda/DB_USERNAME"),
		Password: getDBSecretValueFromSSM("/lambda/DB_PASSWORD"),
		Host:     getDBSecretValueFromSSM("/lambda/DB_HOST"),
		Database: getDBSecretValueFromSSM("/lambda/DB_DATABASE"),
		Port:     "3306",
		Charset:  "utf8"}
}

func getDBSecretValueFromSSM(key string) string {
	sess := session.Must(session.NewSession())
	ssmSdk := ssm.New(sess)
	res, err := ssmSdk.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}
	return *res.Parameter.Value
}
