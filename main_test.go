package main

import (
	"effectivemobile/FIO"
	"effectivemobile/initializers"
	"testing"
)

func TestCreateFIO(t *testing.T) {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
	fio := FIO.FIO{ID: 1, Name: "Danat", Surname: "Fednova", Patronymic: "Alekseevich",
		Age: 54, Gender: "female", Nation: "KZ"}

	err := initializers.DB.Create(&fio).Error
	if err != nil {
		t.Errorf("ФИО не удалилось")
	}
}

func TestGetFIO(t *testing.T) {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
	id := 1
	expected := FIO.FIO{ID: uint(id), Name: "Danat", Surname: "Fednova", Patronymic: "Alekseevich",
		Age: 44, Gender: "female", Nation: "KZ"}
	var fio FIO.FIO
	initializers.DB.First(&fio, id)
	if expected != fio {
		t.Errorf("Получено другое значение из бд")
	}
}

func TestGetAllFIOs(t *testing.T) {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
	minimumSize := 1
	var fios []FIO.FIO
	initializers.DB.Find(&fios)
	if len(fios) < minimumSize {
		t.Errorf("Пустой массив ФИО")
	}
}

func TestUpdateFio(t *testing.T) {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
	id := 1
	expected := FIO.FIO{ID: uint(id), Name: "Danat", Surname: "Fednova", Patronymic: "Alekseevich",
		Age: 54, Gender: "female", Nation: "KZ"}

	var newFIO FIO.FIO
	var fio FIO.FIO
	initializers.DB.First(&fio, id)
	fio.SetAge(54)
	initializers.DB.Save(&fio)

	initializers.DB.First(&newFIO, id)
	if expected != fio {
		t.Errorf("Изменение возраста не произошло(")
	}
}

func TestDeleteFio(t *testing.T) {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
	id := 1

	err := initializers.DB.Delete(&FIO.FIO{}, id).Error
	if err != nil {
		t.Errorf("ФИО не удалилось")
	}
}
