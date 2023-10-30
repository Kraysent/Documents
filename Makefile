REGISTRY = `terraform -chdir=infra output -raw document-registry-id`

all:
	go build .

build-docker:
	docker build --network=host -t documents -f Dockerfile .
	docker tag documents cr.yandex/$(REGISTRY)/documents

frontend-build-docker:
	cd frontend && make build-docker

push-docker:
	docker push cr.yandex/$(REGISTRY)/documents

frontend-push-docker:
	cd frontend && make push-docker

deploy:
	cd infra && ./deploy.sh

full-deploy: build-docker frontend-build-docker push-docker frontend-push-docker deploy
