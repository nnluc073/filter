package filter

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Search structured representation of a search query.
type Search struct {
	Query    string
	Operator *Operator
	Fields   []string
}

// Scope returns the GORM scopes with the search query.
func (s *Search) Scope(schema *schema.Schema) func(*gorm.DB) *gorm.DB {
	if len(s.Fields) == 0 {
		return nil
	}

	return func(tx *gorm.DB) *gorm.DB {
		searchQuery := tx.Session(&gorm.Session{NewDB: true})

		for _, field := range s.Fields {

			f, sch, joinName := getField(field, schema, nil)
			if f == nil {
				continue
			}

			if joinName != "" {
				if err := tx.Statement.Parse(tx.Statement.Model); err != nil {
					tx.AddError(err)
					return tx
				}
				tx = join(tx, joinName, schema)
			}

			filter := &Filter{
				Field:    f.DBName,
				Operator: s.Operator,
				Args:     []string{s.Query},
				Or:       true,
			}

			tableName := tx.Statement.Quote(tableFromJoinName(sch.Table, joinName)) + "."
			searchQuery = s.Operator.Function(searchQuery, filter, tableName+tx.Statement.Quote(f.DBName), f.DataType)
		}

		return tx.Where(searchQuery)
	}
}
