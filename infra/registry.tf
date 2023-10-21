resource "yandex_container_registry" "documents-registry" {
  name = "documents-registry"
  folder_id = var.folder_id
}

output "document-registry-id" {
  value = yandex_container_registry.documents-registry.id
}
