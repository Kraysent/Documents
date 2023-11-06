SERVER_IP=docarchive.space
scp docker-compose-prod.yaml kraysent@$SERVER_IP:~/docker-compose.yaml
scp ../configs/prod.yaml kraysent@$SERVER_IP:~/configs/prod.yaml
scp .env kraysent@$SERVER_IP:~/infra/.env
scp start.sh kraysent@$SERVER_IP:~/start.sh
scp nginx/nginx.conf kraysent@$SERVER_IP:~/configs/nginx.conf
scp -r .well-known kraysent@$SERVER_IP:~/configs/.well-known
ssh kraysent@$SERVER_IP "chmod +x start.sh"
ssh kraysent@$SERVER_IP "docker-compose pull"
