package redis_registry

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {
	if _, err := Connect(context.Background(), "127.0.0.1:6379", "123456", 0); err != nil {
		t.Fatal(err)
	}

}

func TestScriptSetNXAndGet(t *testing.T) {
	ctx := context.Background()
	rkey := "aa"
	value1 := "bb1"
	value2 := "bb1"
	if r, err := Connect(context.Background(), "127.0.0.1:6379", "123456", 0); err != nil {
		t.Fatal(err)
	} else {
		if result := r.rdb.Del(ctx, rkey); result.Err() != nil {
			t.Fatal(result.Err())
		}
		// aa=bb1 return 1, nil
		if result := r.rdb.EvalSha(ctx, r.setNXAndGetScript, []string{rkey}, value1, 60); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if arr[1] != nil {
			t.Fatal("result 2 must nil")
		} else if v1 != 1 {
			t.Fatal("result 1 must equal 1")
		}
		// aa=bb1 return 0, bb1
		if result := r.rdb.EvalSha(ctx, r.setNXAndGetScript, []string{rkey}, value1, 60); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if v2, ok := arr[1].(string); !ok {
			t.Fatal("result 2 must string")
		} else if v1 != 0 {
			t.Fatal("result 1 must equal 1")
		} else if v2 != value1 {
			t.Fatalf("result 1 must equal %s", value1)
		}
		// aa=bb2 return 0, bb1
		if result := r.rdb.EvalSha(ctx, r.setNXAndGetScript, []string{rkey}, value2, 60); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if v2, ok := arr[1].(string); !ok {
			t.Fatal("result 2 must string")
		} else if v1 != 0 {
			t.Fatal("result 1 must equal 1")
		} else if v2 != value1 {
			t.Fatalf("result 1 must equal %s", value1)
		}
	}
}

func TestScriptCompareAndDeleteScript(t *testing.T) {
	ctx := context.Background()
	rkey := "aa"
	value1 := "bb1"
	value2 := "bb2"
	if r, err := Connect(context.Background(), "127.0.0.1:6379", "123456", 0); err != nil {
		t.Fatal(err)
	} else {
		// del aa
		if result := r.rdb.Del(ctx, rkey); result.Err() != nil {
			t.Fatal(result.Err())
		}
		// del aa return 0, nil
		if result := r.rdb.EvalSha(ctx, r.compareAndDeleteScript, []string{rkey}, value1); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if arr[1] != nil {
			t.Fatal("result 2 must nil")
		} else if v1 != 0 {
			t.Fatal("result 1 must equal 0")
		}
		// set aa bb1
		if result := r.rdb.Set(ctx, rkey, value1, 0); result.Err() != nil {
			t.Fatal(result.Err())
		}
		// del aa return 0, bb1
		if result := r.rdb.EvalSha(ctx, r.compareAndDeleteScript, []string{rkey}, value2); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if v2, ok := arr[1].(string); !ok {
			t.Fatal("result 2 must string")
		} else if v1 != 0 {
			t.Fatal("result 1 must equal 0")
		} else if v2 != value1 {
			t.Fatalf("result 2 must equal %s", value1)
		}
		// del aa return 1, bb1
		if result := r.rdb.EvalSha(ctx, r.compareAndDeleteScript, []string{rkey}, value1); result.Err() != nil {
			t.Fatal(result.Err())
		} else if arr, err := result.Slice(); err != nil {
			t.Fatal(err)
		} else if len(arr) != 2 {
			t.Fatal("result len must equal 2")
		} else if v1, ok := arr[0].(int64); !ok {
			t.Fatal("result 1 must int64")
		} else if v2, ok := arr[1].(string); !ok {
			t.Fatal("result 2 must string")
		} else if v1 != 1 {
			t.Fatal("result 1 must equal 1")
		} else if v2 != value1 {
			t.Fatalf("result 2 must equal %s", value1)
		}
	}
}
