.PHONY: lint-vue
lint-vue:
	cd ui && npm run lint

.PHONY: lint-vue-fix
lint-vue-fix:
	cd ui && npm run lint-fix

.PHONY: deploy
deploy:
	cd ui && npm run build
	rm -rf build && mkdir -p build/ui && cp -r ui/dist/* build/ui
	gcloud app deploy --verbosity debug

.PHONY: tail-logs
tail-logs:
	gcloud app logs tail -s default