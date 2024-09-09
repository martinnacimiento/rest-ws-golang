package tests

import (
	"net/http/httptest"
	"testing"

	"tincho.dev/rest-ws/utils"
)

func TestGetPagination(t *testing.T) {
	tests := []struct {
		query      string
		wantOffset int64
		wantLimit  int64
		wantErr    bool
	}{
		// Caso por defecto: sin parámetros de offset y limit
		{"", 0, 10, false},
		// Caso con parámetros válidos
		{"offset=5&limit=20", 5, 20, false},
		// Caso con offset inválido
		{"offset=invalid&limit=20", 0, 0, true},
		// Caso con limit inválido
		{"offset=5&limit=invalid", 0, 0, true},
		// Caso con solo offset válido
		{"offset=8", 8, 10, false},
		// Caso con solo limit válido
		{"limit=25", 0, 25, false},
	}

	for _, tt := range tests {
		// Crear una solicitud con la query simulada
		req := httptest.NewRequest("GET", "/?"+tt.query, nil)

		// Llamar a la función GetPagination
		gotOffset, gotLimit, err := utils.GetPagination(req)

		// Verificar si se esperaba un error
		if (err != nil) != tt.wantErr {
			t.Errorf("GetPagination() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}

		// Verificar los valores de offset y limit
		if gotOffset != tt.wantOffset {
			t.Errorf("GetPagination() offset = %v, want %v", gotOffset, tt.wantOffset)
		}

		if gotLimit != tt.wantLimit {
			t.Errorf("GetPagination() limit = %v, want %v", gotLimit, tt.wantLimit)
		}
	}
}
