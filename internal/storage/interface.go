package storage

import (
	"context"

	"block-flow/internal/models"
)

// Storage defines the interface for persisting flows and configurations
type Storage interface {
	// Flow operations
	SaveFlow(ctx context.Context, flow *models.Flow) error
	LoadFlow(ctx context.Context, flowID string) (*models.Flow, error)
	LoadAllFlows(ctx context.Context) ([]*models.Flow, error)
	DeleteFlow(ctx context.Context, flowID string) error
	FlowExists(ctx context.Context, flowID string) bool

	// Flow execution operations
	SaveFlowExecution(ctx context.Context, execution *models.FlowExecution) error
	LoadFlowExecution(ctx context.Context, executionID string) (*models.FlowExecution, error)
	LoadFlowExecutions(ctx context.Context, flowID string) ([]*models.FlowExecution, error)
	DeleteFlowExecution(ctx context.Context, executionID string) error

	// Configuration operations
	SaveConfig(ctx context.Context, key string, value interface{}) error
	LoadConfig(ctx context.Context, key string, target interface{}) error
	DeleteConfig(ctx context.Context, key string) error

	// Health and maintenance
	Health(ctx context.Context) error
	Close() error
}
