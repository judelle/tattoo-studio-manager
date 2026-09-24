package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Appointment struct {
	ID           int
	ClientID     int
	Date         string
	Time         string
	Description  string
	PlannedPrice int
	FinalPrice   int
	Deposit      int
	Status       string
	SketchPaths  []string
}

type Client struct {
	ID      int
	Name    string
	IsAdult bool
	Phone   string
	Social  string
}

func main() {
	clients := make([]Client, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("1 — Добавить клиента")
		fmt.Println("2 — Показать клиентов")
		fmt.Println("3 — Найти клиента по ID")
		fmt.Println("4 — Изменить клиента")
		fmt.Println("5 — Удалить клиента")
		fmt.Println("0 — Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		fmt.Println("-------------------")
		switch choice {
		case "1":
			newClient := addClient(scanner, len(clients)+1)
			clients = append(clients, newClient)
		case "2":
			printClients(clients)
		case "3":
			findClient(scanner, clients)
		case "4":
			changeClient(scanner, clients)
		case "5":
			clients = deleteClientChoice(scanner, clients)
		case "0":
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func printClient(client Client) {
	fmt.Println("-------------------")
	fmt.Println("ID:", client.ID)
	fmt.Println("Имя:", client.Name)

	if client.IsAdult {
		fmt.Println("Совершеннолетний")
	} else {
		fmt.Println("Несовершеннолетний")
	}

	fmt.Println("Телефон:", client.Phone)
	fmt.Println("Соц. сеть:", client.Social)
	fmt.Println("-------------------")
}

func findClientByID(clients []Client, id int) (Client, error) {
	for _, client := range clients {
		if client.ID == id {
			return client, nil
		}
	}

	return Client{}, fmt.Errorf("клиент с ID %d не найден", id)
}

func findClientIndexByID(clients []Client, id int) (int, error) {
	for i, client := range clients {
		if client.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("клиент с ID %d не найден", id)
}

func addClient(scanner *bufio.Scanner, id int) Client {

	fmt.Print("Введите имя: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())

	fmt.Print("Клиент совершеннолетний? Д/Н: ")
	scanner.Scan()
	adultInput := strings.TrimSpace(strings.ToLower(scanner.Text()))
	isAdult := false
	if adultInput == "д" {
		isAdult = true
	}

	fmt.Print("Введите номер телефона: ")
	scanner.Scan()
	phone := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите соц. сеть: ")
	scanner.Scan()
	social := strings.TrimSpace(scanner.Text())

	client := Client{
		ID:      id,
		Name:    name,
		IsAdult: isAdult,
		Phone:   phone,
		Social:  social,
	}
	return client
}

func printClients(clients []Client) {
	if len(clients) == 0 {
		fmt.Println("Вы еще не добавили клиентов")
		return
	}

	for _, client := range clients {
		printClient(client)
	}
}

func findClient(scanner *bufio.Scanner, clients []Client) {
	fmt.Print("Введите ID клиента: ")
	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())

	fmt.Println("-------------------")

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("ID должен быть числом")
		return
	}

	client, err := findClientByID(clients, id)
	if err != nil {
		fmt.Println(err)
		fmt.Println("-------------------")
		return
	}
	printClient(client)

}

func changeClient(scanner *bufio.Scanner, clients []Client) {
	fmt.Print("Введите ID клиента для изменения: ")
	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())

	fmt.Println("-------------------")

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("ID должен быть числом")
		return
	}

	index, err := findClientIndexByID(clients, id)
	if err != nil {
		fmt.Println(err)
		fmt.Println("-------------------")
		return
	}
	for {
		fmt.Println()
		fmt.Println("1 — Изменить имя")
		fmt.Println("2 — Изменить номер")
		fmt.Println("3 — Изменить соц. Сети")
		fmt.Println("4 — Изменить статус совершеннолетнего")
		fmt.Println("0 — Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()
		infoToChange := strings.TrimSpace(strings.ToLower(scanner.Text()))
		fmt.Println("-------------------")

		switch infoToChange {
		case "1":
			fmt.Println("Введите новое имя")
			scanner.Scan()
			newName := strings.TrimSpace(scanner.Text())
			clients[index].Name = newName
		case "2":
			fmt.Println("Введите новый номер")
			scanner.Scan()
			newPhone := strings.TrimSpace(scanner.Text())
			clients[index].Phone = newPhone
		case "3":
			fmt.Println("Введите новую соц. сеть")
			scanner.Scan()
			newSocial := strings.TrimSpace(scanner.Text())
			clients[index].Social = newSocial
		case "4":
			fmt.Print("Клиент совершеннолетний? Д/Н: ")
			scanner.Scan()
			newIsAdult := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if newIsAdult == "д" {
				clients[index].IsAdult = true
			} else {
				clients[index].IsAdult = false
			}
		case "0":
			printClient(clients[index])
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}

}

func deleteClientChoice(scanner *bufio.Scanner, clients []Client) []Client {
	for {
		fmt.Print("Введите ID клиента для удаления: ")
		scanner.Scan()

		input := strings.TrimSpace(scanner.Text())

		fmt.Println("-------------------")

		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("ID должен быть числом")
			continue
		}

		index, err := findClientIndexByID(clients, id)
		if err != nil {
			fmt.Println(err)
			fmt.Println("-------------------")
			continue
		}

		for {
			fmt.Println()
			fmt.Println("1 — Удалить клиента")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			scanner.Scan()
			scanerDelete := strings.TrimSpace(strings.ToLower(scanner.Text()))
			fmt.Println("-------------------")

			switch scanerDelete {
			case "1":
				fmt.Println("Вы уверены?")
				fmt.Println("1 — Да")
				fmt.Println("0 — Нет")
				scanner.Scan()
				if strings.TrimSpace(strings.ToLower(scanner.Text())) == "1" {
					clients = slices.Delete(clients, index, index+1)
					fmt.Println("Клиент удалён")
					return clients
				}
			case "0":
				return clients
			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}
}
