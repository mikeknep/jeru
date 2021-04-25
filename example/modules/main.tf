resource "random_pet" "good_boy" {}


resource "random_string" "short_lower" {
  length = 5
  lower  = true
  upper  = false
}
resource "random_string" "long_lower" {
  length = 15
  lower  = true
  upper  = false
}
resource "random_string" "short_upper" {
  length = 5
  lower  = false
  upper  = true
}
resource "random_string" "long_upper" {
  length = 15
  lower  = false
  upper  = true
}


resource "random_password" "superuser" {
  length      = 32
  min_lower   = 5
  min_numeric = 5
  min_special = 5
  min_upper   = 5
}
resource "random_password" "admin" {
  length      = 32
  min_lower   = 5
  min_numeric = 0
  min_special = 0
  min_upper   = 5
}
resource "random_password" "readonly" {
  length      = 24
  min_lower   = 5
  min_numeric = 0
  min_special = 0
  min_upper   = 5
}
