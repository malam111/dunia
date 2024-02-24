package dunia

import (		
	"errors"
	"net/http"
)

type Crud interface {
	Insert() (error)
	Update() (error)
	Delete() (error)
}


func Init() (*http.ServeMux, error) {
	sermux := http.NewServeMux()
	if sermux == nil {
		return sermux, errors.New("Serve Mux initialization failed")
	}
	return sermux, nil
}


