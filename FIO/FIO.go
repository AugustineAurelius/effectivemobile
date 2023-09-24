package FIO

import (
	enrichment2 "effectivemobile/FIO/enrichment"
)

type FIO struct {
	ID         uint `gorm:"primaryKey; autoIncrement"`
	Name       string
	Surname    string
	Patronymic string
	Age        float64
	Gender     string
	Nation     string
}

// Конструктор
func NewFIO(name string, surname string) FIO {
	return FIO{Name: name, Surname: surname,
		Age: enrichment2.AddAge(name), Gender: enrichment2.AddGender(name),
		Nation: enrichment2.AddNation(name)}
}

// Геттеры и Сеттеры
func (f FIO) GetName() string {
	return f.Name
}
func (f FIO) GetSurname() string {
	return f.Surname
}
func (f FIO) GetPatronymic() string {
	return f.Patronymic
}
func (f FIO) GetAge() float64 {
	return f.Age
}
func (f FIO) GetGender() string {
	return f.Gender
}
func (f FIO) GetNation() string {
	return f.Nation
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
func (f *FIO) SetAge(age float64) {
	f.Age = age
}
func (f *FIO) SetGender(gender string) {
	f.Gender = gender
}
func (f *FIO) SetNation(nation string) {
	f.Nation = nation
}
