output "instance_id_list" {
  value = "${ucloud_instance.web.*.id}"
}

output "eip_id_list" {
  value = "${ucloud_eip.default.*.id}"
}
