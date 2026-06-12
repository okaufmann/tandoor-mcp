package e2e_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealplan"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealtype"
	"github.com/compilercomplied/tandoor-mcp/src/tools/auto_plan"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_meal_plan"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_meal_type"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_meal_plans"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestMealPlanE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange: Create a MealType first
		mtRes, mtErr := infra.CallTool(ctx, fixture.Client, "create_meal_type", create_meal_type.Args{
			Name: "Dinner",
		})
		infra.AssertToolSuccess(t, mtRes, mtErr)
		mType := infra.ParseToolResponse[mealtype.MealTypeResponse](t, mtRes)

		recipeID := 1
		recipeName := "Tandoori Chicken"
		servings := 2.5
		fromDateStr := "2026-06-12T12:00:00Z"
		title := "Special Dinner"

		createArgs := create_meal_plan.Args{
			Title:        &title,
			RecipeID:     &recipeID,
			RecipeName:   &recipeName,
			Servings:     servings,
			FromDate:     fromDateStr,
			MealTypeID:   mType.ID,
			MealTypeName: mType.Name,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_meal_plan", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[mealplan.MealPlanResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected meal plan ID > 0, got 0")
		}
		if logEntry.Title == nil || *logEntry.Title != title {
			t.Errorf("expected Title %q, got %v", title, logEntry.Title)
		}
		if logEntry.Recipe == nil || logEntry.Recipe.ID != recipeID {
			t.Errorf("expected Recipe ID %d, got %v", recipeID, logEntry.Recipe)
		}
		if logEntry.Servings != servings {
			t.Errorf("expected Servings %v, got %v", servings, logEntry.Servings)
		}
		if logEntry.MealType.ID != mType.ID {
			t.Errorf("expected MealType ID %d, got %d", mType.ID, logEntry.MealType.ID)
		}

		// Act
		listArgs := get_meal_plans.Args{}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_meal_plans", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[mealplan.MealPlanListResponse](t, listRes)

		if logsList.Count == 0 {
			t.Errorf("expected count > 0")
		}
		found := false
		for _, entry := range logsList.Results {
			if entry.ID == logEntry.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find meal plan ID %d in get_meal_plans response", logEntry.ID)
		}
	})

	t.Run("HappyPath_AutoPlan", func(t *testing.T) {
		// Arrange: Create a MealType first
		mtRes, mtErr := infra.CallTool(ctx, fixture.Client, "create_meal_type", create_meal_type.Args{
			Name: "Lunch",
		})
		infra.AssertToolSuccess(t, mtRes, mtErr)
		mType := infra.ParseToolResponse[mealtype.MealTypeResponse](t, mtRes)

		startStr := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		endStr := time.Now().Add(72 * time.Hour).Format(time.RFC3339)

		autoArgs := auto_plan.Args{
			StartDate:   startStr,
			EndDate:     endStr,
			MealTypeID:  mType.ID,
			Servings:    2.0,
			AddShopping: false,
		}

		// Act
		autoRes, autoErr := infra.CallTool(ctx, fixture.Client, "auto_plan", autoArgs)

		// Assert
		infra.AssertToolSuccess(t, autoRes, autoErr)
		autoPlanResp := infra.ParseToolResponse[mealplan.AutoMealPlanResponse](t, autoRes)

		if autoPlanResp.MealTypeID != mType.ID {
			t.Errorf("expected meal type ID %d, got %d", mType.ID, autoPlanResp.MealTypeID)
		}
	})

	t.Run("ValidationError_MissingFromDate", func(t *testing.T) {
		// Arrange
		args := create_meal_plan.Args{
			FromDate:     "", // invalid
			MealTypeID:   1,
			MealTypeName: "Dinner",
			Servings:     2,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_meal_plan", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "from_date is required") {
			t.Errorf("expected validation message, got %q", errText)
		}
	})
}
