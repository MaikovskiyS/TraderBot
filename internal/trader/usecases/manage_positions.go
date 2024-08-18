package usecases

import (
	"context"
)

func (u *useCases) ManagePositions(ctx context.Context) error {
	resp, err := u.Traiding.GetPositions(ctx)
	if err != nil {
		return err
	}

	for _, position := range resp.Positions {
		err := u.Traiding.ManagePosition(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}
