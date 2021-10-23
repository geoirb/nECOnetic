up:
	docker-compose -f deployment/docker-compose.yml up -d --build

migrate-station:
	docker run							\
		--rm   							\
		-v data-service:/data-service 	\
		-w /data-service				\
		golang:latest					\
		go run tools/migrator/stations/main.go

migrate-profile-temperature:
	docker run							\
		--rm   							\
		-v data-service:/data-service 	\
		-w /data-service				\
		golang:latest					\
		go run tools/migrator/profiler/temperature/main.go

migrate-profile-wind:
	docker run							\
		--rm   							\
		-v data-service:/data-service 	\
		-w /data-service				\
		golang:latest					\
		go run tools/migrator/profiler/wind/main.go

migrate-eco-data:
	docker run							\
		--rm   							\
		-v data-service:/data-service 	\
		-w data-service					\
		golang:latest					\
		go run tools/migrator/eco-data/main.go



