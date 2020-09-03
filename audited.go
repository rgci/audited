// Package audited is used to log last UpdatedBy and CreatedBy for your models
//
package audited

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
}

// GetID get id
func (model User) GetID() uint {
	return model.ID
}

// AuditedModel make Model Auditable, embed `audited.AuditedModel` into your model as anonymous field to make the model auditable
//    type User struct {
//      gorm.Model
//      audited.AuditedModel
//    }
type AuditedModel struct {
	CreatedByID int
	UpdatedByID int
	DeletedByID int
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	UpdatedBy   User `gorm:"foreignKey:UpdatedByID"`
	DeletedBy   User `gorm:"foreignKey:DeletedByID"`
}

// GetCreatedBy get created by
func (model AuditedModel) GetCreatedBy() int {
	return model.CreatedByID
}

// SetCreatedBy set created by
func (model AuditedModel) SetCreatedBy(i interface{}) {
	model.CreatedByID = i.(int)
}

// GetUpdatedBy get updated by
func (model AuditedModel) GetUpdatedBy() int {
	return model.UpdatedByID
}

// SetUpdatedBy set created by
func (model AuditedModel) SetUpdatedBy(i interface{}) {
	model.UpdatedByID = i.(int)
}

// GetDeletedBy get deleted by
func (model AuditedModel) GetDeletedBy() int {
	return model.DeletedByID
}

// SetDeletedBy set created by
func (model AuditedModel) SetDeletedBy(i interface{}) {
	model.DeletedByID = i.(int)
}
