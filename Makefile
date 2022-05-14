.PHONY: deploy
deploy:
	gcloud functions deploy DayLogHandler --runtime go116 --trigger-http --allow-unauthenticated