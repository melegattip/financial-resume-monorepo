package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// inMemorySavingsRepo is a simple in-memory implementation of ports.SavingsGoalRepository for tests
type inMemorySavingsRepo struct {
	goals        map[string]*domain.SavingsGoal            // key: goal.ID
	byUser       map[string]map[string]*domain.SavingsGoal // userID -> goalID -> goal
	transactions map[string][]*domain.SavingsTransaction   // goalID -> txs
}

func newInMemorySavingsRepo() *inMemorySavingsRepo {
	return &inMemorySavingsRepo{
		goals:        map[string]*domain.SavingsGoal{},
		byUser:       map[string]map[string]*domain.SavingsGoal{},
		transactions: map[string][]*domain.SavingsTransaction{},
	}
}

func (r *inMemorySavingsRepo) Create(ctx context.Context, goal *domain.SavingsGoal) error {
	r.goals[goal.ID] = goal
	if _, ok := r.byUser[goal.UserID]; !ok {
		r.byUser[goal.UserID] = map[string]*domain.SavingsGoal{}
	}
	r.byUser[goal.UserID][goal.ID] = goal
	return nil
}

func (r *inMemorySavingsRepo) GetByID(ctx context.Context, userID, goalID string) (*domain.SavingsGoal, error) {
	userGoals := r.byUser[userID]
	if userGoals == nil {
		return nil, nil
	}
	if g, ok := userGoals[goalID]; ok {
		return g, nil
	}
	return nil, nil
}

func (r *inMemorySavingsRepo) List(ctx context.Context, userID string) ([]*domain.SavingsGoal, error) {
	userGoals := r.byUser[userID]
	var out []*domain.SavingsGoal
	for _, g := range userGoals {
		out = append(out, g)
	}
	return out, nil
}

func (r *inMemorySavingsRepo) ListByStatus(ctx context.Context, userID string, status domain.SavingsGoalStatus) ([]*domain.SavingsGoal, error) {
	all, _ := r.List(ctx, userID)
	var out []*domain.SavingsGoal
	for _, g := range all {
		if g.Status == status {
			out = append(out, g)
		}
	}
	return out, nil
}

func (r *inMemorySavingsRepo) ListByCategory(ctx context.Context, userID string, category domain.SavingsGoalCategory) ([]*domain.SavingsGoal, error) {
	all, _ := r.List(ctx, userID)
	var out []*domain.SavingsGoal
	for _, g := range all {
		if g.Category == category {
			out = append(out, g)
		}
	}
	return out, nil
}

func (r *inMemorySavingsRepo) Update(ctx context.Context, goal *domain.SavingsGoal) error {
	r.goals[goal.ID] = goal
	if _, ok := r.byUser[goal.UserID]; !ok {
		r.byUser[goal.UserID] = map[string]*domain.SavingsGoal{}
	}
	r.byUser[goal.UserID][goal.ID] = goal
	return nil
}

func (r *inMemorySavingsRepo) Delete(ctx context.Context, userID, goalID string) error {
	if r.byUser[userID] != nil {
		delete(r.byUser[userID], goalID)
	}
	delete(r.goals, goalID)
	delete(r.transactions, goalID)
	return nil
}

func (r *inMemorySavingsRepo) CreateTransaction(ctx context.Context, t *domain.SavingsTransaction) error {
	r.transactions[t.GoalID] = append(r.transactions[t.GoalID], t)
	return nil
}

func (r *inMemorySavingsRepo) GetTransactionsByGoal(ctx context.Context, userID, goalID string) ([]*domain.SavingsTransaction, error) {
	return r.transactions[goalID], nil
}

func (r *inMemorySavingsRepo) GetTransactionsByUser(ctx context.Context, userID string) ([]*domain.SavingsTransaction, error) {
	var out []*domain.SavingsTransaction
	for _, txs := range r.transactions {
		out = append(out, txs...)
	}
	return out, nil
}

// noopNotification implements ports.SavingsGoalNotificationService without side effects
type noopNotification struct{}

func (n *noopNotification) NotifyGoalAchieved(ctx context.Context, goal *domain.SavingsGoal) error {
	return nil
}
func (n *noopNotification) NotifyGoalOverdue(ctx context.Context, goal *domain.SavingsGoal) error {
	return nil
}
func (n *noopNotification) NotifyMilestoneReached(ctx context.Context, goal *domain.SavingsGoal, milestone float64) error {
	return nil
}
func (n *noopNotification) NotifyAutoSaveExecuted(ctx context.Context, goal *domain.SavingsGoal, amount float64) error {
	return nil
}

func TestSavingsGoalsEndToEndScenarios(t *testing.T) {
	repo := newInMemorySavingsRepo()
	uc := NewSavingsGoalUseCase(repo, &noopNotification{})
	ctx := context.Background()

	userID := "user_test_1"

	// a) Crear una meta
	createReq := ports.CreateSavingsGoalRequest{
		UserID:       userID,
		Name:         "Viaje",
		Description:  "",
		TargetAmount: 50000,
		Category:     domain.SavingsGoalCategoryVacation,
		Priority:     domain.SavingsGoalPriorityMedium,
		TargetDate:   time.Now().AddDate(0, 1, 0),
		IsAutoSave:   false,
	}
	created, err := uc.CreateGoal(ctx, createReq)
	if err != nil {
		t.Fatalf("create goal failed: %v", err)
	}

	goalID := created.ID

	// b) cambiar solo el icono (image_url)
	icon := "data:text/plain;charset=utf-8,%F0%9F%9A%80"
	_, err = uc.UpdateGoal(ctx, ports.UpdateSavingsGoalRequest{UserID: userID, GoalID: goalID, ImageURL: &icon})
	if err != nil {
		t.Fatalf("update icon failed: %v", err)
	}

	// c) cambiar solo el nombre
	newName := "Viaje 2025"
	_, err = uc.UpdateGoal(ctx, ports.UpdateSavingsGoalRequest{UserID: userID, GoalID: goalID, Name: &newName})
	if err != nil {
		t.Fatalf("update name failed: %v", err)
	}

	// d) cambiar solo el monto objetivo
	newTarget := 60000.0
	_, err = uc.UpdateGoal(ctx, ports.UpdateSavingsGoalRequest{UserID: userID, GoalID: goalID, TargetAmount: &newTarget})
	if err != nil {
		t.Fatalf("update target amount failed: %v", err)
	}

	// e) cambiar solo la fecha
	newDate := time.Now().AddDate(0, 2, 0)
	_, err = uc.UpdateGoal(ctx, ports.UpdateSavingsGoalRequest{UserID: userID, GoalID: goalID, TargetDate: &newDate})
	if err != nil {
		t.Fatalf("update target date failed: %v", err)
	}

	// f) depósito y validar GET actualizado
	depResp, err := uc.AddSavings(ctx, ports.AddSavingsRequest{UserID: userID, GoalID: goalID, Amount: 133.0, Description: "dep 1"})
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	got, err := uc.GetGoal(ctx, ports.GetSavingsGoalRequest{UserID: userID, GoalID: goalID})
	if err != nil {
		t.Fatalf("get after deposit failed: %v", err)
	}
	if got.CurrentAmount != depResp.NewCurrentAmount {
		t.Fatalf("current amount mismatch after deposit: got %.2f want %.2f", got.CurrentAmount, depResp.NewCurrentAmount)
	}

	// g) retiro y validar actualizado
	_, err = uc.WithdrawSavings(ctx, ports.WithdrawSavingsRequest{UserID: userID, GoalID: goalID, Amount: 50.0, Description: "ret 1"})
	if err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}
	got2, err := uc.GetGoal(ctx, ports.GetSavingsGoalRequest{UserID: userID, GoalID: goalID})
	if err != nil {
		t.Fatalf("get after withdraw failed: %v", err)
	}
	expected := got.CurrentAmount - 50.0
	if absFloat(got2.CurrentAmount-expected) > 0.0001 {
		t.Fatalf("current amount mismatch after withdraw: got %.2f want %.2f", got2.CurrentAmount, expected)
	}

	// h) eliminar meta
	if err := uc.DeleteGoal(ctx, ports.DeleteSavingsGoalRequest{UserID: userID, GoalID: goalID}); err != nil {
		t.Fatalf("delete goal failed: %v", err)
	}
	final, err := repo.GetByID(ctx, userID, goalID)
	if err != nil || final != nil {
		t.Fatalf("expected goal to be deleted, got: %v err: %v", final, err)
	}
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
