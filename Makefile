migrate: 
	migrate -path=./migrations -database=$$MOVIES_DB_DSN up

rollback: 
	migrate -path=./migrations -database=$$MOVIES_DB_DSN down