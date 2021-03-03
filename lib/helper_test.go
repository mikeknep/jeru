package lib

var (
	im = "terraform import resource.a id"
	mv = "terraform state mv module.a module.b"
	rm = "terraform state rm resource.a"

	approve   = func() (bool, error) { return true, nil }
	unapprove = func() (bool, error) { return false, nil }
)
