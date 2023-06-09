start: 
	docker compose up -d

stop:
	docker compose down

clean:
	docker compose down
	docker volume rm tanzutrends_db 

build: 
	docker compose build
	docker compose up -d --force-recreate
	