// Package audited is used to log last UpdatedBy and CreatedBy for your models
//
package audited

// AuditedModel make Model Auditable, embed `audited.AuditedModel` into your model as anonymous field to make the model auditable
//    type User struct {
//      gorm.Model
//      audited.AuditedModel
//    }
type AuditedModel struct {
	CreatedBy int
	UpdatedBy int
	DeletedBy int
}

// SetCreatedBy set created by
func (model *AuditedModel) SetCreatedBy(createdBy int) {
	model.CreatedBy = createdBy
}

// GetCreatedBy get created by
func (model AuditedModel) GetCreatedBy() int {
	return model.CreatedBy
}

// SetUpdatedBy set updated by
func (model *AuditedModel) SetUpdatedBy(updatedBy int) {
	model.UpdatedBy = updatedBy
}

// GetUpdatedBy get updated by
func (model AuditedModel) GetUpdatedBy() int {
	return model.UpdatedBy
}

// SetDeletedBy set deleted by
func (model *AuditedModel) SetDeletedBy(deletedBy int) {
	model.DeletedBy = deletedBy
}

// GetDeletedBy get deleted by
func (model AuditedModel) GetDeletedBy() int {
	return model.DeletedBy
}
