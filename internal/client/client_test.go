package client

import (
	"errors"
	"testing"
)

func TestFindByID(t *testing.T) {
	clients := []Client{
		{
			ID:    "1",
			Name:  "Александр",
			Phone: "+79990000001",
		},
		{
			ID:    "2",
			Name:  "Анна",
			Phone: "+79990000002",
		},
	}

	tests := []struct {
		name      string
		id        string
		wantFound bool
		wantID    string
	}{
		{
			name:      "client found",
			id:        "1",
			wantFound: true,
			wantID:    "1",
		},
		{
			name:      "client not found",
			id:        "999",
			wantFound: false,
		},
		{
			name:      "empty ID",
			id:        "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := FindByID(clients, tt.id)

			if found != tt.wantFound {
				t.Fatalf(
					"FindByID() found = %v, want %v",
					found,
					tt.wantFound,
				)
			}

			if tt.wantFound && got.ID != tt.wantID {
				t.Errorf(
					"FindByID() ID = %q, want %q",
					got.ID,
					tt.wantID,
				)
			}
		})
	}
}

func TestFindIndexByID(t *testing.T) {
	clients := []Client{
		{ID: "1", Name: "Александр"},
		{ID: "2", Name: "Анна"},
	}

	tests := []struct {
		name      string
		id        string
		wantIndex int
		wantFound bool
	}{
		{
			name:      "first client",
			id:        "1",
			wantIndex: 0,
			wantFound: true,
		},
		{
			name:      "second client",
			id:        "2",
			wantIndex: 1,
			wantFound: true,
		},
		{
			name:      "client not found",
			id:        "999",
			wantIndex: -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, found := FindIndexByID(clients, tt.id)

			if found != tt.wantFound {
				t.Fatalf(
					"FindIndexByID() found = %v, want %v",
					found,
					tt.wantFound,
				)
			}

			if index != tt.wantIndex {
				t.Errorf(
					"FindIndexByID() index = %d, want %d",
					index,
					tt.wantIndex,
				)
			}
		})
	}
}

func TestFindByPhone(t *testing.T) {
	clients := []Client{
		{
			ID:    "1",
			Name:  "Александр",
			Phone: "+79990000001",
		},
		{
			ID:    "2",
			Name:  "Анна",
			Phone: "+79990000002",
		},
	}

	tests := []struct {
		name      string
		phone     string
		wantFound bool
		wantID    string
	}{
		{
			name:      "phone found",
			phone:     "+79990000002",
			wantFound: true,
			wantID:    "2",
		},
		{
			name:      "phone not found",
			phone:     "+79999999999",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := FindByPhone(clients, tt.phone)

			if found != tt.wantFound {
				t.Fatalf(
					"FindByPhone() found = %v, want %v",
					found,
					tt.wantFound,
				)
			}

			if tt.wantFound && got.ID != tt.wantID {
				t.Errorf(
					"FindByPhone() ID = %q, want %q",
					got.ID,
					tt.wantID,
				)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	existingClients := []Client{
		{
			ID:      "existing",
			Name:    "Анна",
			IsAdult: true,
			Phone:   "+79990000001",
			Social:  "@anna",
		},
	}

	tests := []struct {
		name    string
		input   CreateInput
		wantErr error
	}{
		{
			name: "valid client",
			input: CreateInput{
				Name:    "Александр",
				IsAdult: true,
				Phone:   "+79990000002",
				Social:  "@alex",
			},
		},
		{
			name: "empty name",
			input: CreateInput{
				Name:  "",
				Phone: "+79990000002",
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "empty phone",
			input: CreateInput{
				Name:  "Александр",
				Phone: "",
			},
			wantErr: ErrPhoneRequired,
		},
		{
			name: "phone already exists",
			input: CreateInput{
				Name:  "Другой клиент",
				Phone: "+79990000001",
			},
			wantErr: ErrPhoneAlreadyExists,
		},
		{
			name: "empty social allowed",
			input: CreateInput{
				Name:    "Александр",
				IsAdult: false,
				Phone:   "+79990000003",
				Social:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(existingClients, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Create() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if tt.wantErr != nil {
				return
			}

			if got.ID == "" {
				t.Error("Create() generated empty ID")
			}

			if got.Name != tt.input.Name {
				t.Errorf(
					"Create() Name = %q, want %q",
					got.Name,
					tt.input.Name,
				)
			}

			if got.Phone != tt.input.Phone {
				t.Errorf(
					"Create() Phone = %q, want %q",
					got.Phone,
					tt.input.Phone,
				)
			}

			if got.IsAdult != tt.input.IsAdult {
				t.Errorf(
					"Create() IsAdult = %v, want %v",
					got.IsAdult,
					tt.input.IsAdult,
				)
			}

			if got.Social != tt.input.Social {
				t.Errorf(
					"Create() Social = %q, want %q",
					got.Social,
					tt.input.Social,
				)
			}
		})
	}
}

func TestCreatePhoneAlreadyExistsReturnsExistingClient(t *testing.T) {
	clients := []Client{
		{
			ID:    "existing",
			Name:  "Анна",
			Phone: "+79990000001",
		},
	}

	got, err := Create(clients, CreateInput{
		Name:  "Другой клиент",
		Phone: "+79990000001",
	})

	if !errors.Is(err, ErrPhoneAlreadyExists) {
		t.Fatalf(
			"Create() error = %v, want %v",
			err,
			ErrPhoneAlreadyExists,
		)
	}

	if got.ID != "existing" {
		t.Errorf(
			"Create() returned client ID = %q, want %q",
			got.ID,
			"existing",
		)
	}
}

func TestChangeName(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		newName  string
		wantErr  error
		wantName string
	}{
		{
			name:     "change name",
			id:       "1",
			newName:  "Новое имя",
			wantName: "Новое имя",
		},
		{
			name:     "empty name",
			id:       "1",
			newName:  "",
			wantErr:  ErrNameRequired,
			wantName: "Старое имя",
		},
		{
			name:     "client not found",
			id:       "999",
			newName:  "Новое имя",
			wantErr:  ErrClientNotFound,
			wantName: "Старое имя",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []Client{
				{
					ID:   "1",
					Name: "Старое имя",
				},
			}

			err := ChangeName(clients, tt.id, tt.newName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"ChangeName() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if clients[0].Name != tt.wantName {
				t.Errorf(
					"ChangeName() Name = %q, want %q",
					clients[0].Name,
					tt.wantName,
				)
			}
		})
	}
}

func TestChangePhone(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		newPhone  string
		wantErr   error
		wantPhone string
	}{
		{
			name:      "change phone",
			id:        "1",
			newPhone:  "+79990000003",
			wantPhone: "+79990000003",
		},
		{
			name:      "same phone allowed",
			id:        "1",
			newPhone:  "+79990000001",
			wantPhone: "+79990000001",
		},
		{
			name:      "empty phone",
			id:        "1",
			newPhone:  "",
			wantErr:   ErrPhoneRequired,
			wantPhone: "+79990000001",
		},
		{
			name:      "phone belongs to another client",
			id:        "1",
			newPhone:  "+79990000002",
			wantErr:   ErrPhoneAlreadyExists,
			wantPhone: "+79990000001",
		},
		{
			name:      "client not found",
			id:        "999",
			newPhone:  "+79990000003",
			wantErr:   ErrClientNotFound,
			wantPhone: "+79990000001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []Client{
				{
					ID:    "1",
					Name:  "Первый",
					Phone: "+79990000001",
				},
				{
					ID:    "2",
					Name:  "Второй",
					Phone: "+79990000002",
				},
			}

			err := ChangePhone(clients, tt.id, tt.newPhone)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"ChangePhone() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if clients[0].Phone != tt.wantPhone {
				t.Errorf(
					"ChangePhone() Phone = %q, want %q",
					clients[0].Phone,
					tt.wantPhone,
				)
			}
		})
	}
}

