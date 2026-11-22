## change to your database address
DATABASE=postgres://postgres:12345@localhost:5432/cinemanetwork?sslmode=disable

dbu: ## Migrates the database to the latest version.
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20  -verbose up

dbd: ## Migrates the database down by first version.
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20 -verbose down

dbdrop: ## Drops the database.
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20 -verbose drop -f

dbseq: ## Creates a new migration file with the given name. ex: name=create_users_table
	migrate create -ext sql -seq -dir ./migrations $(name)

dbver: ## Prints the current migration version.
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20 version

dbf: ## Force the migrates the database to the given version. ex: version=1
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20 -verbose force $(version)

dbt: ## Goto migrates the database to the given version. ex: version=1
	migrate -database="$(DATABASE)" -path=./migrations -lock-timeout=20 -verbose goto $(version)

dbreset: dbdrop dbu ## Reset database to the latest version.