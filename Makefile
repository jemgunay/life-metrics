.PHONY: local
local:
	npm install --legacy-peer-deps
	cd ui && npm install dotenv-webpack --save-dev
	echo 'VUE_APP_API_HOST=http://localhost:8080' > ui/.env.local

.PHONY: lint-go
lint-go:
	golint ./... | grep -v "vendor/"

.PHONY: lint-vue
lint-vue:
	cd ui && npm run lint

.PHONY: lint-vue-fix
lint-vue-fix:
	cd ui && npm run lint-fix

.PHONY: deploy-vue
build-vue:
	cd ui && rm -rf dist && npm run build

.PHONY: deploy-go
deploy-go:
	gcloud builds submit --tag gcr.io/life-metrics-316018/life-metrics
	gcloud run deploy life-metrics --image gcr.io/life-metrics-316018/life-metrics --platform managed #--version=staging

.PHONY: deploy
deploy: build-vue deploy-go

.PHONY: g-tail-logs
g-tail-logs:
	gcloud app logs tail -s default