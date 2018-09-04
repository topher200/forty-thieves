.PHONY: run-webcmd
run-webcmd:
	cd webcmd && go run main.go

.PHONY: run-solvercmd
run-solvercmd:
	cd solvercmd && go install
	solvercmd

.PHONY: create-db
create-db:
	pgmgr db create

.PHONY: migrate-db
migrate-db:
	pgmgr db migrate

.PHONY: test
test:
	./test.sh

.PHONY: recreate-test-db
recreate-test-db:
	pgmgr --config-file .pgmgr.test.json db drop
	pgmgr --config-file .pgmgr.test.json db create
	pgmgr --config-file .pgmgr.test.json db migrate

.PHONY: test-with-fresh-database
test-with-fresh-database: | recreate-test-db test
