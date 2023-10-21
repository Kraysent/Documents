SERVER_IP=51.250.78.192
scp docker-compose-prod.yaml kraysent@$SERVER_IP:~/docker-compose.yaml
scp ../configs/prod.yaml kraysent@$SERVER_IP:~/configs/prod.yaml
scp .env kraysent@$SERVER_IP:~/infra/.env
