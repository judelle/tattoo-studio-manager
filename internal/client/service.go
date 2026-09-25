package client

func FindByID(clients []Client, id string) (Client, bool) {
	for _, client := range clients {
		if client.ID == id {
			return client, true
		}
	}

	return Client{}, false
}

func FindIndexByID(clients []Client, id string) (int, bool) {
	for i, client := range clients {
		if client.ID == id {
			return i, true
		}
	}

	return -1, false
}

func FindByPhone(clients []Client, phone string) (Client, bool) {
	for _, client := range clients {
		if client.Phone == phone {
			return client, true
		}
	}
	return Client{}, false
}
