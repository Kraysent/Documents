resource "yandex_mdb_postgresql_cluster" "documents" {
  name        = "documents"
  environment = "PRODUCTION"
  network_id  = yandex_vpc_network._DOCUMENTS_NETS_.id

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
    name             = "documents-host-a"
    subnet_id        = yandex_vpc_subnet._DOCUMENTS_PG_NETS_.id
    assign_public_ip = true
  }
}

resource "yandex_vpc_network" "_DOCUMENTS_NETS_" {
  name = "_DOCUMENTS_NETS_"
}

resource "yandex_vpc_subnet" "_DOCUMENTS_PG_NETS_" {
  name           = "_DOCUMENTS_PG_NETS_"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network._DOCUMENTS_NETS_.id
  v4_cidr_blocks = ["10.5.0.0/24"]
}

resource "yandex_mdb_postgresql_user" "documents" {
  cluster_id = yandex_mdb_postgresql_cluster.documents.id
  name       = "documents"
  password   = var.password
}

resource "yandex_mdb_postgresql_database" "documentsdb" {
  cluster_id = yandex_mdb_postgresql_cluster.documents.id
  name       = "documentsdb"
  owner      = "documents"
}

output "postgres_mdb_id" {
  value = yandex_mdb_postgresql_cluster.documents.id
}
