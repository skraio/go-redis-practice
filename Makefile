include .test.env

.PHONY: run
run:
	@echo 'Запуск сервера API...'
	go run ./cmd/api
