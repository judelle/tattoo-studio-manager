package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		fmt.Println("0 — Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
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
			clients = deleteClientChoice(scanner, clients)
		case "6":
			newAppointment, ok := addAppointment(scanner, clients, appointments)
			if !ok {
				fmt.Println("Запись не добавлена")
				continue
			}
			appointments = append(appointments, newAppointment)
			fmt.Println("Запись успешно добавлена")

		case "0":
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
