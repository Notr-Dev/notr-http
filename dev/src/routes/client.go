package routes

type ClientData struct {
	Name string
}

func loadClientData() ClientData {
	return ClientData{
		Name: "Test Client",
	}
}
