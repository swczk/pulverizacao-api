package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aplicacao struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TalhaoID          primitive.ObjectID `json:"talhaoId" bson:"talhao_id"`
	EquipamentoID     primitive.ObjectID `json:"equipamentoId" bson:"equipamento_id"`
	TipoAplicacaoID   primitive.ObjectID `json:"tipoAplicacaoId" bson:"tipo_aplicacao_id"`
	DataInicio        time.Time          `json:"dataInicio" bson:"data_inicio"`
	DataFim           *time.Time         `json:"dataFim,omitempty" bson:"data_fim,omitempty"`
	Dosagem           float64            `json:"dosagem" bson:"dosagem"`
	VolumeAplicado    *float64           `json:"volumeAplicado,omitempty" bson:"volume_aplicado,omitempty"`
	Operador          string             `json:"operador" bson:"operador"`
	CondicaoClimatica string             `json:"condicaoClimatica" bson:"condicao_climatica"`
	Observacoes       string             `json:"observacoes" bson:"observacoes"`
	Finalizada        bool               `json:"finalizada" bson:"finalizada"`
	CreatedAt         time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updated_at"`

	// Campos populados via lookup
	Talhao        *Talhao        `json:"talhao,omitempty" bson:"talhao,omitempty"`
	Equipamento   *Equipamento   `json:"equipamento,omitempty" bson:"equipamento,omitempty"`
	TipoAplicacao *TipoAplicacao `json:"tipoAplicacao,omitempty" bson:"tipo_aplicacao,omitempty"`
}

type AplicacaoInput struct {
	TalhaoID          primitive.ObjectID `json:"talhaoId"`
	EquipamentoID     primitive.ObjectID `json:"equipamentoId"`
	TipoAplicacaoID   primitive.ObjectID `json:"tipoAplicacaoId"`
	DataInicio        time.Time          `json:"dataInicio"`
	Dosagem           float64            `json:"dosagem"`
	VolumeAplicado    *float64           `json:"volumeAplicado,omitempty"`
	Operador          string             `json:"operador"`
	CondicaoClimatica string             `json:"condicaoClimatica"`
	Observacoes       string             `json:"observacoes"`
}

type AplicacaoUpdateInput struct {
	TalhaoID          *primitive.ObjectID `json:"talhaoId,omitempty"`
	EquipamentoID     *primitive.ObjectID `json:"equipamentoId,omitempty"`
	TipoAplicacaoID   *primitive.ObjectID `json:"tipoAplicacaoId,omitempty"`
	DataInicio        *time.Time          `json:"dataInicio,omitempty"`
	DataFim           *time.Time          `json:"dataFim,omitempty"`
	Dosagem           *float64            `json:"dosagem,omitempty"`
	VolumeAplicado    *float64            `json:"volumeAplicado,omitempty"`
	Operador          *string             `json:"operador,omitempty"`
	CondicaoClimatica *string             `json:"condicaoClimatica,omitempty"`
	Observacoes       *string             `json:"observacoes,omitempty"`
	Finalizada        *bool               `json:"finalizada,omitempty"`
}

type Talhao struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nome                   string             `json:"nome" bson:"nome"`
	AreaHectares           float64            `json:"areaHectares" bson:"area_hectares"`
	Cultura                string             `json:"cultura" bson:"cultura"`
	Variedade              string             `json:"variedade" bson:"variedade"`
	CoordenadasGeograficas string             `json:"coordenadasGeograficas" bson:"coordenadas_geograficas"`
	CreatedAt              time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt              time.Time          `json:"updatedAt" bson:"updated_at"`
}

type Equipamento struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nome             string             `json:"nome" bson:"nome"`
	Modelo           string             `json:"modelo" bson:"modelo"`
	Fabricante       string             `json:"fabricante" bson:"fabricante"`
	AnoFabricacao    int                `json:"anoFabricacao" bson:"ano_fabricacao"`
	LarguraBarra     float64            `json:"larguraBarra" bson:"largura_barra"`
	CapacidadeTanque float64            `json:"capacidadeTanque" bson:"capacidade_tanque"`
	NumeroSerie      string             `json:"numeroSerie" bson:"numero_serie"`
	CreatedAt        time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updated_at"`
}

type TipoAplicacao struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nome          string             `json:"nome" bson:"nome"`
	Descricao     string             `json:"descricao" bson:"descricao"`
	VazaoPadrao   float64            `json:"vazaoPadrao" bson:"vazao_padrao"`
	TipoProduto   string             `json:"tipoProduto" bson:"tipo_produto"`
	UnidadeMedida string             `json:"unidadeMedida" bson:"unidade_medida"`
	CreatedAt     time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updated_at"`
}
