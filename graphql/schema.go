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

	talhaoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Talhao",
		Fields: graphql.Fields{
			"id":                     &graphql.Field{Type: graphql.String},
			"nome":                   &graphql.Field{Type: graphql.String},
			"areaHectares":           &graphql.Field{Type: graphql.Float},
			"cultura":                &graphql.Field{Type: graphql.String},
			"variedade":              &graphql.Field{Type: graphql.String},
			"coordenadasGeograficas": &graphql.Field{Type: graphql.String},
			"createdAt":              &graphql.Field{Type: dateTimeType},
			"updatedAt":              &graphql.Field{Type: dateTimeType},
		},
	})

	equipamentoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Equipamento",
		Fields: graphql.Fields{
			"id":               &graphql.Field{Type: graphql.String},
			"nome":             &graphql.Field{Type: graphql.String},
			"modelo":           &graphql.Field{Type: graphql.String},
			"fabricante":       &graphql.Field{Type: graphql.String},
			"anoFabricacao":    &graphql.Field{Type: graphql.Int},
			"larguraBarra":     &graphql.Field{Type: graphql.Float},
			"capacidadeTanque": &graphql.Field{Type: graphql.Float},
			"numeroSerie":      &graphql.Field{Type: graphql.String},
			"createdAt":        &graphql.Field{Type: dateTimeType},
			"updatedAt":        &graphql.Field{Type: dateTimeType},
		},
	})

	tipoAplicacaoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "TipoAplicacao",
		Fields: graphql.Fields{
			"id":            &graphql.Field{Type: graphql.String},
			"nome":          &graphql.Field{Type: graphql.String},
			"descricao":     &graphql.Field{Type: graphql.String},
			"vazaoPadrao":   &graphql.Field{Type: graphql.Float},
			"tipoProduto":   &graphql.Field{Type: graphql.String},
			"unidadeMedida": &graphql.Field{Type: graphql.String},
			"createdAt":     &graphql.Field{Type: dateTimeType},
			"updatedAt":     &graphql.Field{Type: dateTimeType},
		},
	})

	aplicacaoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Aplicacao",
		Fields: graphql.Fields{
			"id":                &graphql.Field{Type: graphql.String},
			"talhaoId":          &graphql.Field{Type: graphql.String},
			"equipamentoId":     &graphql.Field{Type: graphql.String},
			"tipoAplicacaoId":   &graphql.Field{Type: graphql.String},
			"dataInicio":        &graphql.Field{Type: dateTimeType},
			"dataFim":           &graphql.Field{Type: dateTimeType},
			"dosagem":           &graphql.Field{Type: graphql.Float},
			"volumeAplicado":    &graphql.Field{Type: graphql.Float},
			"operador":          &graphql.Field{Type: graphql.String},
			"condicaoClimatica": &graphql.Field{Type: graphql.String},
			"observacoes":       &graphql.Field{Type: graphql.String},
			"finalizada":        &graphql.Field{Type: graphql.Boolean},
			"createdAt":         &graphql.Field{Type: dateTimeType},
			"updatedAt":         &graphql.Field{Type: dateTimeType},
			"talhao":            &graphql.Field{Type: talhaoType},
			"equipamento":       &graphql.Field{Type: equipamentoType},
			"tipoAplicacao":     &graphql.Field{Type: tipoAplicacaoType},
		},
	})

	aplicacaoInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "AplicacaoInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"talhaoId":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"equipamentoId":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"tipoAplicacaoId":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"dataInicio":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(dateTimeType)},
			"dosagem":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
			"volumeAplicado":    &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"operador":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"condicaoClimatica": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"observacoes":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	aplicacaoUpdateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "AplicacaoUpdateInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"talhaoId":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"equipamentoId":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"tipoAplicacaoId":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"dataInicio":        &graphql.InputObjectFieldConfig{Type: dateTimeType},
			"dataFim":           &graphql.InputObjectFieldConfig{Type: dateTimeType},
			"dosagem":           &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"volumeAplicado":    &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"operador":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"condicaoClimatica": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"observacoes":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"finalizada":        &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		},
	})
)

func CreateSchema(db *mongo.Database) (graphql.Schema, error) {
	resolver := &Resolver{db: db}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"aplicacao": &graphql.Field{
				Type: aplicacaoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.GetAplicacao,
			},
			"aplicacoes": &graphql.Field{
				Type: graphql.NewList(aplicacaoType),
				Args: graphql.FieldConfigArgument{
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 10},
					"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
				},
				Resolve: resolver.GetAplicacoes,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createAplicacao": &graphql.Field{
				Type: aplicacaoType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(aplicacaoInputType)},
				},
				Resolve: resolver.CreateAplicacao,
			},
			"updateAplicacao": &graphql.Field{
				Type: aplicacaoType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(aplicacaoUpdateInputType)},
				},
				Resolve: resolver.UpdateAplicacao,
			},
			"deleteAplicacao": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.DeleteAplicacao,
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

