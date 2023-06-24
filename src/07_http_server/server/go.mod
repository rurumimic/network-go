module server

go 1.20

require (
	github.com/awoodbeck/gnp v0.0.0-20230225045246-30fd6b8da810
	http_server/handlers v0.0.0-00010101000000-000000000000
)

replace http_server/handlers => ../handlers
