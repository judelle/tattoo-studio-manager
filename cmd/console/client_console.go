package main

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
	"tattoo-studio/internal/client"

	"github.com/google/uuid"
)

func printClient(client client.Client) {
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
}

func addClient(scanner *bufio.Scanner, clients []client.Client) client.Client {
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
			client, isFound := client.FindClientByPhone(clients, phone)
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

	client := client.Client{
		ID:      uuid.NewString(),
		Name:    name,
		IsAdult: isAdult,
		Phone:   phone,
		Social:  social,
	}
	return client
}

func printClients(clients []client.Client) {
	if len(clients) == 0 {
		fmt.Println("Вы еще не добавили клиентов")
		return
	}
	fmt.Printf("Список из %d клиентов: \n", len(clients))
	for _, client := range clients {
		printClient(client)
	}
}

func findClient(scanner *bufio.Scanner, clients []client.Client) {
	fmt.Print("Введите ID клиента: ")
	scanner.Scan()

	fmt.Println("-------------------")
	id := strings.TrimSpace(scanner.Text())
	client, isFound := client.FindClientByID(clients, id)
	if !isFound {
		fmt.Printf("клиент с ID %s не найден\n", id)
		fmt.Println("-------------------")
		return
	}
	printClient(client)

}

func changeClient(scanner *bufio.Scanner, clients []client.Client) {
	fmt.Print("Введите ID клиента для изменения: ")
	scanner.Scan()

	fmt.Println("-------------------")
	id := strings.TrimSpace(scanner.Text())
	index, isFound := client.FindClientIndexByID(clients, id)
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
					client, isFound := client.FindClientByPhone(clients, newPhone)
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

func deleteClientChoice(scanner *bufio.Scanner, clients []client.Client) []client.Client {
	for {
		fmt.Println("Для выхода - 0")
		fmt.Print("Введите ID клиента для удаления: ")
		scanner.Scan()
		if strings.TrimSpace(scanner.Text()) == "0" {
			return clients
		}
		fmt.Println("-------------------")
		id := strings.TrimSpace(scanner.Text())
		index, isFound := client.FindClientIndexByID(clients, id)
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
