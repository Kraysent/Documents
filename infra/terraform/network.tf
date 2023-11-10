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

data "yandex_cm_certificate_content" "docarchive_cert" {
  certificate_id = "fpqorhic5ov3967opovp"
}

output "certificate_chain" {
  value = data.yandex_cm_certificate_content.docarchive_cert.certificates
}

output "certificate_key" {
  value     = data.yandex_cm_certificate_content.docarchive_cert.private_key
  sensitive = true
}
