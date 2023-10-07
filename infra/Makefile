init:
	terraform init

plan:
	TF_VAR_password=${YC_MONGO_PASSWORD} TF_VAR_token=${YC_TOKEN} TF_VAR_cloud_id=${YC_CLOUD_ID} TF_VAR_folder_id=${YC_FOLDER_ID} terraform plan

apply:
	TF_VAR_password=${YC_MONGO_PASSWORD} TF_VAR_token=${YC_TOKEN} TF_VAR_cloud_id=${YC_CLOUD_ID} TF_VAR_folder_id=${YC_FOLDER_ID} terraform apply
