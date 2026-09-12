migrate-up:
	goose -dir migrations postgres "$(DATABASE_URI)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URI)" down
