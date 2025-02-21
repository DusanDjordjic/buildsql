package examples

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	buildsql "github.com/DusanDjordjic/buildsql"
)

/*
	Lets say that we support sorting and filtering by varios parametes
	in our GetOrders endpoint. We can imagine that after parsing the request we
	can generate a GetOrdersFilters struct that will store every detail about what user wants to see.
	After that we can have a Build method that will build the QueryArgs that we can then
	use to query our database
*/

type ScheduleCompatisonType string

const (
	INVALID ScheduleCompatisonType = ""
	LT      ScheduleCompatisonType = "lt"
	LTE     ScheduleCompatisonType = "lte"
	EQ      ScheduleCompatisonType = "eq"
	GT      ScheduleCompatisonType = "gt"
	GTE     ScheduleCompatisonType = "gte"
	BETWEEN ScheduleCompatisonType = "between"
)

type OrderStatus int32

const (
	ORDER_INVALID   OrderStatus = 0
	ORDER_CREATED   OrderStatus = 1
	ORDER_REJECTED  OrderStatus = 2
	ORDER_CONFIRMED OrderStatus = 3
	ORDER_DELIVERED OrderStatus = 4
)

type Order struct {
	ID           uint64
	ScheduledFor time.Time
	Status       OrderStatus
	UserID       uint64
}

type OrderFilters struct {
	buildsql.QueryBuilder
	Status                  OrderStatus
	ScheduledStart          *time.Time
	ScheduledEnd            *time.Time
	ScheduledComparisonType ScheduleCompatisonType
}

func (filters *OrderFilters) Parse() error {
	// Do the parsing here and populate the struct
	filters.Status = ORDER_CREATED
	if scheduleStart, err := time.Parse("2006-01-02T15:06", "2025-01-20T16:30"); err == nil {
		filters.ScheduledStart = new(time.Time)
		*filters.ScheduledStart = scheduleStart
	}

	if scheduleEnd, err := time.Parse("2006-01-02T15:06", "2025-01-21T00:00"); err == nil {
		filters.ScheduledEnd = new(time.Time)
		*filters.ScheduledEnd = scheduleEnd
	}

	filters.ScheduledComparisonType = BETWEEN
	return nil
}

func (filters *OrderFilters) GenerateQueryArgs(userID uint64) buildsql.QueryArgs {
	// Reset the arg counter to 1
	filters.Reset()

	q := new(buildsql.DeafultQueryArgs)

	q.Wheres = append(q.Wheres, fmt.Sprintf("user_id = %s", filters.NextArg()))
	q.Args = append(q.Args, userID)

	if filters.Status != ORDER_INVALID {
		q.Wheres = append(q.Wheres, fmt.Sprintf("status = %s", filters.NextArg()))
		q.Args = append(q.Args, filters.Status)
	}

	if filters.ScheduledComparisonType == LT && filters.ScheduledStart != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for < %s", filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart)
	} else if filters.ScheduledComparisonType == LTE && filters.ScheduledStart != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for <= %s", filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart)
	} else if filters.ScheduledComparisonType == EQ && filters.ScheduledStart != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for = %s", filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart)
	} else if filters.ScheduledComparisonType == GT && filters.ScheduledStart != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for > %s", filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart)
	} else if filters.ScheduledComparisonType == GTE && filters.ScheduledStart != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for >= %s", filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart)
	} else if filters.ScheduledComparisonType == BETWEEN && filters.ScheduledStart != nil && filters.ScheduledEnd != nil {
		q.Wheres = append(q.Wheres, fmt.Sprintf("scheduled_for BETWEEN %s AND %s", filters.NextArg(), filters.NextArg()))
		q.Args = append(q.Args, *filters.ScheduledStart, *filters.ScheduledEnd)
	}

	// deafult soring
	q.OrderBys = append(q.OrderBys, "scheduled_for DESC")

	return q
}

// This function is called by every handler to get the query builder
// so if you want to change something you just do it here
func GetQueryBuilder() buildsql.QueryBuilder {
	return &buildsql.IntQueryBuilder{ArgCounter: 1}
}

// This is our route handler for example
func GetOrders() {
	filters := OrderFilters{}
	filters.QueryBuilder = GetQueryBuilder()

	err := filters.Parse( /* Pass some arguments if you need them */ )
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse order filters, %s", err.Error())
		// return error to client
		panic("failed to parse filters")
	}

	// Lets say that users can see only their orders
	// we will get userID from token or somewhere else and add it to query
	var userID uint64 = 10
	q := filters.GenerateQueryArgs(userID)

	var db *sql.DB = nil /* use real db connection */
	orders, err := ServiceGetOrders(db, q)
	if err != nil {
		// Internal server order
		panic("failed to get orders")
	}
	_ = orders
	// send back JSON response or render a template etc.
}

func ServiceGetOrders(tx *sql.DB, q buildsql.QueryArgs) ([]Order, error) {
	// I usualy load queries from a file on build time using go:embed but you can have them hardcoded as well

	q.SetBaseQuery("SELECT order.id, order.status, order,scheduled_for FROM orders")
	rows, err := tx.Query(q.GenerateSQL(), q.GetArgs()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		// Scan order and add it to list
		order := Order{}
		err := rows.Scan(&order.ID, &order.ScheduledFor, &order.Status)
		if err != nil {
			return nil, err
		}
	}

	return orders, err
}
