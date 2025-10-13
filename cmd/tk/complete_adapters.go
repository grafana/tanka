package main

import (
	completev1 "github.com/posener/complete"
	completev2 "github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// v2ToV1Predictor adapts a v2 Predictor to v1 Predictor interface
// This is needed because go-clix/cli still uses complete v1
type v2ToV1Predictor struct {
	v2Pred completev2.Predictor
}

func (p v2ToV1Predictor) Predict(args completev1.Args) []string {
	return p.v2Pred.Predict(args.Last)
}

func predictFilesV1(pattern string) completev1.Predictor {
	return v2ToV1Predictor{v2Pred: predict.Files(pattern)}
}

func predictDirsV1(pattern string) completev1.Predictor {
	return v2ToV1Predictor{v2Pred: predict.Dirs(pattern)}
}
