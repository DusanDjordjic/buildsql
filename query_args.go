package gosqlbuilder

import (
	"fmt"
	"strings"
)

/*
Base query is something like
SELECT user.id, user.email, COUNT(order.id) FROM users user
INNER JOIN orders order ON user.id = order.user_id

Then you can dynamically add some filters (Wheres) and Sorters (Order bys)
You can also add limit and offset all with parameters

Idea is to have static queries written in a file and embed them using go:embed
Then set that text to be a base query and dynamically add wheres, order bys, limit and offset
From an API endpoint

I don't like using ORMs like gorm because its anoying to work with when you need something more complex.
Maybe I am just bad tho.
*/

type QueryArgs interface {
	SetBaseQuery(string)
	GenerateSQL() string
	GetArgs() []any
}

type DeafultQueryArgs struct {
	BaseQuery string
	Wheres    []string
	OrderBys  []string
	Limit     uint64
	Offset    uint64
	Args      []any
}

func (q *DeafultQueryArgs) GenerateSQL() string {
	out := q.BaseQuery
	if len(q.Wheres) > 0 {
		where := strings.Join(q.Wheres, " AND ")
		out += fmt.Sprintf("WHERE %s\n", where)
	}

	if len(q.OrderBys) > 0 {
		orderBy := strings.Join(q.OrderBys, ",")
		out += fmt.Sprintf("ORDER BY %s\n", orderBy)

	}

	if q.Limit > 0 {
		out += fmt.Sprintf("LIMIT %d\n", q.Limit)
	}

	if q.Offset > 0 {
		out += fmt.Sprintf("OFFSET %d", q.Offset)
	}

	return out
}

func (q *DeafultQueryArgs) SetBaseQuery(base string) {
	q.BaseQuery = base
}

func (q *DeafultQueryArgs) GetArgs() []any {
	return q.Args
}
