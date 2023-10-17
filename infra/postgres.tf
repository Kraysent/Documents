resource "yandex_mdb_postgresql_cluster" "documents" {
  name = "documents"
  environment = "PRODUCTION"
  network_id  = "enpaghvmmpik3acvm435"

  config {
    version = 15
    resources {
      resource_preset_id = "b1.medium"
      disk_type_id       = "network-ssd"
      disk_size          = 20
    }
  }

  host {
    zone             = "ru-central1-a"
    subnet_id        = "e9btpiqfmo4j3gc0eca6"
    assign_public_ip = true
  }
}