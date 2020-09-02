package audited

import (
	"gorm.io/gorm"
)

type auditableInterface interface {
	SetCreatedBy(createdBy int)
	GetCreatedBy() int
	SetUpdatedBy(updatedBy int)
	GetUpdatedBy() int
	SetDeletedBy(updatedBy int)
	GetDeletedBy() int
}

func isAuditable(db *gorm.DB) (isAuditable bool) {
	if db.Statement.Model == nil {
		return false
	}
	_, isAuditable = db.Statement.Model.(auditableInterface)
	return
}

func getCurrentUser(db *gorm.DB) (uint, bool) {
	var user interface{}
	var hasUser bool

	user, hasUser = db.Get("audited:current_user")
	v, ok := user.(uint)
	if ok != true {
		return 0, ok
	}
	if hasUser {
		return v, true
	}

	return 0, false
}

func assignCreatedBy(db *gorm.DB) {
	if !isAuditable(db) {
		return
	}
	if user, ok := getCurrentUser(db); ok {
		db.Statement.SetColumn("CreatedBy", user)
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
			db.Statement.SetColumn("UpdatedBy", user)
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
