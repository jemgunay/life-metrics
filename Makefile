.PHONY: lint-vue
lint-vue:
	cd ui && npm run lint

.PHONY: lint-vue-fix
lint-vue-fix:
	cd ui && npm run lint-fix

.PHONY: deploy-vue
deploy-vue:
	cd ui && npm run build

.PHONY: deploy-go
deploy-go:
	gcloud builds submit --tag gcr.io/life-metrics-316018/life-metrics
	gcloud run deploy life-metrics --image gcr.io/life-metrics-316018/life-metrics --platform managed

.PHONY: deploy
deploy: deploy-vue deploy-go

.PHONY: g-tail-logs
g-tail-logs:
	gcloud app logs tail -s default