package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var (
		host = "127.0.0.1"
		// host = "10.108.195.0"
		user     = "postgres"
		password = "LKCxjisYkxvsci68aylEmasXWCo9DSu8CEk6LFzhnHzCjM640IOaBoSPTnUcnPPv"
		dbname   = "postgres"
		port     = 5432
	)
	dsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(dsn, host, user, password, dbname, port)), &gorm.Config{})
	fmt.Printf("%#v\nerr=%s", db, err)
	fmt.Printf("db.config = %#v\n", db.Config)
	// defer db.Close()

	employee := Company{
		Name:     "Lily",
		Age:      20,
		Address:  "Shenzhen",
		Salary:   20000,
		JoinDate: time.Now(),
	}
	// db.AutoMigrate(&employee)
	result := db.Create(&employee)
	fmt.Printf("employee.id = %d, result = %#v\n, error = %v", employee.ID, result, result.Error)

	var insertRecords []fakeTable
	for i := 0; i < 10; i++ {
		insertRecords = append(insertRecords,
			fakeTable{
				Name:  fmt.Sprintf("name%d", i),
				Email: fmt.Sprintf("test%d@test.com", i),
				// you don't need to set CreatedAt, UpdatedAt
			},
		)
	}

	for _, record := range insertRecords {
		result := db.Create(&record)
		fmt.Printf("record.id = %d, error = %v\n", record.ID, result.Error)
	}
	// err = gormbulk.BulkInsert(db, insertRecords, 3000)
	// if err != nil {
	// do something
	// }

	// columns you want to exclude from Insert, specify as an argument
	// err = gormbulk.BulkInsert(db, insertRecords, 3000, "Email")
	// if err != nil {
	// do something
	// }
}

type fakeTable struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (fakeTable) TableName() string {
	return "fake_table"
}

type Company struct {
	ID       uint `gorm:"autoIncrement"`
	Name     string
	Age      uint8
	Address  string
	Salary   int64
	JoinDate time.Time
}

// TableName 会将 Company 的表名重写为 `company`
func (Company) TableName() string {
	return "company"
}
