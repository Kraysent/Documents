REGISTRY = `terraform -chdir=infra output -raw document-registry-id`

all:
	go build .

build-docker:
	docker build -t documents -f Dockerfile .
	docker tag documents cr.yandex/$(REGISTRY)/documents

push-docker:
	docker push cr.yandex/$(REGISTRY)/documents

deploy:
	cd infra && ./deploy.sh

full-deploy: build-docker push-docker deploy
