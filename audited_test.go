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

type AuditedUser struct {
	gorm.Model
	Name        string
	CreatedByID uint
	UpdatedByID uint
	DeletedByID uint
	CreatedBy   *AuditedUser `gorm:"foreignKey:CreatedByID"`
	UpdatedBy   *AuditedUser `gorm:"foreignKey:UpdatedByID"`
	DeletedBy   *AuditedUser `gorm:"foreignKey:DeletedByID"`
}

// SetCreatedBy set created by
func (model AuditedUser) SetCreatedBy(i interface{}) {
	model.CreatedByID = i.(uint)
}

// SetUpdatedBy set created by
func (model AuditedUser) SetUpdatedBy(i interface{}) {
	model.UpdatedByID = i.(uint)
}

// SetDeletedBy set created by
func (model AuditedUser) SetDeletedBy(i interface{}) {
	model.DeletedByID = i.(uint)
}

type Product struct {
	audited.Model
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

type Company struct {
	gorm.Model
	Name        string
	CreatedByID uint
	UpdatedByID uint
	DeletedByID uint
	CreatedBy   *audited.User `gorm:"foreignKey:CreatedByID"`
	UpdatedBy   *audited.User `gorm:"foreignKey:UpdatedByID"`
	DeletedBy   *audited.User `gorm:"foreignKey:DeletedByID"`
}

// SetCreatedBy set created by
func (model Company) SetCreatedBy(i interface{}) {
	model.CreatedByID = i.(uint)
}

// SetUpdatedBy set created by
func (model Company) SetUpdatedBy(i interface{}) {
	model.UpdatedByID = i.(uint)
}

// SetDeletedBy set created by
func (model Company) SetDeletedBy(i interface{}) {
	model.DeletedByID = i.(uint)
}

var db *gorm.DB

func testDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	return db, err
}

func TestMain(m *testing.M) {
	db, _ = testDB()
	db.AutoMigrate(audited.User{}, &OtherUser{}, &AuditedUser{}, &Product{}, &Company{})
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

	ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), auditUser)
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

	ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), auditUser)
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

func TestMissingAuditUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), nil)
	db = db.WithContext(ctx)
	user := audited.User{Name: "test"}
	db.Create(&user)

	product := Product{
		Name: "test",
	}
	db.Create(&product)

	assert.Equal(t, product.CreatedByID, uint(0))

	product.Name = "product_new"
	db.Save(&product)
	assert.Equal(t, product.UpdatedByID, uint(0))
}

func TestCreateCompanySuccess(t *testing.T) {
	auditUser := audited.User{Name: "audit"}
	db.Create(&auditUser)

	ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), auditUser)
	db = db.WithContext(ctx)
	user := audited.User{Name: "test"}
	db.Create(&user)

	company := Company{Name: "test"}
	db.Create(&company)

	assert.Equal(t, company.CreatedByID, auditUser.ID)

	company.Name = "product_new"
	db.Save(&company)

	assert.Equal(t, company.UpdatedByID, auditUser.ID)
}

func TestCreateAuditedUserSuccess(t *testing.T) {
	auditUser := AuditedUser{Name: "audit"}
	db.Create(&auditUser)

	ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), auditUser)
	db = db.WithContext(ctx)
	user := AuditedUser{Name: "test"}
	db.Create(&user)

	assert.Equal(t, user.CreatedByID, auditUser.ID)

	auditUser.Name = "test2"
	db.Save(&auditUser)

	assert.Equal(t, auditUser.UpdatedByID, auditUser.ID)
}
