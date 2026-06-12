package mealplan

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create creates a new meal plan by calling POST /api/meal-plan/
func Create(ctx context.Context, c *tandoor.Client, params MealPlanParam) (*MealPlanResponse, error) {
	res, err := tandoor.Request[MealPlanResponse](ctx, c, "POST", "/api/meal-plan/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal plan: %w", err)
	}
	return res, nil
}

// List lists meal plans by calling GET /api/meal-plan/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*MealPlanListResponse, error) {
	qb := tandoor.NewQuery().
		Add("from_date", params.FromDate).
		Add("to_date", params.ToDate).
		Add("meal_type", params.MealType).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[MealPlanListResponse](ctx, c, "GET", "/api/meal-plan/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal plans: %w", err)
	}
	return res, nil
}

// AutoPlan generates an automatic meal plan by calling POST /api/auto-plan/
func AutoPlan(ctx context.Context, c *tandoor.Client, params AutoMealPlanParam) (*AutoMealPlanResponse, error) {
	res, err := tandoor.Request[AutoMealPlanResponse](ctx, c, "POST", "/api/auto-plan/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to auto generate meal plan: %w", err)
	}
	return res, nil
}
