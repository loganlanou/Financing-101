package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/loganlanou/Financing-101/internal/database"
	"log/slog"
)

// LearningModule represents a learning course module
type LearningModule struct {
	ID          string
	Title       string
	Description string
	Category    string
	SortOrder   int
	LessonCount int
	CreatedAt   time.Time
}

// Lesson represents an individual lesson within a module
type Lesson struct {
	ID        string
	ModuleID  string
	Title     string
	Content   string
	Summary   string
	SortOrder int
	CreatedAt time.Time
}

// GlossaryTerm represents a financial term definition
type GlossaryTerm struct {
	ID         string
	Term       string
	Definition string
	Category   string
	CreatedAt  time.Time
}

// LearningTip represents a daily learning tip
type LearningTip struct {
	ID         string
	Title      string
	Content    string
	Category   string
	LearnURL   string
	ActiveDate time.Time
	CreatedAt  time.Time
}

// LearnService provides access to learning content
type LearnService struct {
	log     *slog.Logger
	queries *database.Queries
}

func NewLearnService(log *slog.Logger, queries *database.Queries) *LearnService {
	return &LearnService{log: log, queries: queries}
}

// GetModules returns all learning modules
func (s *LearnService) GetModules(ctx context.Context) ([]LearningModule, error) {
	rows, err := s.queries.ListLearningModules(ctx)
	if err != nil {
		return nil, err
	}

	modules := make([]LearningModule, 0, len(rows))
	for _, row := range rows {
		// Get lesson count for each module
		count, err := s.queries.CountLessonsByModule(ctx, row.ID)
		if err != nil {
			s.log.Warn("failed to count lessons", slog.String("module_id", row.ID), slog.Any("err", err))
			count = 0
		}

		modules = append(modules, LearningModule{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Category:    row.Category,
			SortOrder:   int(row.SortOrder),
			LessonCount: int(count),
			CreatedAt:   row.CreatedAt,
		})
	}

	return modules, nil
}

// GetModulesByCategory returns modules filtered by category
func (s *LearnService) GetModulesByCategory(ctx context.Context, category string) ([]LearningModule, error) {
	rows, err := s.queries.ListLearningModulesByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	modules := make([]LearningModule, 0, len(rows))
	for _, row := range rows {
		count, err := s.queries.CountLessonsByModule(ctx, row.ID)
		if err != nil {
			count = 0
		}

		modules = append(modules, LearningModule{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Category:    row.Category,
			SortOrder:   int(row.SortOrder),
			LessonCount: int(count),
			CreatedAt:   row.CreatedAt,
		})
	}

	return modules, nil
}

// GetModule returns a single module with its lessons
func (s *LearnService) GetModule(ctx context.Context, id string) (*LearningModule, []Lesson, error) {
	row, err := s.queries.GetLearningModule(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	lessonRows, err := s.queries.GetLessonsByModule(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	lessons := make([]Lesson, 0, len(lessonRows))
	for _, lr := range lessonRows {
		lessons = append(lessons, Lesson{
			ID:        lr.ID,
			ModuleID:  lr.ModuleID,
			Title:     lr.Title,
			Content:   lr.Content,
			Summary:   lr.Summary,
			SortOrder: int(lr.SortOrder),
			CreatedAt: lr.CreatedAt,
		})
	}

	module := &LearningModule{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Category:    row.Category,
		SortOrder:   int(row.SortOrder),
		LessonCount: len(lessons),
		CreatedAt:   row.CreatedAt,
	}

	return module, lessons, nil
}

// GetGlossary returns all glossary terms
func (s *LearnService) GetGlossary(ctx context.Context) ([]GlossaryTerm, error) {
	rows, err := s.queries.ListGlossaryTerms(ctx)
	if err != nil {
		return nil, err
	}

	terms := make([]GlossaryTerm, 0, len(rows))
	for _, row := range rows {
		terms = append(terms, GlossaryTerm{
			ID:         row.ID,
			Term:       row.Term,
			Definition: row.Definition,
			Category:   row.Category,
			CreatedAt:  row.CreatedAt,
		})
	}

	return terms, nil
}

// SearchGlossary searches glossary terms
func (s *LearnService) SearchGlossary(ctx context.Context, query string) ([]GlossaryTerm, error) {
	rows, err := s.queries.SearchGlossaryTerms(ctx, sql.NullString{String: query, Valid: query != ""})
	if err != nil {
		return nil, err
	}

	terms := make([]GlossaryTerm, 0, len(rows))
	for _, row := range rows {
		terms = append(terms, GlossaryTerm{
			ID:         row.ID,
			Term:       row.Term,
			Definition: row.Definition,
			Category:   row.Category,
			CreatedAt:  row.CreatedAt,
		})
	}

	return terms, nil
}

// GetTodaysTip returns today's learning tip, or a random one if none is set
func (s *LearnService) GetTodaysTip(ctx context.Context) (*LearningTip, error) {
	row, err := s.queries.GetTodaysLearningTip(ctx)
	if err != nil {
		// Fall back to random tip
		row, err = s.queries.GetRandomLearningTip(ctx)
		if err != nil {
			return nil, err
		}
	}

	learnURL := ""
	if row.LearnUrl.Valid {
		learnURL = row.LearnUrl.String
	}

	var activeDate time.Time
	if row.ActiveDate.Valid {
		activeDate = row.ActiveDate.Time
	}

	return &LearningTip{
		ID:         row.ID,
		Title:      row.Title,
		Content:    row.Content,
		Category:   row.Category,
		LearnURL:   learnURL,
		ActiveDate: activeDate,
		CreatedAt:  row.CreatedAt,
	}, nil
}
