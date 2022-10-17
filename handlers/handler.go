package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"recipe-microservice/models"
	"time"
)

type RecipeHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

//function to create a new RecipeHandler struct

func NewRecipeHandler(ctx context.Context, collection *mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// endpoint methods for each operation about on a recipe instance

func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	//create a list of recipes and bind recipes in cur to that slice
	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.IndentedJSON(http.StatusOK, recipes)
	return
}

func (handler *RecipeHandler) CreateRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.BindJSON(&recipe); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	result, err := handler.collection.InsertOne(handler.ctx, recipe)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusCreated, result)
}

func (handler *RecipeHandler) GetRecipeByIDHandler(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", id}}

	var recipe models.Recipe
	result := handler.collection.FindOne(handler.ctx, filter)
	if err := result.Decode(&recipe); err != nil {
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"error": "Sorry Resource Not found",
			})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	searchKey := c.Param("key")

	filter := bson.M{"$or": bson.M{
		{"tags", searchKey},
		{"ingredients", searchKey},
		{"instructions", searchKey},
	}}

	cursor, err := handler.collection.Find(handler.ctx, filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer cursor.Close(handler.ctx)

	recipeResult := make([]models.Recipe, 0)
	for cursor.Next(handler.ctx) {
		var recipe models.Recipe
		cursor.Decode(&recipe)
		recipeResult = append(recipeResult, recipe)
	}
	if len(recipeResult) == 0 {
		c.IndentedJSON(http.StatusNotFound,
			gin.H{"error": "Sorry no document contains your search query"})
		return
	}
	c.IndentedJSON(http.StatusOK, recipeResult)
}
