package FIO

import (
	enrichment2 "effectivemobile/FIO/enrichment"
)

type FIO struct {
	ID         uint    `gorm:"primaryKey; autoIncrement"`
	Name       string  `gorm:"embedded"`
	Surname    string  `gorm:"embedded"`
	Patronymic string  `gorm:"embedded"`
	Age        float64 `gorm:"embedded"`
	Gender     string  `gorm:"embedded"`
	Nation     string  `gorm:"embedded"`
}

func NewFIO(name string, surname string) FIO {
	return FIO{Name: name, Surname: surname,
		Age: enrichment2.AddAge(name), Gender: enrichment2.AddGender(name),
		Nation: enrichment2.AddNation(name)}
}

func (f FIO) GetName() string {
	return f.Name
}
func (f FIO) GetSurname() string {
	return f.Surname
}
func (f FIO) GetPatronymic() string {
	return f.Patronymic
}

func (f *FIO) SetName(name string) {
	f.Name = name
}
func (f *FIO) SetSurname(surname string) {
	f.Surname = surname
}
func (f *FIO) SetPatronymic(patronymic string) {
	f.Patronymic = patronymic
}
