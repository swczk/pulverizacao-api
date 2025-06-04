// MongoDB initialization script
db = db.getSiblingDB('pulverizacao');

// Create collections
db.createCollection('aplicacoes');
db.createCollection('talhoes');
db.createCollection('equipamentos');
db.createCollection('tipos_aplicacao');

// Create indexes for better performance
db.aplicacoes.createIndex({ "talhao_id": 1 });
db.aplicacoes.createIndex({ "equipamento_id": 1 });
db.aplicacoes.createIndex({ "tipo_aplicacao_id": 1 });
db.aplicacoes.createIndex({ "data_inicio": 1 });
db.aplicacoes.createIndex({ "finalizada": 1 });
db.aplicacoes.createIndex({ "created_at": 1 });

db.talhoes.createIndex({ "nome": 1 });
db.talhoes.createIndex({ "cultura": 1 });

db.equipamentos.createIndex({ "nome": 1 });
db.equipamentos.createIndex({ "fabricante": 1 });
db.equipamentos.createIndex({ "modelo": 1 });

db.tipos_aplicacao.createIndex({ "nome": 1 });
db.tipos_aplicacao.createIndex({ "tipo_produto": 1 });

// Insert sample data
const talhaoId = new ObjectId();
const equipamentoId = new ObjectId();
const tipoAplicacaoId = new ObjectId();

// Sample Talhao
db.talhoes.insertOne({
  _id: talhaoId,
  nome: "Talhão Norte",
  area_hectares: 25.5,
  cultura: "Soja",
  variedade: "Intacta RR2 PRO",
  coordenadas_geograficas: "-25.2456,-51.3658",
  created_at: new Date(),
  updated_at: new Date()
});

// Sample Equipamento
db.equipamentos.insertOne({
  _id: equipamentoId,
  nome: "Pulverizador Autopropelido",
  modelo: "Uniport 3030",
  fabricante: "Jacto",
  ano_fabricacao: 2023,
  largura_barra: 30.0,
  capacidade_tanque: 3000.0,
  numero_serie: "JAC2023001",
  created_at: new Date(),
  updated_at: new Date()
});

// Sample Tipo Aplicacao
db.tipos_aplicacao.insertOne({
  _id: tipoAplicacaoId,
  nome: "Herbicida Pós-Emergente",
  descricao: "Aplicação de herbicida após a emergência da cultura",
  vazao_padrao: 150.0,
  tipo_produto: "Herbicida",
  unidade_medida: "L/ha",
  created_at: new Date(),
  updated_at: new Date()
});

// Sample Aplicacao
db.aplicacoes.insertOne({
  talhao_id: talhaoId,
  equipamento_id: equipamentoId,
  tipo_aplicacao_id: tipoAplicacaoId,
  data_inicio: new Date("2025-05-28T08:30:00Z"),
  dosagem: 2.5,
  volume_aplicado: 150.0,
  operador: "João Silva",
  condicao_climatica: "Ensolarado, vento fraco",
  observacoes: "Aplicação de rotina",
  finalizada: false,
  created_at: new Date(),
  updated_at: new Date()
});

print('Database initialized with sample data');