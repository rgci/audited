package audited_test

import (
	"os"
	"testing"

	"github.com/rgci/audited"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name string
	audited.AuditedModel
}

var db *gorm.DB

func testDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	return db, err
}

func TestMain(m *testing.M) {
	db, _ = testDB()
	db.AutoMigrate(audited.User{}, &Product{})
	audited.RegisterCallbacks(db)
	code := m.Run()
	DB, _ := db.DB()
	DB.Close()
	os.Remove("./test.db")
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	user := audited.User{Name: "grande"}
	db.Save(&user)
	db := db.Set("audited:current_user", user)

	product := Product{Name: "product1"}
	db.Debug().Save(&product)

	if product.GetCreatedBy() != int(user.ID) {
		t.Errorf("created_by is not equal current user")
	}

	product.Name = "product_new"
	db.Save(&product)
	if product.GetUpdatedBy() != int(user.ID) {
		t.Errorf("updated_by is not equal current user")
	}
}
