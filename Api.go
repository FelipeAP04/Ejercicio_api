package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Incidente struct {
	ID              int    `json:"id" bson:"id"`
	Empleado        string `json:"empleado" bson:"empleado"`
	TipoEquipo      string `json:"tipo_equipo" bson:"tipo_equipo"`
	DetalleProblema string `json:"detalle_problema" bson:"detalle_problema"`
	DiaProblema     string `json:"dia_problema" bson:"dia_problema"`
}

var coll *mongo.Collection

func connectToMongoDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://Lipito:EneR3-18@apiincidentes.5zad9yj.mongodb.net/?retryWrites=true&w=majority&appName=ApiIncidentes").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	go func() {
		defer func() {
			if err = client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	}()

	// Test connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Set collection
	coll = client.Database("db_incidentes").Collection("incidentes")
}

func getIncidentes(c *gin.Context) {
	var incidentes []Incidente
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	if err = cursor.All(context.TODO(), &incidentes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}
	c.JSON(http.StatusOK, incidentes)
}

func postIncidente(c *gin.Context) {
	var newIncidente Incidente
	if err := c.BindJSON(&newIncidente); err == nil {
		newIncidente.ID = int(time.Now().Unix())
		_, err := coll.InsertOne(context.TODO(), newIncidente)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
			return
		}
		c.JSON(http.StatusCreated, newIncidente)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
	}
}

func getIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var incidente Incidente
	err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&incidente)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incidente not found"})
		return
	}
	c.JSON(http.StatusOK, incidente)
}

func putIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedIncidente Incidente
	if err := c.BindJSON(&updatedIncidente); err == nil {
		updatedIncidente.ID = id
		_, err := coll.UpdateOne(context.TODO(), bson.D{{"id", id}}, bson.D{{"$set", updatedIncidente}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
			return
		}
		c.JSON(http.StatusOK, updatedIncidente)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
	}
}

func deleteIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := coll.DeleteOne(context.TODO(), bson.D{{"id", id}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Incidente deleted"})
}

func main() {
	connectToMongoDB()

	router := gin.Default()
	router.GET("/incidentes", getIncidentes)
	router.POST("/incidentes", postIncidente)
	router.GET("/incidentes/:id", getIncidenteById)
	router.PUT("/incidentes/:id", putIncidenteById)
	router.DELETE("/incidentes/:id", deleteIncidenteById)
	router.Run("localhost:8080")
}
