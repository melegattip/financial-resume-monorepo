package repository

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB implementa un mock de la base de datos para pruebas
type MockDB struct {
	mock.Mock
	db *gorm.DB
}

// NewMockDB crea una nueva instancia de MockDB
func NewMockDB() *MockDB {
	return &MockDB{
		db: &gorm.DB{},
	}
}

// Create simula la creación de un registro
func (m *MockDB) Create(value interface{}) *gorm.DB {
	m.Called(value)
	return m.db
}

// First simula la búsqueda del primer registro
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	m.Called(out, where)
	return m.db
}

// Find simula la búsqueda de registros
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	m.Called(out, where)
	return m.db
}

// Save simula el guardado de un registro
func (m *MockDB) Save(value interface{}) *gorm.DB {
	m.Called(value)
	return m.db
}

// Delete simula la eliminación de registros
func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	m.Called(value, where)
	return m.db
}

// Error retorna el error simulado
func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

// RowsAffected retorna el número de filas afectadas simulado
func (m *MockDB) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

// Begin simula el inicio de una transacción
func (m *MockDB) Begin() *gorm.DB {
	m.Called()
	return m.db
}

// Commit simula el commit de una transacción
func (m *MockDB) Commit() *gorm.DB {
	m.Called()
	return m.db
}

// Where simula la cláusula WHERE
func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return m.db
}

// Model simula la selección del modelo
func (m *MockDB) Model(value interface{}) *gorm.DB {
	m.Called(value)
	return m.db
}

// Rollback simula el rollback de una transacción
func (m *MockDB) Rollback() *gorm.DB {
	m.Called()
	return m.db
}
