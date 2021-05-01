resource "random_pet" "good_boy" {}

resource "random_string" "main" {
  count = 9

  length = 3 * (count.index + 1)
}

resource "random_password" "letters" {
  count = 3

  length  = 3 * (count.index + 1)
  lower   = true
  upper   = true
  number  = false
  special = false
}

resource "random_password" "numbers" {
  count = 3

  length  = 3 * (count.index + 1)
  lower   = false
  upper   = false
  number  = true
  special = false
}

resource "random_password" "special" {
  count = 3

  length  = 3 * (count.index + 1)
  lower   = false
  upper   = false
  number  = false
  special = true
}
