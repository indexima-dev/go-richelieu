package generator

type Plan struct {
	Schemas []Schema `json:"schemas"`
}

type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

type Table struct {
	Name    string   `json:"name"`
	Rows    int      `json:"rows"`
	Files   int      `json:"files"`
	Mode    string   `json:"mode,omitempty"` // Can be one of RANDOM, ALTERNATE, BLOCK
	Columns []Column `json:"columns"`
}

func (plan *Plan) HasSchema(name string) int {
	for index := range plan.Schemas {
		if plan.Schemas[index].Name == name {
			return index
		}
	}
	return -1
}

func (schema *Schema) HasTable(name string) int {
	for index := range schema.Tables {
		if schema.Tables[index].Name == name {
			return index
		}
	}
	return -1
}

func (table *Table) HasColumn(name string) int {
	for index := range table.Columns {
		if table.Columns[index].Name == name {
			return index
		}
	}
	return -1
}

// Merge parts of plan
func MergePlanParts(plan *Plan, added *Plan) {
	if added == nil {
		return
	}
	for _, s := range added.Schemas {
		var i = plan.HasSchema(s.Name)
		if i < 0 { // Schema is not already in this plan, just add it
			plan.Schemas = append(plan.Schemas, s)
		} else { // Schema is already present, merge the tables
			mergeSchemaParts(&plan.Schemas[i], &s)
		}
	}
}

// Merge parts of schema
func mergeSchemaParts(schema *Schema, added *Schema) {
	if added == nil {
		return
	}
	for _, t := range added.Tables {
		var i = schema.HasTable(t.Name)
		if i < 0 { // Table is not already in this plan, just add it
			schema.Tables = append(schema.Tables, t)
		} else { // Table is already present, merge the tables
			mergeTableParts(&schema.Tables[i], &t)
		}
	}

}

// Merge parts of table
func mergeTableParts(table *Table, added *Table) {
	if added == nil {
		return
	}
	if added.Rows != 0 {
		table.Rows = added.Rows
	}
	if added.Files != 0 {
		table.Files = added.Files
	}
	if added.Mode != "" {
		table.Mode = added.Mode
	}
	if len(added.Columns) > 0 {
		table.Columns = added.Columns
	}
}

// Merge parts of column
func mergeColumnParts(column *Column, added *Column) {
	if added == nil {
		return
	}
	if added.Type != "" {
		column.Type = added.Type
	}
	if added.Offset != 0 {
		column.Offset = added.Offset
	}
	if added.Distinct != 0 {
		column.Distinct = added.Distinct
	}
	if added.Mode != "" {
		column.Mode = added.Mode
	}
	if added.Start != "" {
		column.Start = added.Start
	}
	if added.End != "" {
		column.End = added.End
	}
	if added.Prefix != "" {
		column.Prefix = added.Prefix
	}
	if added.ValuesList != "" {
		column.ValuesList = added.ValuesList
	}
	if added.BlockStep != 0 {
		column.BlockStep = added.BlockStep
	}
}
