up:
	docker-compose -f deployment/docker-compose.yml up -d --build

clean:
	docker-compose -f deployment/docker-compose.yml down
	docker volume rm deployment_econetic-db