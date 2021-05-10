.PHONY: lint-vue
lint-vue:
	cd ui && npm run lint

.PHONY: lint-vue-fix
lint-vue-fix:
	cd ui && npm run lint-fix