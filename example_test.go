package wlog

import (
	"context"
	"strconv"
	"testing"
)

func TestExample(t *testing.T) {
	// 在<ProjectMemo>找到了参考 👍🏻: 测试本地日志

	// 测试本地日志
	LDev.Log().Debug("开发消息 1")

	// 创建新的 Factory
	factory, err := NewFactory(createStderrLogger())
	if err != nil {
		t.Fatalf("创建新的 factory 失败: %v", err)
	}

	ctx := context.Background()
	// 使用新创建的 Factory
	log := factory.NewBuilder(ctx).Name("common").Leaf()

	log.Info("开发消息 2")

	// 使用全局方法
	Common("ok").Info("使用默认 wlog 实例打印")
	Common("ok").WithField("dev", true).Info("使用开发 wlog 实例打印")

	// 使用 ByCtx
	Leaf(ctx, "l1").Info("使用 Leaf 条目打印")

	// 使用 Field (Leaf 策略)
	Common("ok").WithField("key1", "value1").Info("使用 Field")

	// 使用 Fields (Leaf 策略)
	Common("ok").WithFields(Fields{"key2": "value2", "key3": "value3"}).Info("使用 Fields")

	// 使用 Branch
	branchLog, branchCtx := Branch(ctx, "newBranch")
	branchLog.Info("使用 Branch")
	// 使用更新后的上下文
	Leaf(branchCtx, "branchedLeaf").Info("使用来自 Branch 的上下文")

	// 使用 Detach
	detachedLog, detachedCtx := Detach(ctx, "detached")
	detachedLog.Info("使用 Detach")
	// 使用更新后的上下文
	Leaf(detachedCtx, "detachedLeaf").Info("使用来自 Detach 的上下文")
}

func TestMFP(t *testing.T) {
	ctx := context.Background()

	l1, ctx1 := Branch(ctx, "l1")
	l1.Info("l1")

	l2, ctx2 := Branch(ctx1, "l2")
	l2.Info("l2")

	l3, ctx3 := Detach(ctx2, "l3")
	l3.Info("l3")

	l4, _ := Detach(ctx3, "l4")
	l4.Info("l4")
}

func BenchmarkExample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Common("ok").Info("使用默认 wlog 实例打印")
	}
}

func BenchmarkWithField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Common("ok").WithField("key", "value").Info("打印单个字段")
	}
}

func BenchmarkWithFields(b *testing.B) {
	fields := Fields{"key1": "value1", "key2": "value2", "key3": "value3"}
	for i := 0; i < b.N; i++ {
		Common("ok").WithFields(fields).Info("打印多个字段")
	}
}

func BenchmarkBranch(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		log, _ := Branch(ctx, "newBranch")
		log.Info("使用新分支打印")
	}
}

func BenchmarkChainedBranchesDeep3(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		_, ctx = Branch(ctx, "branch"+strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		Leaf(ctx, "leaf").Info("链式分支后打印 - 3 层")
	}
}

func BenchmarkChainedBranchesDeep10(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, ctx = Branch(ctx, "branch"+strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		Leaf(ctx, "leaf").Info("链式分支后打印 - 10 层")
	}
}

func BenchmarkChainedBranchesDeep100(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		_, ctx = Branch(ctx, "branch"+strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		Leaf(ctx, "leaf").Info("链式分支后打印 - 100 层")
	}
}

func BenchmarkDisableExample(b *testing.B) {
	DevEnabled.Store(false)
	for i := 0; i < b.N; i++ {
		LDev.Log("ok").Info("使用默认 wlog 实例打印")
	}
}

func BenchmarkDisableExample2(b *testing.B) {
	DevEnabled.Store(false)
	d := LDev.Log("ok")
	for i := 0; i < b.N; i++ {
		d.Info("使用默认 wlog 实例打印")
	}
}
