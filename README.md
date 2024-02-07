## Test
podman run -p 4444:5432 --name wymjdb -e POSTGRES_PASSWORD=12456 -e POSTGRES_USER=wymj -e POSTGRES_DB=wymj_db_test -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres
