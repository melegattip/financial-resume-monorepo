package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		expectedError error
	}{
		{
			name:          "Test Category",
			userID:        "user_id",
			expectedError: nil,
		},
		{
			name:          "",
			userID:        "user_id",
			expectedError: errors.New("category name cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := NewCategoryBuilder().
				SetName(tt.name).
				SetUserID(tt.userID).
				Build()

			err := category.Validate()
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.NotNil(t, category)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, category)
			assert.Equal(t, tt.name, category.Name)
			assert.Equal(t, tt.userID, category.UserID)
		})
	}
}

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category *Category
		wantErr  bool
	}{
		{
			name: "Valid category",
			category: &Category{
				ID:   "123",
				Name: "Groceries",
			},
			wantErr: false,
		},
		{
			name: "Invalid - empty title",
			category: &Category{
				ID:   "123",
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrEmptyCategoryName, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCategoryBuilder(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		categoryName string
		description  string
	}{
		{
			name:         "Valid category",
			id:           "123",
			categoryName: "Groceries",
			description:  "Food and household items",
		},
		{
			name:         "Empty description",
			id:           "456",
			categoryName: "Transportation",
			description:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := NewCategoryBuilder().
				SetID(tt.id).
				SetName(tt.categoryName).
				Build()

			assert.Equal(t, tt.id, category.ID)
			assert.Equal(t, tt.categoryName, category.Name)
		})
	}
}

func TestCategoryBuilder_CustomDates(t *testing.T) {
	category := NewCategoryBuilder().
		SetID("789").
		SetName("Entertainment").
		Build()

	assert.Equal(t, "789", category.ID)
	assert.Equal(t, "Entertainment", category.Name)
}
