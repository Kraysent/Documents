REGISTRY = `terraform -chdir=../infra/terraform output -raw document-registry-id`

build-docker:
	docker build --network=host -t documents-frontend -f Dockerfile .
	docker tag documents-frontend cr.yandex/$(REGISTRY)/documents-frontend

push-docker:
	docker push cr.yandex/$(REGISTRY)/documents-frontend
