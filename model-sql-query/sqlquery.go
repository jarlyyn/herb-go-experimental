package query

type Query interface {
	QueryCommand() string
	QueryArgs() []interface{}
}

func New(command string, args ...interface{}) *PlainQuery {
	return &PlainQuery{
		Command: command,
		Args:    args,
	}
}

type PlainQuery struct {
	Command string
	Args    []interface{}
}

func (q *PlainQuery) QueryCommand() string {
	return q.Command
}
func (q *PlainQuery) QueryArgs() []interface{} {
	return q.Args
}

func NewFromQuery() *FromQuery {
	return &FromQuery{
		Tables: [][2]string{},
	}

}

type FromQuery struct {
	Tables [][2]string
}

func (q *FromQuery) Add(tableName string, alias string) *FromQuery {
	q.Tables = append(q.Tables, [2]string{tableName, alias})
	return q
}
func (q *FromQuery) QueryCommand() string {
	var command = ""
	if len(q.Tables) == 0 {
		return command
	}
	command = "From "
	for k := range q.Tables {
		command += q.Tables[k][0] + " as " + q.Tables[k][1] + " , "
	}
	return command[:len(command)-3]
}
func (q *FromQuery) QueryArgs() []interface{} {
	return []interface{}{}
}

type QueryData struct {
	Field string
	Data  interface{}
	Raw   string
}

func NewInsertQuery(tableName string) *InsertQuery {
	return &InsertQuery{
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type InsertQuery struct {
	Prefix    *PlainQuery
	TableName string
	Data      []QueryData
}

func (q *InsertQuery) Add(field string, data interface{}) *InsertQuery {
	q.Data = append(q.Data, QueryData{Field: field, Data: data})
	return q
}
func (q *InsertQuery) AddRaw(field string, raw string) *InsertQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *InsertQuery) QueryCommand() string {
	var command = "INSERT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " INTO " + q.TableName
	var values = ""
	var columns = ""
	for k := range q.Data {
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
		columns += q.Data[k].Field + " , "
	}
	command += " ("
	command += columns[:len(columns)-3]
	command += " )"

	command += " VALUES ( "
	command += values[:len(values)-3]
	command += " )"
	return command
}
func (q *InsertQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		if q.Data[k].Data != nil {
			args = append(args, q.Data[k].Data)
		}
	}
	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	result = append(result, args...)
	return result
}

type DeleteQuery struct {
	TableName string
	Prefix    *PlainQuery
}

func (q *DeleteQuery) QueryCommand() string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " FROM " + q.TableName
	return command
}

func (q *DeleteQuery) QueryArgs() []interface{} {
	return q.Prefix.QueryArgs()
}

func NewDeleteQuery(tableName string) *DeleteQuery {
	return &DeleteQuery{
		Prefix:    New(""),
		TableName: tableName,
	}
}

func NewUpdateQuery(tableName string) *UpdateQuery {
	return &UpdateQuery{
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type UpdateQuery struct {
	Prefix    *PlainQuery
	TableName string
	Data      []QueryData
}

func (q *UpdateQuery) Add(field string, data interface{}) *UpdateQuery {
	q.Data = append(q.Data, QueryData{Field: field, Data: data})
	return q
}
func (q *UpdateQuery) AddRaw(field string, raw string) *UpdateQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *UpdateQuery) QueryCommand() string {
	var command = "UPDATE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " " + q.TableName
	command += " SET"
	var values = ""
	for k := range q.Data {
		values += q.Data[k].Field + " = "
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
	}
	command += values[:len(values)-3]
	return command
}
func (q *UpdateQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		if q.Data[k].Data != nil {
			args = append(args, q.Data[k].Data)
		}
	}
	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	result = append(result, args...)
	return result
}
func NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Prefix: New(""),
		Fields: []string{},
	}
}

type SelectQuery struct {
	Prefix *PlainQuery
	Fields []string
}

func (q *SelectQuery) Add(fields ...string) *SelectQuery {
	q.Fields = append(q.Fields, fields...)
	return q
}

func (q *SelectQuery) QueryCommand() string {
	var command = "SELECT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	var colunms = " "
	for k := range q.Fields {
		colunms += q.Fields[k] + " , "
	}
	command += colunms[:len(colunms)-3]
	return command
}
func (q *SelectQuery) QueryArgs() []interface{} {
	return q.Prefix.QueryArgs()
}
func (q *SelectQuery) Result() *SelectResult {
	return NewSelectResult(q.Fields)
}

func NewSelectResult(fields []string) *SelectResult {
	return &SelectResult{
		Fields: fields,
		Args:   make([]interface{}, len(fields)),
	}

}

type SelectResult struct {
	Fields []string
	Args   []interface{}
}

func (r *SelectResult) Bind(field string, arg interface{}) *SelectResult {
	for k := range r.Fields {
		if r.Fields[k] == field {
			r.Args[k] = arg
			return r
		}
	}
	return r
}

func NewWhereQuery() *WhereQurey {
	return &WhereQurey{
		Condition: New(""),
	}
}

type WhereQurey struct {
	Condition *PlainQuery
}

func (q *WhereQurey) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}
func (q *WhereQurey) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
func NewSelect(tableName string) *Select {
	return &Select{
		Select: NewSelectQuery(),
		From:   NewFromQuery(),
		Where:  NewWhereQuery(),
		Other:  New(""),
	}
}

type Select struct {
	Select *SelectQuery
	From   *FromQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (s *Select) Query() *PlainQuery {
	return Concat(s.Select, s.From, s.Where, s.Other)
}

func NewDelete(TableName string) *Delete {
	return &Delete{
		Delete: NewDeleteQuery(TableName),
		Where:  NewWhereQuery(),
		Other:  New(""),
	}
}

type Delete struct {
	Delete *DeleteQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (d *Delete) Query() *PlainQuery {
	return Concat(d.Delete, d.Where, d.Other)
}

type Insert struct {
	Insert *InsertQuery
	Other  *PlainQuery
}

func (i *Insert) Query() *PlainQuery {
	return Concat(i.Insert, i.Other)
}

type Update struct {
	Update *UpdateQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (u *Update) Query() *PlainQuery {
	return Concat(u.Update, u.Where, u.Other)
}
