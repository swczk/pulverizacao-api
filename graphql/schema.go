package graphql

import (
	"context"
	"time"

	"pulverizacao-api/models"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dateTimeType = graphql.NewScalar(graphql.ScalarConfig{
		Name: "DateTime",
		Serialize: func(value interface{}) interface{} {
			switch t := value.(type) {
			case time.Time:
				return t.Format(time.RFC3339)
			case *time.Time:
				if t != nil {
					return t.Format(time.RFC3339)
				}
				return nil
			}
			return nil
		},
		ParseValue: func(value interface{}) interface{} {
			switch str := value.(type) {
			case string:
				t, _ := time.Parse(time.RFC3339, str)
				return t
			}
			return nil
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				t, _ := time.Parse(time.RFC3339, valueAST.Value)
				return t
			}
			return nil
		},
	})

	geoPointType = graphql.NewObject(graphql.ObjectConfig{
		Name: "GeoPoint",
		Fields: graphql.Fields{
			"latitude":  &graphql.Field{Type: graphql.Float},
			"longitude": &graphql.Field{Type: graphql.Float},
			"timestamp": &graphql.Field{Type: dateTimeType},
			"altitude":  &graphql.Field{Type: graphql.Float},
			"speed":     &graphql.Field{Type: graphql.Float},
			"accuracy":  &graphql.Field{Type: graphql.Float},
		},
	})

	geoTrajetoriaType = graphql.NewObject(graphql.ObjectConfig{
		Name: "GeoTrajetoria",
		Fields: graphql.Fields{
			"aplicacaoId":         &graphql.Field{Type: graphql.String},
			"pontoInicial":        &graphql.Field{Type: geoPointType},
			"pontoFinal":          &graphql.Field{Type: geoPointType},
			"trajetoria":          &graphql.Field{Type: graphql.NewList(geoPointType)},
			"areaCobertura":       &graphql.Field{Type: graphql.Float},
			"distanciaPercorrida": &graphql.Field{Type: graphql.Float},
			"createdAt":           &graphql.Field{Type: dateTimeType},
			"updatedAt":           &graphql.Field{Type: dateTimeType},
		},
	})

	geoPointInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "GeoPointInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"latitude":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
			"longitude": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
			"timestamp": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(dateTimeType)},
			"altitude":  &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"speed":     &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"accuracy":  &graphql.InputObjectFieldConfig{Type: graphql.Float},
		},
	})

	geoTrajetoriaInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "GeoTrajetoriaInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"aplicacaoId":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"pontoInicial":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(geoPointInputType)},
			"pontoFinal":          &graphql.InputObjectFieldConfig{Type: geoPointInputType},
			"trajetoria":          &graphql.InputObjectFieldConfig{Type: graphql.NewList(geoPointInputType)},
			"areaCobertura":       &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"distanciaPercorrida": &graphql.InputObjectFieldConfig{Type: graphql.Float},
		},
	})

	geoTrajetoriaUpdateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "GeoTrajetoriaUpdateInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"pontoFinal":          &graphql.InputObjectFieldConfig{Type: geoPointInputType},
			"novosPontos":         &graphql.InputObjectFieldConfig{Type: graphql.NewList(geoPointInputType)},
			"areaCobertura":       &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"distanciaPercorrida": &graphql.InputObjectFieldConfig{Type: graphql.Float},
		},
	})
)

func CreateSchema(db *mongo.Database) (graphql.Schema, error) {
	resolver := &Resolver{db: db}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"geoTrajetoria": &graphql.Field{
				Type: geoTrajetoriaType,
				Args: graphql.FieldConfigArgument{
					"aplicacaoId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.GetGeoTrajetoria,
			},
			"geoTrajetorias": &graphql.Field{
				Type: graphql.NewList(geoTrajetoriaType),
				Args: graphql.FieldConfigArgument{
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 10},
					"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
				},
				Resolve: resolver.GetGeoTrajetorias,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createGeoTrajetoria": &graphql.Field{
				Type: geoTrajetoriaType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(geoTrajetoriaInputType)},
				},
				Resolve: resolver.CreateGeoTrajetoria,
			},
			"updateGeoTrajetoria": &graphql.Field{
				Type: geoTrajetoriaType,
				Args: graphql.FieldConfigArgument{
					"aplicacaoId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"input":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(geoTrajetoriaUpdateInputType)},
				},
				Resolve: resolver.UpdateGeoTrajetoria,
			},
			"deleteGeoTrajetoria": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"aplicacaoId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.DeleteGeoTrajetoria,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

