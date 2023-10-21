resource "yandex_vpc_network" "_DOCUMENTS_NETS_" {
  name = "_DOCUMENTS_NETS_"
}

resource "yandex_vpc_subnet" "_DOCUMENTS_PG_NETS_" {
  name           = "_DOCUMENTS_PG_NETS_"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network._DOCUMENTS_NETS_.id
  v4_cidr_blocks = ["10.5.0.0/24"]
}

resource "yandex_vpc_subnet" "_DOCUMENTS_VM_NETS_" {
  name           = "_DOCUMENTS_VM_NETS_"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network._DOCUMENTS_NETS_.id
  v4_cidr_blocks = ["10.4.0.0/16"]
}


