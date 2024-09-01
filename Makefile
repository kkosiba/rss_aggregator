run:
	docker compose up \
		--force-recreate \
		--build

cleanup:
	docker compose down \
		--remove-orphans \
		--rmi local \
		--volumes
