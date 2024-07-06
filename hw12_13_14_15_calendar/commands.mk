

goose-down:
	goose -dir migrations postgres "user=otus_user password=otus_password dbname=calendar sslmode=disable port=5432" down

goose-up:
	goose -dir migrations postgres "user=otus_user password=otus_password dbname=calendar sslmode=disable port=5432" up

run-calendar:
	go run ./cmd/calendar/*.go --config=./configs/calendar_config.yaml

run-scheduler:
	go run ./cmd/scheduler/*.go --config=./configs/scheduler_config.yaml

run-sender:
	go run ./cmd/sender/*.go --config=./configs/sender_config.yaml

connect-to-grpc-server:
	grpcui -plaintext localhost:50051

