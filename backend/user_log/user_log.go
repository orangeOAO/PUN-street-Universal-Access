package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// 只少符合三種(數字、符號、英文大寫、英文小寫)
func PasswordCheck(passwd string) error {
	if len(passwd) < 6 {
		return errors.New("password too short")
	}

	indNum := [4]bool{false, false, false, false}
	spCode := "!@#$%^&*_-"

	for _, char := range passwd {
		switch {
		case 'A' <= char && char <= 'Z':
			indNum[0] = true
		case 'a' <= char && char <= 'z':
			indNum[1] = true
		case '0' <= char && char <= '9':
			indNum[2] = true
		case strings.ContainsRune(spCode, char):
			indNum[3] = true
		default:
			return errors.New("Unsupported character")
		}
	}

	codeCount := 0
	for _, i := range indNum {
		if i {
			codeCount++
		}
	}

	if codeCount < 3 {
		return errors.New("Too simple password")
	}

	return nil
}

// 手機號碼格式
func validPhoneNumber(phoneNumber string) error {
	regex := `^8869[0-9]{8}$|^09[0-9]{8}$`

	if match, _ := regexp.MatchString(regex, phoneNumber); match {
		return nil
	}
	return errors.New("Wrong Number format")

}

// 找重複email 有重複false
func searchMailAddress(db *sql.DB, address string) error {
	rows, err := db.Query("SELECT id, Email FROM UserData")
	if err != nil {
		fmt.Println("Error running query:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		if err := rows.Scan(&id, &email); err != nil {
			return errors.New("Error scanning row")
		}

		if strings.EqualFold(address, email) {
			return errors.New("This Email has been registered")
		}
	}

	return nil
}
func DecideWrite(tmp [4]bool) bool {

	var ans bool = true
	for i := range tmp {
		ans = ans && tmp[i]
	}
	return ans
}
func main() {

	_TFwrite := [4]bool{true, true, true, true} //TFwrite in db
	// 建立数据库连接
	connStr := "user=orange dbname=user_data sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()
	type UserData struct {
		Name, Account, Address, PhoneNumber, Password string
		Authority                                     int

		Birthday time.Time
	}
	//stirng2Time
	date, err := time.Parse("2006/01/02", "2002/12/17")
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}
	data := UserData{
		Name:        "XDD",
		Account:     "k98007@gmail.com",
		Address:     "Taiwan",
		Birthday:    date,
		PhoneNumber: "0912455672",
		Password:    "OAOrange1",
		Authority:   0,
	}

	if err := validPhoneNumber(data.PhoneNumber); err != nil {
		fmt.Println(err)
		_TFwrite[0] = false
	}

	if err := PasswordCheck(data.Password); err != nil {
		fmt.Println(err)
		_TFwrite[1] = false
	}

	if _, err := mail.ParseAddress(data.Account); err != nil {
		fmt.Println(err)
		_TFwrite[2] = false
	} else {
		if err2 := searchMailAddress(db, data.Account); err2 != nil {
			fmt.Println(err2)
			_TFwrite[3] = false

		}
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}

	// 执行插入操作
	if DecideWrite(_TFwrite) {
		_, err = tx.Exec("INSERT INTO UserData (Name, Password, Email, Address, PhoneNumber, Birthday, Authority) VALUES ($1, $2, $3, $4, $5, $6, $7)", data.Name, data.Password, data.Account, data.Address, data.PhoneNumber, data.Birthday, data.Authority)

		if err != nil {
			fmt.Println("Error inserting data:", err)
			tx.Rollback()
			return
		}
		// 提交事务
		err = tx.Commit()
		if err != nil {
			fmt.Println("Error committing transaction:", err)
			return
		}
	}

}
