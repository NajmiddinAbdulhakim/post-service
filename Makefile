CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"
migrate-gen-up:
	migrate -source file:migrations -database 'postgres://najmiddin:1234@localhost:5432/userdb?sslmode=disable' up
migrate-gen-down:
	migrate -source file:migrations -database 'postgres://najmiddin:1234@localhost:5432/userdb?sslmode=disable' down
migrate-gen-create:
	migrate create -ext sql -dir migrations -seq create_first_table