func TestChangeSocial(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		newSocial  string
		wantErr    error
		wantSocial string
	}{
		{
			name:       "change social",
			id:         "1",
			newSocial:  "@new",
			wantSocial: "@new",
		},
		{
			name:       "empty social allowed",
			id:         "1",
			newSocial:  "",
			wantSocial: "",
		},
		{
			name:       "client not found",
			id:         "999",
			newSocial:  "@new",
			wantErr:    ErrClientNotFound,
			wantSocial: "@old",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []Client{
				{
					ID:     "1",
					Social: "@old",
				},
			}

			err := ChangeSocial(clients, tt.id, tt.newSocial)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"ChangeSocial() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if clients[0].Social != tt.wantSocial {
				t.Errorf(
					"ChangeSocial() Social = %q, want %q",
					clients[0].Social,
					tt.wantSocial,
				)
			}
		})
	}
}

func TestChangeAdultStatus(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		newStatus   bool
		wantErr     error
		wantIsAdult bool
	}{
		{
			name:        "change to adult",
			id:          "1",
			newStatus:   true,
			wantIsAdult: true,
		},
		{
			name:        "change to minor",
			id:          "1",
			newStatus:   false,
			wantIsAdult: false,
		},
		{
			name:        "client not found",
			id:          "999",
			newStatus:   true,
			wantErr:     ErrClientNotFound,
			wantIsAdult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []Client{
				{
					ID:      "1",
					IsAdult: false,
				},
			}

			err := ChangeAdultStatus(
				clients,
				tt.id,
				tt.newStatus,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"ChangeAdultStatus() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if clients[0].IsAdult != tt.wantIsAdult {
				t.Errorf(
					"ChangeAdultStatus() IsAdult = %v, want %v",
					clients[0].IsAdult,
					tt.wantIsAdult,
				)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Run("delete existing client", func(t *testing.T) {
		clients := []Client{
			{ID: "1", Name: "Первый"},
			{ID: "2", Name: "Второй"},
			{ID: "3", Name: "Третий"},
		}

		got, err := Delete(clients, "2")
		if err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf(
				"Delete() len = %d, want 2",
				len(got),
			)
		}

		if got[0].ID != "1" {
			t.Errorf(
				"Delete() first ID = %q, want %q",
				got[0].ID,
				"1",
			)
		}

		if got[1].ID != "3" {
			t.Errorf(
				"Delete() second ID = %q, want %q",
				got[1].ID,
				"3",
			)
		}
	})

	t.Run("client not found", func(t *testing.T) {
		clients := []Client{
			{ID: "1"},
			{ID: "2"},
		}

		got, err := Delete(clients, "999")

		if !errors.Is(err, ErrClientNotFound) {
			t.Fatalf(
				"Delete() error = %v, want %v",
				err,
				ErrClientNotFound,
			)
		}

		if len(got) != 2 {
			t.Errorf(
				"Delete() modified slice on error, len = %d, want 2",
				len(got),
			)
		}
	})

	t.Run("empty clients", func(t *testing.T) {
		clients := []Client{}

		got, err := Delete(clients, "1")

		if !errors.Is(err, ErrClientNotFound) {
			t.Fatalf(
				"Delete() error = %v, want %v",
				err,
				ErrClientNotFound,
			)
		}

		if len(got) != 0 {
			t.Errorf(
				"Delete() len = %d, want 0",
				len(got),
			)
		}
	})
}
