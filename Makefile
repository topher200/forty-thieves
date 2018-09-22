.PHONY: run-webcmd
run-webcmd: install-dependencies
	cd webcmd && go run main.go

.PHONY: run-solvercmd
run-solvercmd: install-dependencies
	cd solvercmd && go install
	solvercmd

.PHONY: create-db
create-db:
	pgmgr db create

.PHONY: migrate-db
migrate-db:
	pgmgr db migrate

.PHONY: recreate-test-db
recreate-test-db:
	pgmgr --config-file .pgmgr.test.json db drop | true
	pgmgr --config-file .pgmgr.test.json db create
	pgmgr --config-file .pgmgr.test.json db migrate

.PHONY: test
test: | recreate-test-db run-tests

.PHONY: install-dependencies
install-dependencies:
	dep ensure

# helper func, you should call 'test' instead
.PHONY: run-tests
run-tests: install-dependencies
	./test.sh
