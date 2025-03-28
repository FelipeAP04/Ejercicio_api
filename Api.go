package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Incidente struct {
	ID              int    `json:"id"`
	Empleado        string `json:"empleado"`
	TipoEquipo      string `json:"tipo_equipo"`
	DetalleProblema string `json:"detalle_problema"`
	DiaProblema     string `json:"dia_problema"`
}

var posts = []Incidente{
	{
		ID:              1,
		Empleado:        "Juan Pérez",
		TipoEquipo:      "Computadora",
		DetalleProblema: "No enciende la pantalla",
		DiaProblema:     "2023-10-15",
	},
	{
		ID:              2,
		Empleado:        "María Gómez",
		TipoEquipo:      "Impresora",
		DetalleProblema: "Atascamiento de papel frecuente",
		DiaProblema:     "2023-10-16",
	},
	{
		ID:              3,
		Empleado:        "Carlos Ruiz",
		TipoEquipo:      "Red",
		DetalleProblema: "Conexión WiFi intermitente",
		DiaProblema:     "2023-10-17",
	},
	{
		ID:              4,
		Empleado:        "Ana López",
		TipoEquipo:      "Computadora",
		DetalleProblema: "Teclado no responde",
		DiaProblema:     "2023-10-18",
	},
	{
		ID:              5,
		Empleado:        "Pedro Sánchez",
		TipoEquipo:      "Teléfono",
		DetalleProblema: "No tiene tono de llamada",
		DiaProblema:     "2023-10-19",
	},
}

func getIncidentes(c *gin.Context) {
	c.JSON(http.StatusOK, posts)
}

func main() {
	router := gin.Default()
	router.GET("/incidentes", getIncidentes)
	router.Run("localhost:8080")
}
