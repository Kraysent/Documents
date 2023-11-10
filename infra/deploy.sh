SERVER_IP=docarchive.space
scp docker-compose.yaml kraysent@$SERVER_IP:~/docker-compose.yaml
scp ../configs/prod.yaml kraysent@$SERVER_IP:~/configs/prod.yaml
scp .env kraysent@$SERVER_IP:~/.env
scp start.sh kraysent@$SERVER_IP:~/start.sh
scp stop.sh kraysent@$SERVER_IP:~/stop.sh
scp nginx/nginx.conf kraysent@$SERVER_IP:~/configs/nginx.conf
scp grafana/config.river kraysent@$SERVER_IP:~/configs/config.river
ssh kraysent@$SERVER_IP "chmod +x start.sh"
ssh kraysent@$SERVER_IP "docker-compose pull"
