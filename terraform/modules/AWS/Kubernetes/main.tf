resource "aws_instance" "simulator_master_instances" {
  count                       = "${var.number_of_master_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${var.access_key_name}"
  instance_type               = "${var.master_instance_type}"
  vpc_security_group_ids      = ["${var.control_plane_sg_id}"]
  associate_public_ip_address = false
  subnet_id                   = "${var.private_subnet_id}"
  user_data                   = "${data.template_file.master_cloud_config.rendered}"
  iam_instance_profile        = "${var.iam_instance_profile_id}"
}

resource "aws_instance" "simulator_node_instances" {
  count                       = "${var.number_of_cluster_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${var.access_key_name}"
  instance_type               = "${var.cluster_nodes_instance_type}"
  vpc_security_group_ids      = ["${var.control_plane_sg_id}"]
  associate_public_ip_address = false
  subnet_id                   = "${var.private_subnet_id}"
  user_data                   = "${data.template_file.node_cloud_config.rendered}"
  depends_on                  = ["aws_instance.simulator_master_instances"]
  iam_instance_profile        = "${var.iam_instance_profile_id}"
}
