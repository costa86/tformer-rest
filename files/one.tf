
variable "name_one" {
  type    = string
  default = "one more"
}

variable "name_count_one" {
  type    = number
  default = 11
}

resource "random_pet" "name_one" {
  prefix = var.name_one
  length = var.name_count_one
}

output "name_one" {
  value = random_pet.name_one
}