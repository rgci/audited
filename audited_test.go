package audited_test

import (
	"context"
	"os"
	"testing"

	"github.com/rgci/audited"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OtherUser struct {
	gorm.Model
	Name string
}

type Product struct {
	audited.AuditedModel
	LinkedUserID int
	LinkedUser   audited.User `gorm:"foreignKey:LinkedUserID"`
	Name         string
}

// SetCreatedBy set created by
func (model Product) SetCreatedBy(i interface{}) {
	model.CreatedByID = i.(uint)
}

// SetUpdatedBy set created by
func (model Product) SetUpdatedBy(i interface{}) {
	model.UpdatedByID = i.(uint)
}

// SetDeletedBy set created by
func (model Product) SetDeletedBy(i interface{}) {
	model.DeletedByID = i.(uint)
}

var db *gorm.DB

func testDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	return db, err
}

func TestMain(m *testing.M) {
	db, _ = testDB()
	db.AutoMigrate(audited.User{}, &Product{}, &OtherUser{})
	db.Use(audited.New())
	code := m.Run()
	DB, _ := db.DB()
	DB.Close()
	os.Remove("./test.db")
	os.Exit(code)
}

func TestCreateUserSuccess(t *testing.T) {
	auditUser := audited.User{Name: "audit"}
	db.Create(&auditUser)

	ctx := context.WithValue(context.Background(), "gorm:audited:current_user", auditUser)
	db = db.WithContext(ctx)
	user := audited.User{Name: "test"}
	db.Create(&user)

	product := Product{
		Name: "test",
	}
	db.Create(&product)

	assert.Equal(t, product.CreatedByID, auditUser.ID)

	product.Name = "product_new"
	db.Save(&product)

	assert.Equal(t, product.UpdatedByID, auditUser.ID)
}

func TestCreateUserFail(t *testing.T) {
	auditUser := OtherUser{Name: "audit"}
	db.Create(&auditUser)

	ctx := context.WithValue(context.Background(), "gorm:audited:current_user", auditUser)
	db = db.WithContext(ctx)
	user := audited.User{Name: "test"}
	db.Create(&user)

	product := Product{
		Name: "test",
	}
	db.Create(&product)

	assert.NotEqual(t, product.CreatedByID, auditUser.ID)

	product.Name = "product_new"
	db.Save(&product)
	assert.NotEqual(t, product.UpdatedByID, auditUser.ID)
}
