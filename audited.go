// Package audited is used to log last UpdatedBy and CreatedBy for your models
//
package audited

import "gorm.io/gorm"

// User basic user model used by default AuditModel
type User struct {
	gorm.Model
	Name string
}

// Model make Model Auditable, embed `audited.Model` into your model as anonymous field to make the model auditable
// If you want a different user model just create your own base model and make sure it implements the auditableInterface
//    type User struct {
//      gorm.Model
//      audited.AuditedModel
//    }
type Model struct {
	gorm.Model
	CreatedByID uint
	UpdatedByID uint
	DeletedByID uint
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	UpdatedBy   User `gorm:"foreignKey:UpdatedByID"`
	DeletedBy   User `gorm:"foreignKey:DeletedByID"`
}

// SetCreatedBy set created by
func (model Model) SetCreatedBy(i interface{}) {
	model.CreatedByID = i.(uint)
}

// SetUpdatedBy set created by
func (model Model) SetUpdatedBy(i interface{}) {
	model.UpdatedByID = i.(uint)
}

// SetDeletedBy set created by
func (model Model) SetDeletedBy(i interface{}) {
	model.DeletedByID = i.(uint)
}
