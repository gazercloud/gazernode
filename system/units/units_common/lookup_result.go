package units_common

type LookupResultColumn struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type LookupResultRow struct {
	Cells []string `json:"cells"`
}

type LookupResult struct {
	Entity    string               `json:"entity"`
	KeyColumn string               `json:"key_column"`
	Columns   []LookupResultColumn `json:"columns"`
	Rows      []LookupResultRow    `json:"rows"`
}

func NewLookupResult() *LookupResult {
	var c LookupResult
	c.Columns = make([]LookupResultColumn, 0)
	c.Rows = make([]LookupResultRow, 0)
	return &c
}

func (c *LookupResult) AddColumn(name string, displayName string) {
	c.Columns = append(c.Columns, LookupResultColumn{
		Name:        name,
		DisplayName: displayName,
	})
}

func (c *LookupResult) AddRow1(cell1 string) {
	var row LookupResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	c.Rows = append(c.Rows, row)
}

func (c *LookupResult) AddRow2(cell1 string, cell2 string) {
	var row LookupResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	row.Cells = append(row.Cells, cell2)
	c.Rows = append(c.Rows, row)
}

func (c *LookupResult) AddRow3(cell1 string, cell2 string, cell3 string) {
	var row LookupResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	row.Cells = append(row.Cells, cell2)
	row.Cells = append(row.Cells, cell3)
	c.Rows = append(c.Rows, row)
}
