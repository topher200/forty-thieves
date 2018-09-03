.PHONY: create-db
create-db:
	pgmgr db create

.PHONY: migrate-db
migrate-db:
	pgmgr db migrate
