package lookup

type ResultColumn struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Hidden      bool   `json:"hidden"`
}

type ResultRow struct {
	Cells []string `json:"cells"`
}

type Result struct {
	Entity    string         `json:"entity"`
	KeyColumn string         `json:"key_column"`
	Columns   []ResultColumn `json:"columns"`
	Rows      []ResultRow    `json:"rows"`
}

func (c *Result) AddColumn(name string, displayName string, hidden bool) {
	c.Columns = append(c.Columns, ResultColumn{
		Name:        name,
		DisplayName: displayName,
		Hidden:      hidden,
	})
}

func (c *Result) AddRow1(cell1 string) {
	var row ResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	c.Rows = append(c.Rows, row)
}

func (c *Result) AddRow2(cell1 string, cell2 string) {
	var row ResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	row.Cells = append(row.Cells, cell2)
	c.Rows = append(c.Rows, row)
}

func (c *Result) AddRow3(cell1 string, cell2 string, cell3 string) {
	var row ResultRow
	row.Cells = make([]string, 0)
	row.Cells = append(row.Cells, cell1)
	row.Cells = append(row.Cells, cell2)
	row.Cells = append(row.Cells, cell3)
	c.Rows = append(c.Rows, row)
}
