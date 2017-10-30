//
// Date: 10/19/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
  "os"
  //"log" 
  "time"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
)

// Database interface
type Datastore interface {

  // Accounts
  GetAllAcounts() []Account
  GetAccountById(uint) (Account, error)
  GetAccountByName(string) (Account, error)
  CreateAccount(*Account) error
  GetAccountsTotalBalance() float64
  GetAccountsPricePerUnit() float64

  // Account Marks
  MarkAccountByDate(uint, time.Time, float64) error
  GetMarksByAccountById(uint) []AccountMarks
  GetMarksByAccountByIdAndDate(uint, time.Time) (AccountMarks, error)

  // Units
  AddUnits(time.Time, float64, string) error
  GetUnitsTotalCount() float64

  // Ledgers
  GetAllLedgers() []Ledger
  CreateLedger(uint, time.Time, float64, string, string) (*Ledger, error)   

  // LedgerCategory
  GetLedgerCategoryById(uint) (LedgerCategory, error)

  // Marks
  GetAllMarks() []Mark
  MarkByDate(time.Time) error
  GetMarkByDate(time.Time) (*Mark, error)
  GetMarkAccountUnitsByAccountId(uint) float64
  
}

type DB struct {
  *gorm.DB
}

//
// Setup the db connection.
//
func NewDB() (*DB, error) {

  // Connect to Mysql
  db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + "/" + os.Getenv("DB_DATABASE") + "?charset=utf8&parseTime=True&loc=UTC")
  
  if err != nil {
    return nil, err
  }

  // Ping make sure the server is up.
  if err = db.DB().Ping(); err != nil {
    return nil, err
  }  

  // Enable
  //db.LogMode(true)
  //db.SetLogger(log.New(os.Stdout, "\r\n", 0))  

  // Run migrations
  db.AutoMigrate(&Unit{})
  db.AutoMigrate(&Mark{}) 
  db.AutoMigrate(&Ledger{})    
  db.AutoMigrate(&Account{})
  db.AutoMigrate(&AccountMarks{})
  db.AutoMigrate(&AccountUnits{})
  db.AutoMigrate(&LedgerCategory{})

  // Return db connection.
  return &DB{db}, nil
}

/* End File */