type Resolver struct {
	db *mongo.Database
}

func (r *Resolver) GetGeoTrajetoria(p graphql.ResolveParams) (interface{}, error) {
	aplicacaoIdStr := p.Args["aplicacaoId"].(string)
	aplicacaoID, err := primitive.ObjectIDFromHex(aplicacaoIdStr)
	if err != nil {
		return nil, err
	}

	collection := r.db.Collection("geo_trajetorias")
	filter := bson.M{"aplicacao_id": aplicacaoID}

	var geoTrajetoria models.GeoTrajetoria
	err = collection.FindOne(context.TODO(), filter).Decode(&geoTrajetoria)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return geoTrajetoria, nil
}

func (r *Resolver) GetGeoTrajetorias(p graphql.ResolveParams) (interface{}, error) {
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	collection := r.db.Collection("geo_trajetorias")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$skip", Value: int64(offset)}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var geoTrajetorias []models.GeoTrajetoria
	if err = cursor.All(context.TODO(), &geoTrajetorias); err != nil {
		return nil, err
	}

	return geoTrajetorias, nil
}

func (r *Resolver) CreateGeoTrajetoria(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})

	aplicacaoID, err := primitive.ObjectIDFromHex(input["aplicacaoId"].(string))
	if err != nil {
		return nil, err
	}

	// Processar ponto inicial
	pontoInicialInput := input["pontoInicial"].(map[string]interface{})
	pontoInicial := models.GeoPoint{
		Latitude:  pontoInicialInput["latitude"].(float64),
		Longitude: pontoInicialInput["longitude"].(float64),
		Timestamp: pontoInicialInput["timestamp"].(time.Time),
	}

	if altitude, ok := pontoInicialInput["altitude"].(float64); ok {
		pontoInicial.Altitude = &altitude
	}
	if speed, ok := pontoInicialInput["speed"].(float64); ok {
		pontoInicial.Speed = &speed
	}
	if accuracy, ok := pontoInicialInput["accuracy"].(float64); ok {
		pontoInicial.Accuracy = &accuracy
	}

	geoTrajetoria := models.GeoTrajetoria{
		AplicacaoID:  aplicacaoID,
		PontoInicial: pontoInicial,
		Trajetoria:   []models.GeoPoint{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Processar ponto final se fornecido
	if pontoFinalInput, ok := input["pontoFinal"].(map[string]interface{}); ok {
		pontoFinal := models.GeoPoint{
			Latitude:  pontoFinalInput["latitude"].(float64),
			Longitude: pontoFinalInput["longitude"].(float64),
			Timestamp: pontoFinalInput["timestamp"].(time.Time),
		}
		if altitude, ok := pontoFinalInput["altitude"].(float64); ok {
			pontoFinal.Altitude = &altitude
		}
		if speed, ok := pontoFinalInput["speed"].(float64); ok {
			pontoFinal.Speed = &speed
		}
		if accuracy, ok := pontoFinalInput["accuracy"].(float64); ok {
			pontoFinal.Accuracy = &accuracy
		}
		geoTrajetoria.PontoFinal = &pontoFinal
	}

	// Processar trajetória se fornecida
	if trajetoriaInput, ok := input["trajetoria"].([]interface{}); ok {
		for _, pontoInterface := range trajetoriaInput {
			pontoInput := pontoInterface.(map[string]interface{})
			ponto := models.GeoPoint{
				Latitude:  pontoInput["latitude"].(float64),
				Longitude: pontoInput["longitude"].(float64),
				Timestamp: pontoInput["timestamp"].(time.Time),
			}
			if altitude, ok := pontoInput["altitude"].(float64); ok {
				ponto.Altitude = &altitude
			}
			if speed, ok := pontoInput["speed"].(float64); ok {
				ponto.Speed = &speed
			}
			if accuracy, ok := pontoInput["accuracy"].(float64); ok {
				ponto.Accuracy = &accuracy
			}
			geoTrajetoria.Trajetoria = append(geoTrajetoria.Trajetoria, ponto)
		}
	}

	// Área de cobertura e distância percorrida
	if areaCobertura, ok := input["areaCobertura"].(float64); ok {
		geoTrajetoria.AreaCobertura = areaCobertura
	}
	if distanciaPercorrida, ok := input["distanciaPercorrida"].(float64); ok {
		geoTrajetoria.DistanciaPercorrida = distanciaPercorrida
	}

	collection := r.db.Collection("geo_trajetorias")
	_, err = collection.InsertOne(context.TODO(), geoTrajetoria)
	if err != nil {
		return nil, err
	}

	return geoTrajetoria, nil
}

func (r *Resolver) UpdateGeoTrajetoria(p graphql.ResolveParams) (interface{}, error) {
	aplicacaoIdStr := p.Args["aplicacaoId"].(string)
	input := p.Args["input"].(map[string]interface{})

	aplicacaoID, err := primitive.ObjectIDFromHex(aplicacaoIdStr)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	// Atualizar ponto final se fornecido
	if pontoFinalInput, ok := input["pontoFinal"].(map[string]interface{}); ok {
		pontoFinal := models.GeoPoint{
			Latitude:  pontoFinalInput["latitude"].(float64),
			Longitude: pontoFinalInput["longitude"].(float64),
			Timestamp: pontoFinalInput["timestamp"].(time.Time),
		}
		if altitude, ok := pontoFinalInput["altitude"].(float64); ok {
			pontoFinal.Altitude = &altitude
		}
		if speed, ok := pontoFinalInput["speed"].(float64); ok {
			pontoFinal.Speed = &speed
		}
		if accuracy, ok := pontoFinalInput["accuracy"].(float64); ok {
			pontoFinal.Accuracy = &accuracy
		}
		update["$set"].(bson.M)["ponto_final"] = pontoFinal
	}

	// Adicionar novos pontos à trajetória se fornecidos
	if novosPontosInput, ok := input["novosPontos"].([]interface{}); ok {
		var novosPontos []models.GeoPoint
		for _, pontoInterface := range novosPontosInput {
			pontoInput := pontoInterface.(map[string]interface{})
			ponto := models.GeoPoint{
				Latitude:  pontoInput["latitude"].(float64),
				Longitude: pontoInput["longitude"].(float64),
				Timestamp: pontoInput["timestamp"].(time.Time),
			}
			if altitude, ok := pontoInput["altitude"].(float64); ok {
				ponto.Altitude = &altitude
			}
			if speed, ok := pontoInput["speed"].(float64); ok {
				ponto.Speed = &speed
			}
			if accuracy, ok := pontoInput["accuracy"].(float64); ok {
				ponto.Accuracy = &accuracy
			}
			novosPontos = append(novosPontos, ponto)
		}
		update["$push"] = bson.M{"trajetoria": bson.M{"$each": novosPontos}}
	}

	// Atualizar área de cobertura se fornecida
	if areaCobertura, ok := input["areaCobertura"].(float64); ok {
		update["$set"].(bson.M)["area_cobertura"] = areaCobertura
	}

	// Atualizar distância percorrida se fornecida
	if distanciaPercorrida, ok := input["distanciaPercorrida"].(float64); ok {
		update["$set"].(bson.M)["distancia_percorrida"] = distanciaPercorrida
	}

	collection := r.db.Collection("geo_trajetorias")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"aplicacao_id": aplicacaoID}, update)
	if err != nil {
		return nil, err
	}

	// Retornar geo trajetória atualizada
	return r.GetGeoTrajetoria(graphql.ResolveParams{Args: map[string]interface{}{"aplicacaoId": aplicacaoIdStr}})
}

func (r *Resolver) DeleteGeoTrajetoria(p graphql.ResolveParams) (interface{}, error) {
	aplicacaoIdStr := p.Args["aplicacaoId"].(string)
	aplicacaoID, err := primitive.ObjectIDFromHex(aplicacaoIdStr)
	if err != nil {
		return false, err
	}

	collection := r.db.Collection("geo_trajetorias")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"aplicacao_id": aplicacaoID})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
