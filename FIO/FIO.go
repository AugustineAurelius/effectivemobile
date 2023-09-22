package FIO

import (
	enrichment2 "effectivemobile/FIO/enrichment"
)

type FIO struct {
	name       string
	surname    string
	patronymic string
	age        float64
	gender     string
	nation     string
}

func NewFIO(name string, surname string) FIO {
	return FIO{name: name, surname: surname,
		age: enrichment2.AddAge(name), gender: enrichment2.AddGender(name),
		nation: enrichment2.AddNation(name)}
}

func (f FIO) Name() string {
	return f.name
}
func (f FIO) Surname() string {
	return f.surname
}
func (f FIO) Patronymic() string {
	return f.patronymic
}

func (f *FIO) SetName(name string) {
	f.name = name
}
func (f *FIO) SetSurname(surname string) {
	f.surname = surname
}
func (f *FIO) SetPatronymic(patronymic string) {
	f.patronymic = patronymic
}
