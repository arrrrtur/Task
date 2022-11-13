start-project: docker-compose up --build balance

create-database: @docker-compose exec pgdb psql -U postgres -d avito_db -a -f ./data.sql