package domain

type ItemStaticData struct {
	ID       string
	NumCode  string
	CharCode string
	Name     string
}

func NewItemStaticData(id, numCode, charCode string, name string) *ItemStaticData {
	return &ItemStaticData{
		ID:       id,
		NumCode:  numCode,
		CharCode: charCode,
		Name:     name,
	}
}
