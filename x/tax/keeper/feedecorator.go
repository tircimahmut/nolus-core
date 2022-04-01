package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var HUNDRED = sdk.NewDec(100)

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct {
	tk Keeper
}

func NewMempoolFeeDecorator(tk Keeper) MempoolFeeDecorator {
	return MempoolFeeDecorator{
		tk: tk,
	}
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	params := mfd.tk.GetParams(ctx)

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	ctx.Logger().Info(fmt.Sprintf("Mempool: gas %d", gas))

	feeRate := sdk.NewDec(int64(params.FeeRate))

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !simulate {

		// deduct the nomo fee from feeCoins
		taxFee, feeRemaining, err := ApplyFee(feeRate, feeCoins)
		if err != nil {
			return ctx, err
		}

		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}
			// ensure that enough was paid to cover the validator tax after the custom tax was deduced
			if !feeRemaining.IsAnyGTE(requiredFees) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeRemaining.Add(taxFee...), requiredFees.Add(taxFee...))
			}
		}
	}

	return next(ctx, tx, simulate)
}

//TODO Rename this method
func ApplyFeeImpl(feeRate sdk.Dec, feeCoins sdk.Coins) (sdk.Coins, sdk.Coins, error) {
	proceeds := sdk.Coins{}

	if feeRate.IsZero() {
		return proceeds, feeCoins, nil
	}

	if feeCoins.Empty() {
		return sdk.NewCoins(), feeCoins, nil
	}

	// we will deduct the fee from every denomination send
	for _, fee := range feeCoins {
		proceed := sdk.NewCoin(fee.Denom, feeRate.MulInt(fee.Amount).Quo(HUNDRED).TruncateInt())
		proceeds = proceeds.Add(proceed)
	}

	deductedFees, neg := feeCoins.SafeSub(proceeds)
	if neg {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "ApplyFee: insufficient fees; got: %s required: %s", feeCoins, proceeds)
	}

	return proceeds, deductedFees, nil
}