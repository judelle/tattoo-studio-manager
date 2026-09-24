package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type Appointment struct {
	ID           string
	ClientID     string
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
	ID      string
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

func findClientByID(clients []Client, id string) (Client, bool) {
	for _, client := range clients {
		if client.ID == id {
			return client, true
		}
	}

	return Client{}, false
}

func findClientIndexByID(clients []Client, id string) (int, bool) {
	for i, client := range clients {
		if client.ID == id {
			return i, true
		}
	}

	return -1, false
}

func findClientByPhone(clients []Client, phone string) (Client, bool) {
	for _, client := range clients {
		if client.Phone == phone {
			return client, true
		}
	}
	return Client{}, false
}

func addClient(scanner *bufio.Scanner, clients []Client) Client {
	name := ""
	for {
		fmt.Print("Введите имя: ")
		scanner.Scan()
		if name = strings.TrimSpace(scanner.Text()); name == "" {
			fmt.Println("Имя не может быть пустым")
			continue
		}
		break
	}

	isAdult := false
	for {
		fmt.Print("Клиент совершеннолетний? Д/Н: ")
		scanner.Scan()
		adultInput := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if adultInput == "д" {
			isAdult = true
			break
		}
		if adultInput == "н" {
			break
		}
		fmt.Println("Некорректный ввод")

	}

	phone := ""
	for {
		fmt.Print("Введите номер телефона: ")
		scanner.Scan()
		if phone = strings.TrimSpace(scanner.Text()); phone == "" {
			fmt.Println("Номер телефона не может быть пустым")
			continue
		} else {
			client, isFound := findClientByPhone(clients, phone)
			if !isFound {
				break
			} else {
				fmt.Printf("Номер %s принадлежит клиенту: %s, с ID: %s\n",
					phone,
					client.Name,
					client.ID)
			}
		}
	}

	fmt.Print("Введите соц. сеть: ")
	scanner.Scan()
	social := strings.TrimSpace(scanner.Text())

	client := Client{
		ID:      uuid.NewString(),
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

	fmt.Println("-------------------")
	id := strings.TrimSpace(scanner.Text())
	client, isFound := findClientByID(clients, id)
	if !isFound {
		fmt.Printf("клиент с ID %s не найден\n", id)
		fmt.Println("-------------------")
		return
	}
	printClient(client)

}

func changeClient(scanner *bufio.Scanner, clients []Client) {
	fmt.Print("Введите ID клиента для изменения: ")
	scanner.Scan()

	fmt.Println("-------------------")
	id := strings.TrimSpace(scanner.Text())
	index, isFound := findClientIndexByID(clients, id)
	if !isFound {
		fmt.Printf("клиент с ID %s не найден\n", id)
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
			for {
				fmt.Println("Введите новое имя или 0 для отмены")
				scanner.Scan()
				newName := strings.TrimSpace(scanner.Text())
				if newName == "0" {
					break
				}
				if newName == "" {
					fmt.Println("Имя не может быть пустым")
					continue
				}
				clients[index].Name = newName
				break
			}
		case "2":
			for {
				newPhone := ""
				fmt.Print("Введите номер телефона: ")
				scanner.Scan()
				if newPhone = strings.TrimSpace(scanner.Text()); newPhone == "" {
					fmt.Println("Номер телефона не может быть пустым")
					continue
				} else {
					client, isFound := findClientByPhone(clients, newPhone)
					if !isFound || clients[index].ID == client.ID {
						clients[index].Phone = newPhone
						break
					} else {
						fmt.Printf("Номер %s принадлежит клиенту: %s, с ID: %s\n",
							newPhone,
							client.Name,
							client.ID)
					}
				}
			}
		case "3":
			fmt.Println("Введите новую соц. сеть")
			scanner.Scan()
			newSocial := strings.TrimSpace(scanner.Text())
			clients[index].Social = newSocial
		case "4":
			for {
				fmt.Print("Клиент совершеннолетний? Д/Н: ")
				scanner.Scan()
				newIsAdult := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if newIsAdult == "д" {
					clients[index].IsAdult = true
					break
				}
				if newIsAdult == "н" {
					clients[index].IsAdult = false
					break
				}
				fmt.Println("Некорректный ввод")
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
		fmt.Println("Для выхода - 0")
		fmt.Print("Введите ID клиента для удаления: ")
		scanner.Scan()
		if strings.TrimSpace(scanner.Text()) == "0" {
			return clients
		}
		fmt.Println("-------------------")
		id := strings.TrimSpace(scanner.Text())
		index, isFound := findClientIndexByID(clients, id)
		if !isFound {
			fmt.Printf("клиент с ID %s не найден\n", id)
			fmt.Println("-------------------")
			continue
		}
		printClient(clients[index])

		for {
			fmt.Println()
			fmt.Println("1 — Удалить клиента")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			scanner.Scan()
			deleteChoice := strings.TrimSpace(strings.ToLower(scanner.Text()))
			fmt.Println("-------------------")

			switch deleteChoice {
			case "1":
				fmt.Println("Вы уверены?")
				fmt.Println("1 — Да")
				fmt.Println("0 — Нет")
				scanner.Scan()
				if strings.TrimSpace(scanner.Text()) == "1" {
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
