package schema

import (
	"effectivemobile/FIO"
	"effectivemobile/initializers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func GraphqlHandler(c *gin.Context) {
	schema, err := CreateSchema(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	param := struct {
		Query string `json:"query" binding:"required"`
	}{}

	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: param.Query,
	})
	if len(result.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%+v", result.Errors)})
		return
	}

	c.JSON(http.StatusOK, result.Data)
}

func CreateSchema(c *gin.Context) (graphql.Schema, error) {
	fioType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "FIO",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"surname": &graphql.Field{
					Type: graphql.String,
				},
				"patronymic": &graphql.Field{
					Type: graphql.String,
				},
				"age": &graphql.Field{
					Type: graphql.Float,
				},
				"gender": &graphql.Field{
					Type: graphql.String,
				},
				"nation": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	fields := graphql.Fields{
		"FIO": &graphql.Field{
			Type:        fioType,
			Description: "Get fio by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if !ok {
					return nil, fmt.Errorf("Invalid id")
				}

				var user FIO.FIO
				err := initializers.DB.First(&user, id).Error
				if err != nil {
					return nil, err
				}

				return user, nil
			},
		},
		"fio": &graphql.Field{
			Type:        graphql.NewList(fioType),
			Description: "Get all users",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var users []FIO.FIO
				err := initializers.DB.Find(&users).Error
				if err != nil {
					return nil, err
				}

				return users, nil
			},
		},
	}

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
		},
	)
	/*****/
	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* Create new product item
			http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
			*/
			"create": &graphql.Field{
				Type:        fioType,
				Description: "Create new product",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"patronymic": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"nation": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"gender": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rand.Seed(time.Now().UnixNano())
					fio := FIO.FIO{
						ID:         uint(rand.Intn(100000)), // generate random ID
						Name:       params.Args["name"].(string),
						Surname:    params.Args["surname"].(string),
						Patronymic: params.Args["patronymic"].(string),
						Age:        params.Args["age"].(float64),
						Gender:     params.Args["gender"].(string),
						Nation:     params.Args["nation"].(string),
					}
					initializers.DB.Create(&fio)
					return fio, nil
				},
			},

			// Update product by id

			"update": &graphql.Field{
				Type:        fioType,
				Description: "Update product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"patronymic": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"nation": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"gender": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(uint)
					name, nameOk := params.Args["name"].(string)
					surname, surnameOk := params.Args["surname"].(string)
					patronymic, patronymicOk := params.Args["patronymic"].(string)
					age, ageOk := params.Args["age"].(float64)
					gender, genderOk := params.Args["gender"].(string)
					nation, nationOk := params.Args["nation"].(string)
					var fio FIO.FIO

					err := initializers.DB.First(&fio, id).Error
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fio"})
						log.Fatal(err)
					}
					if nameOk {
						fio.SetName(name)
					}
					if surnameOk {
						fio.SetSurname(surname)
					}
					if patronymicOk {
						fio.SetPatronymic(patronymic)
					}
					if ageOk {
						fio.SetAge(age)
					}
					if genderOk {
						fio.SetGender(gender)
					}
					if nationOk {
						fio.SetNation(nation)
					}
					err = initializers.DB.Save(&fio).Error
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fio"})
						log.Fatal(err)
					}

					return fio, nil
				},
			},

			//Delete product by id

			"delete": &graphql.Field{
				Type:        fioType,
				Description: "Delete product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					var fio FIO.FIO
					err := initializers.DB.First(&fio, id).Error
					if err != nil {
						c.JSON(http.StatusNotFound, gin.H{"error": "FIO not found"})
						log.Fatal(err)
					}
					err = initializers.DB.Delete(&FIO.FIO{}, id).Error
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						log.Fatal(err)
					}
					return fio, nil
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Error :%s\n", err)
	}

	return schema, nil
}
