package audited

import (
	"reflect"

	"gorm.io/gorm"
)

type auditableInterface interface {
	GetCreatedBy() int
	SetCreatedBy(i interface{})
	GetUpdatedBy() int
	SetUpdatedBy(i interface{})
	GetDeletedBy() int
	SetDeletedBy(i interface{})
}

func isAuditable(db *gorm.DB) (isAuditable bool) {
	if db.Statement.Model == nil {
		return false
	}
	_, isAuditable = db.Statement.Model.(auditableInterface)
	return
}

func getCurrentUser(db *gorm.DB) (uint64, bool) {
	user, _ := db.Get("audited:current_user")
	rv := reflect.ValueOf(user)
	if rv.Kind() != reflect.Struct {
		return 0, false
	}
	field := rv.FieldByName("ID")
	if field.Kind() != reflect.Uint {
		return 0, false
	}
	return field.Uint(), true
}

func assignCreatedBy(db *gorm.DB) {
	if !isAuditable(db) {
		return
	}
	if user, ok := getCurrentUser(db); ok {
		db.Statement.SetColumn("CreatedByID", user)
	}
}

func assignUpdatedBy(db *gorm.DB) {
	if !isAuditable(db) {
		return
	}
	if user, ok := getCurrentUser(db); ok {
		if attrs, ok := db.InstanceGet("gorm:update_attrs"); ok {
			updateAttrs := attrs.(map[string]interface{})
			updateAttrs["updated_by"] = user
			db.InstanceSet("gorm:update_attrs", updateAttrs)
		} else {
			db.Statement.SetColumn("UpdatedByID", user)
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("audited:assign_created_by") == nil {
		callback.Create().Before("gorm:before_create").Register("audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().Before("gorm:before_update").Register("audited:assign_updated_by", assignUpdatedBy)
	}
}