func (r *Resolver) GetAplicacao(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := r.db.Collection("aplicacoes")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "talhoes"},
			{Key: "localField", Value: "talhao_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "talhao"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "equipamentos"},
			{Key: "localField", Value: "equipamento_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "equipamento"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "tipos_aplicacao"},
			{Key: "localField", Value: "tipo_aplicacao_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tipo_aplicacao"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$talhao"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$equipamento"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$tipo_aplicacao"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var aplicacao models.Aplicacao
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&aplicacao); err != nil {
			return nil, err
		}
	}

	return aplicacao, nil
}

func (r *Resolver) GetAplicacoes(p graphql.ResolveParams) (interface{}, error) {
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	collection := r.db.Collection("aplicacoes")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "talhoes"},
			{Key: "localField", Value: "talhao_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "talhao"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "equipamentos"},
			{Key: "localField", Value: "equipamento_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "equipamento"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "tipos_aplicacao"},
			{Key: "localField", Value: "tipo_aplicacao_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tipo_aplicacao"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$talhao"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$equipamento"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$tipo_aplicacao"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$skip", Value: int64(offset)}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var aplicacoes []models.Aplicacao
	if err = cursor.All(context.TODO(), &aplicacoes); err != nil {
		return nil, err
	}

	return aplicacoes, nil
}

func (r *Resolver) CreateAplicacao(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})

	talhaoID, _ := primitive.ObjectIDFromHex(input["talhaoId"].(string))
	equipamentoID, _ := primitive.ObjectIDFromHex(input["equipamentoId"].(string))
	tipoAplicacaoID, _ := primitive.ObjectIDFromHex(input["tipoAplicacaoId"].(string))

	aplicacao := models.Aplicacao{
		TalhaoID:          talhaoID,
		EquipamentoID:     equipamentoID,
		TipoAplicacaoID:   tipoAplicacaoID,
		DataInicio:        input["dataInicio"].(time.Time),
		Dosagem:           input["dosagem"].(float64),
		Operador:          input["operador"].(string),
		CondicaoClimatica: input["condicaoClimatica"].(string),
		Observacoes:       input["observacoes"].(string),
		Finalizada:        false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if volumeAplicado, ok := input["volumeAplicado"].(float64); ok {
		aplicacao.VolumeAplicado = &volumeAplicado
	}

	collection := r.db.Collection("aplicacoes")
	result, err := collection.InsertOne(context.TODO(), aplicacao)
	if err != nil {
		return nil, err
	}

	aplicacao.ID = result.InsertedID.(primitive.ObjectID)
	return aplicacao, nil
}

func (r *Resolver) UpdateAplicacao(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(string)
	input := p.Args["input"].(map[string]interface{})

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if talhaoID, ok := input["talhaoId"].(string); ok {
		oid, _ := primitive.ObjectIDFromHex(talhaoID)
		update["$set"].(bson.M)["talhao_id"] = oid
	}
	if equipamentoID, ok := input["equipamentoId"].(string); ok {
		oid, _ := primitive.ObjectIDFromHex(equipamentoID)
		update["$set"].(bson.M)["equipamento_id"] = oid
	}
	if tipoAplicacaoID, ok := input["tipoAplicacaoId"].(string); ok {
		oid, _ := primitive.ObjectIDFromHex(tipoAplicacaoID)
		update["$set"].(bson.M)["tipo_aplicacao_id"] = oid
	}
	if dataInicio, ok := input["dataInicio"].(time.Time); ok {
		update["$set"].(bson.M)["data_inicio"] = dataInicio
	}
	if dataFim, ok := input["dataFim"].(time.Time); ok {
		update["$set"].(bson.M)["data_fim"] = dataFim
	}
	if dosagem, ok := input["dosagem"].(float64); ok {
		update["$set"].(bson.M)["dosagem"] = dosagem
	}
	if volumeAplicado, ok := input["volumeAplicado"].(float64); ok {
		update["$set"].(bson.M)["volume_aplicado"] = volumeAplicado
	}
	if operador, ok := input["operador"].(string); ok {
		update["$set"].(bson.M)["operador"] = operador
	}
	if condicaoClimatica, ok := input["condicaoClimatica"].(string); ok {
		update["$set"].(bson.M)["condicao_climatica"] = condicaoClimatica
	}
	if observacoes, ok := input["observacoes"].(string); ok {
		update["$set"].(bson.M)["observacoes"] = observacoes
	}
	if finalizada, ok := input["finalizada"].(bool); ok {
		update["$set"].(bson.M)["finalizada"] = finalizada
		if finalizada {
			update["$set"].(bson.M)["data_fim"] = time.Now()
		}
	}

	collection := r.db.Collection("aplicacoes")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	// Retornar aplicação atualizada
	return r.GetAplicacao(graphql.ResolveParams{Args: map[string]interface{}{"id": id}})
}

func (r *Resolver) DeleteAplicacao(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	collection := r.db.Collection("aplicacoes")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
