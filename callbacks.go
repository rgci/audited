package audited

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

const recursionError = "Recursion detected"
const unknownTypeError = "Unknown type while checking for possible recursion"

// UserKey string value for user key to store in context
const UserKey = "current_user"
const createdColumn = "CreatedBy"
const updatedColumn = "UpdatedBy"
const deletedColumn = "DeletedBy"

// GormAuditKey for storing user context
type GormAuditKey string

// Audited plugin struct
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

func isSameTypeAuditField(db *gorm.DB, f string, u interface{}) bool {
	fieldLookup := db.Statement.Schema.LookUpField(f)
	typeofUser := reflect.TypeOf(u)
	if fieldLookup.IndirectFieldType != typeofUser {
		db.Logger.Error(db.Statement.Context, fmt.Sprintf(
			"Types %v and %v do not match. Audited column %v will not be set",
			fieldLookup.FieldType, typeofUser, f))
		return false
	}
	// catch recursion
	if f == createdColumn {
		primaryFieldName := fieldLookup.Schema.PrioritizedPrimaryField.Name
		reflectedFieldValue := db.Statement.ReflectValue.FieldByName(primaryFieldName)
		reflectedUserValue := reflect.ValueOf(u).FieldByName(primaryFieldName)
		indirectFieldValue := reflect.Indirect(reflectedFieldValue)
		indirectUserValue := reflect.Indirect(reflectedUserValue)
		switch indirectFieldValue.Interface().(type) {
		case uint:
			if indirectFieldValue.Interface().(uint) == indirectUserValue.Interface().(uint) {
				db.Logger.Info(db.Statement.Context, recursionError)
				return false
			}
		default:
			db.Logger.Error(db.Statement.Context, unknownTypeError)
			return false
		}
	}

	return true
}

func assignColumn(db *gorm.DB, c string) {
	if !isAuditable(db) {
		return
	}
	user := db.Statement.Context.Value(GormAuditKey(UserKey))
	if user == nil {
		return
	}
	if ok := isSameTypeAuditField(db, c, user); !ok {
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
