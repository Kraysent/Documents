terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }

  required_version = ">= 0.13"
}

provider "yandex" {
  zone      = "ru-central1-a"
  token     = var.token
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
}

resource "yandex_mdb_mongodb_cluster" "documents" {
  name        = "documents"
  environment = "PRODUCTION"
  network_id  = "enpaghvmmpik3acvm435"

  cluster_config {
    version = "4.4"
  }

  database {
    name = "documents"
  }

  user {
    name     = "superuser"
    password = var.password
    permission {
      database_name = "documents"
      roles         = ["mdbDbAdmin"]
    }
  }

  resources_mongod {
    resource_preset_id = "b1.medium"
    disk_type_id       = "network-ssd"
    disk_size          = 20
  }

  host {
    zone_id          = "ru-central1-a"
    subnet_id        = "e9btpiqfmo4j3gc0eca6"
    assign_public_ip = true
  }

  labels = {
    test = "test_value"
  }
}
