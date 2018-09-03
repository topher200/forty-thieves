.PHONY: run
run:
	cd webcmd && go run main.go

.PHONY: create-db
create-db:
	pgmgr db create

.PHONY: migrate-db
migrate-db:
	pgmgr db migrate
