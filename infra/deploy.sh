SERVER_IP=docarchive.space
terraform -chdir=terraform output -raw certificate_key | ssh kraysent@$SERVER_IP "cat >~/.ssl/docarchive.space.key"
terraform -chdir=terraform output -json certificate_chain | jq -r ".[0]" | ssh kraysent@$SERVER_IP "cat >~/.ssl/docarchive.space.crt"

scp docker-compose.yaml kraysent@$SERVER_IP:~/docker-compose.yaml
scp ../configs/prod.yaml kraysent@$SERVER_IP:~/configs/prod.yaml
scp .env kraysent@$SERVER_IP:~/.env
scp start.sh kraysent@$SERVER_IP:~/start.sh
scp stop.sh kraysent@$SERVER_IP:~/stop.sh
scp nginx/nginx.conf kraysent@$SERVER_IP:~/configs/nginx.conf
scp grafana/config.river kraysent@$SERVER_IP:~/configs/config.river
ssh kraysent@$SERVER_IP "chmod +x start.sh"
ssh kraysent@$SERVER_IP "docker-compose pull"
