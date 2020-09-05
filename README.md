# Audited

Audited is used to record the last User who created and/or updated your [GORM](https://github.com/go-gorm/gorm) model. It does so using a `CreatedBy` and `UpdatedBy` field.

[![GoDoc](https://godoc.org/github.com/rgci/audited?status.svg)](https://godoc.org/github.com/rgci/audited)

### Register GORM Callbacks

Audited utilizes [GORM](https://github.com/jinzhu/gorm) callbacks to log data, so you will need to register callbacks first:

```go
import (
  "gorm.io/gorm"
  "github.com/rgci/audited"
)

db, err := gorm.Open("sqlite3", "demo_db")
db.Use(audited.New())
```

### Making a Model Auditable

Embed `audited.Model` into your model as an anonymous field to make the model auditable:

```go
type Product struct {
  gorm.Model
  Name string
  audited.Model
}
```

### Usage

```go
import "github.com/rgci/audited"
import "gorm.io/gorm"

func main() {
  var db, err = gorm.Open("sqlite3", "demo_db")
  var currentUser = User{ID: 100}
  var product Product

  // Create will set product's `CreatedBy`, `UpdatedBy` to `currentUser`'s primary key if `audited:current_user` is a valid model
  ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey) currentUser)
  db = db.WithContext(ctx)
  db.Create(&product)
  // product.CreatedBy => 100
  // product.UpdatedBy => 100

  // If it is not a valid model, then will set `CreatedBy`, `UpdatedBy` to default value
  ctx := context.WithValue(context.Background(), audited.GormAuditKey(audited.UserKey), nil)
  db = db.WithContext(ctx)
  db.Create(&product)
  // product.CreatedBy => 0
  // product.UpdatedBy => 0
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
