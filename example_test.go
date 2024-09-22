package wlog

import (
	"context"
	"testing"
)

func TestExample(t *testing.T) {
	// test local log
	LDev.Log().Debug("dev msg 1")

	// create new Factory
	wlog, err := NewFactory(createStderrLogger())
	if err != nil {
		t.Fatalf("Failed to create new factory: %v", err)
	}

	ctx := context.Background()
	// use new created Factory
	log, _ := wlog.NewBuilder(ctx).
		WithStrategy(ForkLeaf).
		WithFingerPrints("common").
		Build()
	log.Info("dev msg 2")

	// use global method
	Common("ok").Info("print by default wlog instance")
	Common("ok").WithField("dev", true).Info("print by dev wlog instance")

	// use ByCtx
	Leaf(ctx, "l1").Info("print by Leaf entry")

	// use WithField (Leaf strategy)
	Common("ok").WithField("key1", "value1").Info("Using WithField")

	// use WithFields (Leaf strategy)
	Common("ok").WithFields(Fields{"key2": "value2", "key3": "value3"}).Info("Using WithFields")

	// use BranchField
	branchLog, branchCtx := Common("ok").BranchField(ctx, "branchKey", "branchValue")
	branchLog.Info("Using BranchField")
	// use updated context
	Leaf(branchCtx, "branchLeaf").Info("Using context from BranchField")

	// use BranchFields
	multiBranchLog, multiBranchCtx := Common("ok").BranchFields(ctx, Fields{"multiKey1": "multiValue1", "multiKey2": "multiValue2"})
	multiBranchLog.Info("Using BranchFields")
	// use updated context
	Leaf(multiBranchCtx, "multiBranchLeaf").Info("Using context from BranchFields")

	// use Branch
	branchedLog, branchedCtx := Common("ok").Branch(ctx, "newBranch")
	branchedLog.Info("Using Branch")
	// use updated context
	Leaf(branchedCtx, "branchedLeaf").Info("Using context from Branch")
}

func TestMFP(t *testing.T) {
	ctx := context.Background()

	l1, ctx1 := Branch(ctx, "l1")
	l1.Info("l1")

	l2, ctx2 := Branch(ctx1, "l2")
	l2.Info("l2")

	l3, ctx3 := DetachNew(ctx2, "l3")
	l3.Info("l3")

	l4, _ := DetachNew(ctx3, "l4")
	l4.Info("l4")
}

func BenchmarkExample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Common("ok").Info("print by default wlog instance")
	}
}

func BenchmarkWithField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Common("ok").WithField("key", "value").Info("print with single field")
	}
}

func BenchmarkWithFields(b *testing.B) {
	fields := Fields{"key1": "value1", "key2": "value2", "key3": "value3"}
	for i := 0; i < b.N; i++ {
		Common("ok").WithFields(fields).Info("print with multiple fields")
	}
}

func BenchmarkBranchField(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		log, _ := Common("ok").BranchField(ctx, "key", "value")
		log.Info("print with branch field")
	}
}

func BenchmarkBranchFields(b *testing.B) {
	ctx := context.Background()
	fields := Fields{"key1": "value1", "key2": "value2", "key3": "value3"}
	for i := 0; i < b.N; i++ {
		log, _ := Common("ok").BranchFields(ctx, fields)
		log.Info("print with branch fields")
	}
}

func BenchmarkBranch(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		log, _ := Common("ok").Branch(ctx, "newBranch")
		log.Info("print with new branch")
	}
}

func BenchmarkChainedBranches(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		log1, ctx1 := Common("ok").Branch(ctx, "branch1")
		log2, ctx2 := log1.Branch(ctx1, "branch2")
		log3, _ := log2.Branch(ctx2, "branch3")
		log3.Info("print after chained branches")
	}
}

func BenchmarkChainedBranchFields(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		log1, ctx1 := Common("ok").BranchField(ctx, "key1", "value1")
		log2, ctx2 := log1.BranchField(ctx1, "key2", "value2")
		log3, _ := log2.BranchField(ctx2, "key3", "value3")
		log3.Info("print after chained branch fields")
	}
}

func BenchmarkDisableExample(b *testing.B) {
	DevEnabled.Store(false)
	for i := 0; i < b.N; i++ {
		LDev.Log("ok").Info("print by default wlog instance")
	}
}

func BenchmarkDisableExample2(b *testing.B) {
	DevEnabled.Store(false)
	d := LDev.Log("ok")
	for i := 0; i < b.N; i++ {
		d.Info("print by default wlog instance")
	}
}
