package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func main() {
	clients := make([]client.Client, 0)
	appointments := make([]appointment.Appointment, 0)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("1 — Добавить клиента")
		fmt.Println("2 — Показать клиентов")
		fmt.Println("3 — Найти клиента по ID")
		fmt.Println("4 — Изменить клиента")
		fmt.Println("5 — Удалить клиента")
		fmt.Println("6 — Добавить запись")
		fmt.Println("7 — Найти записи по дате")
		fmt.Println("8 — Перенести запись")
		fmt.Println("9 — Завершить запись")
		fmt.Println("10 — Отменить запись")
		fmt.Println("11 — Удалить запись")
		fmt.Println("12 — История записей клиента")
		fmt.Println("13 — Ближайшие записи")
		fmt.Println("0 — Выход")
		fmt.Print("Выберите действие: ")

		choice, err := readInput(scanner)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Ввод завершён")
				return
			}
			fmt.Println("main error: read choice", err)
			return
		}
		fmt.Println("-------------------")
		switch choice {
		case "1":
			newClient := addClient(scanner, clients)
			clients = append(clients, newClient)
		case "2":
			printClients(clients)
		case "3":
			findClient(scanner, clients)
		case "4":
			changeClient(scanner, clients)
		case "5":
			clients = deleteClientChoice(scanner, clients, appointments)
		case "6":
			newAppointment, ok, err := addAppointment(scanner, clients, appointments)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("addappointment error: ", err)
				return
			}
			if !ok {
				fmt.Println("Запись не добавлена")
				continue
			}
			appointments = append(appointments, newAppointment)
			fmt.Println("Запись успешно добавлена")
		case "7":
			err := printAppointmentByDate(scanner, appointments, clients)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("printAppointmentByDate error: ", err)
				return
			}
		case "8":
			err := reschedulingAppointment(scanner, appointments)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("reschedulingAppointment error: ", err)
				return
			}
		case "9":
			err := completeAppointment(scanner, appointments)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("completeAppointment error: ", err)
				return
			}
		case "10":
			err := cancelAppointment(scanner, appointments)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("cancelAppointment error: ", err)
				return
			}
		case "11":
			appointments, err = deleteAppointmentChoice(scanner, appointments, clients)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("deleteAppointmentChoice error: ", err)
				return
			}
		case "12":
			err := printClientHistory(scanner, appointments, clients)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Ввод завершён")
					return
				}
				fmt.Println("printClientHistory error: ", err)
				return
			}
		case "13":
			printUpcoming(appointments, clients)
		case "0":
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
