# Go SQL Builder

Simple go sql builder made in a way it doesn't hide anything from you.

## How to use it

I started embedding sql queries from .sql files in my golang project but sometimes I need to dynamically add some
filtering or sorting to my queries. It's much simpler approach compared to orms
beacuse its just string concatination.

1. First define your base query something like `SELECT * FROM orders`. This is the part that I have in .sql file
2. Then parse the users request and generate query args
3. Use query args to generate SQL
4. Execute the query

See [simple api example](./examples/simples_api_endpoint.go), it's real code from one of my projects and it works great.

## Next steps

I am not sure whats next, I'll keep using this library and as I find some things that could be added I will add them.
If you have some suggestions please open as issue, or fork and implement a feature yourself
