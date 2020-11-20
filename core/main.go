package main

import (
	"go-admin/core/entitys"
	"go-admin/core/sysInit/sql"
)

func main() {
	sql.RegisterOrm(entitys.User{})
}
