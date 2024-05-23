package git

import (
	"context"
	"fmt"

	"it.smaso/git_swiss/pool"
)

type GitStrategy string

const (
	MERGE_STRATEGY GitStrategy = "merge"
	PULL_STRATEGY  GitStrategy = "pull"
)

func StrategyFromString(strategy string) (GitStrategy, error) {
	switch strategy {
	case "merge":
		return MERGE_STRATEGY, nil
	case "pull":
		return PULL_STRATEGY, nil
	default:
		return MERGE_STRATEGY, fmt.Errorf("strategy %s is not valid", strategy)
	}
}

// Align aligns the target branch with the source branch
func Align(ctx context.Context, path, source, target string, strategy GitStrategy, remote string) error {
	results := pool.RunInParallel(
		func() pool.KeyedResult[string, *BranchResult] {
			return pool.KeyedResult[string, *BranchResult]{
				Key:   "source",
				Value: BranchExists(ctx, path, source),
			}
		},
		func() pool.KeyedResult[string, *BranchResult] {
			return pool.KeyedResult[string, *BranchResult]{
				Key:   "target",
				Value: BranchExists(ctx, path, target),
			}
		},
	)

	sourceRes := pool.FindByKey("source", results).Value
	targetRes := pool.FindByKey("target", results).Value

	if sourceRes == nil && targetRes == nil {
		return fmt.Errorf("source and target branches do not exist")
	} else if sourceRes == nil {
		return fmt.Errorf("source branch does not exist")
	} else if targetRes == nil {
		return fmt.Errorf("target branch does not exist")
	}

	if !sourceRes.IsCurrent {
		if err := Checkout(ctx, path, source); err != nil {
			return fmt.Errorf("failed to checkout to source: %s", err.Error())
		}
	}

	if strategy == MERGE_STRATEGY {
		if err := Pull(ctx, path); err != nil {
			return fmt.Errorf("failed to pull source: %s", err.Error())
		}

		if err := Checkout(ctx, path, target); err != nil {
			return fmt.Errorf("failed to checkout to target: %s", err.Error())
		}
		if err := Merge(ctx, path, source); err != nil {
			return fmt.Errorf("failed to merge source into target: %s", err.Error())
		}
	} else {
		if err := MergePull(ctx, path, remote, source); err != nil {
			return fmt.Errorf("failed to merge with pull strategy: %s", err.Error())
		}
	}

	return nil
}
