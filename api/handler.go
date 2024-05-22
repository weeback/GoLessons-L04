package api

import (
	"CoffeeStore"
	"CoffeeStore/source/db"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

var dbCtrl db.ProductDB

func New() *Handler {
	return &Handler{}
}

type Handler struct {
}

func (h *Handler) RegisterRouter(r *gin.Engine, db db.ProductDB) {
	dbCtrl = db
	r.GET("api/version", h.getAppVersion)
	r.GET("api/products", h.getProducts)
	r.GET("api/products/name", h.getProductsByName)
}

func (h *Handler) getAppVersion(c *gin.Context) {
	// wrap the handler function to avoid the need to pass gin.Context
	func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Version: %s", CoffeeStore.Version); err != nil {
			log.Printf("%s - %s - Error writing response: %v", r.RemoteAddr, r.RequestURI, err)
			return
		}
	}(c.Writer, c.Request)
}

func (h *Handler) getProducts(c *gin.Context) {
	// wrap the handler function to avoid the need to pass gin.Context
	func(w http.ResponseWriter, r *http.Request) {
		name := c.Query("title")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'title' parameter"})
			return
		}
		objectIDHex := "664daf1f16f7338b8935d8d6"

		// Convert the string to an ObjectID
		objectID, err := primitive.ObjectIDFromHex(objectIDHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing ObjectID:"})
			fmt.Printf("Error parsing ObjectID: %v\n", err)
			return
		}
		results, err := dbCtrl.GetProductById(objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "objectID": objectID})
			return
		}
		c.JSON(http.StatusOK, results)
	}(c.Writer, c.Request)
}

func (h *Handler) getProductsByName(c *gin.Context) {
	// wrap the handler function to avoid the need to pass gin.Context
	func(w http.ResponseWriter, r *http.Request) {
		name := c.Query("title")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'title' parameter"})
			return
		}

		results, err := dbCtrl.GetProductByName(context.TODO(), name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "name": name})
			return
		}
		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "name": name})
			return
		}
		c.JSON(http.StatusOK, results)
	}(c.Writer, c.Request)
}
