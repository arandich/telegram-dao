package entity

type Role struct {
	ListRoles map[int]string
}

var Roles = Role{ListRoles: map[int]string{
	1: "Участник",
	2: "Член клуба",
	3: "Администратор",
}}
