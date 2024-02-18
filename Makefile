create-cert: 
	echo "Creating SSL certificate"
	bash postgres.sh
clean-cert:
	echo "Cleaning SSL certificate"
	rm -f server.crt server.key privkey.pem server.req 
migrate-create:
	echo "Creating migration wymj_db"
	migrate create -ext sql -dir pkg/databases/migrations -seq wymj_db
migrate-up:
	echo "Migrating-up database"
	migrate -source file://pkg/databases/migrations -database 'postgres://wymj:123456@localhost:4444/wymj_db_test?sslmode=disable' -verbose up
migrate-down:
	echo "Migrating-down database"
	migrate -source file://pkg/databases/migrations -database 'postgres://wymj:123456@localhost:4444/wymj_db_test?sslmode=disable' -verbose down
