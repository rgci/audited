package audited

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// UserKey string value for user key to store in context
const UserKey = "current_user"
const createdColumn = "CreatedBy"
const updatedColumn = "UpdatedBy"
const deletedColumn = "DeletedBy"

// GormAuditKey for storing user context
type GormAuditKey string

type Audited struct {
	*gorm.DB
}

type auditableInterface interface {
	SetCreatedBy(i interface{})
	SetUpdatedBy(i interface{})
	SetDeletedBy(i interface{})
}

func isAuditable(db *gorm.DB) (isAuditable bool) {
	if db.Statement.Model == nil {
		return false
	}
	_, isAuditable = db.Statement.Model.(auditableInterface)
	return isAuditable
}

func isSameTypeAuditField(db *gorm.DB, f string, u interface{}) error {
	fieldLookup := db.Statement.Schema.LookUpField(f)
	typeofUser := reflect.TypeOf(u)
	if fieldLookup.IndirectFieldType != typeofUser {
		return fmt.Errorf(
			"Types %v and %v do not match. Audited column %v will not be set",
			fieldLookup.FieldType, typeofUser, f)
	}
	return nil
}

func assignColumn(db *gorm.DB, c string) {
	if !isAuditable(db) {
		return
	}
	user := db.Statement.Context.Value(GormAuditKey(UserKey))
	if user == nil {
		return
	}
	if err := isSameTypeAuditField(db, c, user); err != nil {
		db.Logger.Error(db.Statement.Context, err.Error())
		return
	}
	db.Statement.SetColumn(c, user)
}

func assignCreatedBy(db *gorm.DB) { assignColumn(db, createdColumn) }

func assignUpdatedBy(db *gorm.DB) { assignColumn(db, updatedColumn) }

// New Instance of audited plugin
func New() *Audited {
	return &Audited{}
}

// Name of audited plugin
func (a *Audited) Name() string {
	return "gorm:audited"
}

// Initialize initializes plugin
func (a *Audited) Initialize(db *gorm.DB) error {
	callback := db.Callback()
	if callback.Create().Get("gorm:audited:assign_created_by") == nil {
		callback.Create().Before("gorm:before_create").Register("gorm:audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().Before("gorm:before_update").Register("gorm:audited:assign_updated_by", assignUpdatedBy)
	}
	return nil
}
