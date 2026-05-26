package introspect

// Schema describes a database schema or namespace.
type Schema struct {
	Name string `json:"name"`
}

// Relation describes a table or view.
type Relation struct {
	Schema          string `json:"schema,omitempty"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	RowEstimate     *int64 `json:"row_estimate,omitempty"`
	RowEstimateKind string `json:"row_estimate_kind,omitempty"`
}

// Column describes a relation column.
type Column struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Nullable   bool    `json:"nullable"`
	Default    *string `json:"default"`
	PrimaryKey bool    `json:"primary_key"`
}

// Index describes a database index.
type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// ForeignKey describes a foreign-key relationship.
type ForeignKey struct {
	Name              string   `json:"name,omitempty"`
	Columns           []string `json:"columns"`
	ReferencedSchema  string   `json:"referenced_schema,omitempty"`
	ReferencedTable   string   `json:"referenced_table"`
	ReferencedColumns []string `json:"referenced_columns"`
}

// RelationSchema is the full metadata for one relation.
type RelationSchema struct {
	Table       Relation     `json:"table"`
	Columns     []Column     `json:"columns"`
	Indexes     []Index      `json:"indexes"`
	ForeignKeys []ForeignKey `json:"foreign_keys"`
}

// QueryResult contains columns and row values returned by a data or query command.
type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}
