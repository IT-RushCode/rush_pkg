package routes

import (
	h "github.com/IT-RushCode/rush_pkg/handlers/review"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_REVIEW(api fiber.Router, repo *repositories.Repositories) {
	reviewHandler := h.NewReviewHandler(repo)

	review := api.Group("reviews")

	review.Get("/", reviewHandler.GetAllReviews)
	review.Get("/:id", reviewHandler.FindReviewByID)
	review.Post("/", reviewHandler.CreateReview)
	review.Put("/:id", reviewHandler.UpdateReview)
	review.Delete("/:id", reviewHandler.DeleteReview)
}
