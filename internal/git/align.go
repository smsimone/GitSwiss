package git

import (
	"context"
	"fmt"

	"it.smaso/git_utilities/pool"
)

// Align aligns the target branch with the source branch
func Align(ctx context.Context, path, source, target string) error {
	results := pool.RunInParallel[string, *BranchResult](
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

	if err := Pull(ctx, path); err != nil {
		return fmt.Errorf("failed to pull source: %s", err.Error())
	}

	if err := Checkout(ctx, path, target); err != nil {
		return fmt.Errorf("failed to checkout to target: %s", err.Error())
	}
	if err := Merge(ctx, path, source); err != nil {
		return fmt.Errorf("failed to merge source into target: %s", err.Error())
	}

	return nil
}
