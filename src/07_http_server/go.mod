module http_server

go 1.20

replace http_server/handlers => ./handlers

replace http_server/middleware => ./middleware

require (
	http_server/handlers v0.0.0-00010101000000-000000000000
	http_server/middleware v0.0.0-00010101000000-000000000000
)
