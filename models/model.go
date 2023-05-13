package models

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./LabPoC.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type Bank struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Agency  string `json:"agency"`
	Account string `json:"account"`
}

func GetBanks(count int) ([]Bank, error) {

	rows, err := DB.Query("SELECT id, name, agency, account from bank LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bank := make([]Bank, 0)

	for rows.Next() {
		singleBank := Bank{}
		err = rows.Scan(&singleBank.Id, &singleBank.Name, &singleBank.Agency, &singleBank.Account)

		if err != nil {
			return nil, err
		}

		bank = append(bank, singleBank)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return bank, err
}

func GetBankByID(id string) (Bank, error) {

	stmt, err := DB.Prepare("SELECT id, name, agency, account from bank WHERE id = ?")

	if err != nil {
		return Bank{}, err
	}

	bank := Bank{}

	sqlErr := stmt.QueryRow(id).Scan(&bank.Id, &bank.Name, &bank.Agency, &bank.Account)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Bank{}, nil
		}
		return Bank{}, sqlErr
	}
	return bank, nil
}

func AddBank(newBank Bank) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO bank (name, agency, account) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newBank.Name, newBank.Agency, newBank.Account)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateBank(ourBank Bank, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE bank SET name = ?, agency = ?, account = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourBank.Name, ourBank.Agency, ourBank.Account, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteBank(bankId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from bank where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(bankId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
