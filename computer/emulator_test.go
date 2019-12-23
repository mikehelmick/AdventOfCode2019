package computer

import (
	"testing"
)

func run(data []int64, input []int64) chan Output {
	in := make(chan int64, 10)
	out := make(chan Output, 10)

	for _, i := range input {
		in <- i
	}
	e := NewEmulator(data, in, out, false)
	e.Execute()

	close(in)
	return out
}

func checkTerminated(t *testing.T, outc chan Output) {
	if v := <-outc; !v.Done {
		t.Errorf("Machine should have terminated, but didn't.")
	}
}

func TestAdditionPosition(t *testing.T) {
	d := []int64{1, 3, 4, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 109 {
		t.Errorf("Output was %v; want 109", v)
	}
	checkTerminated(t, outc)
}

func TestAdditionImmediate(t *testing.T) {
	d := []int64{1101, 42, 42, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 84 {
		t.Errorf("Output was %v; want 109", v)
	}
	checkTerminated(t, outc)
}

func TestAdditionRelative(t *testing.T) {
	d := []int64{109, 3, 2201, -1, 0, 7, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 2200 {
		t.Errorf("Output was %v; want 109", v)
	}
	checkTerminated(t, outc)
}

func TestMultiplicationPosition(t *testing.T) {
	d := []int64{2, 3, 4, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(5 * 104)
	if v := <-outc; v.Val != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
	checkTerminated(t, outc)
}

func TestMultiplicationImmediate(t *testing.T) {
	d := []int64{1102, 42, 42, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(42 * 42)
	if v := <-outc; v.Val != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
	checkTerminated(t, outc)
}

func TestMultiplicationRelative(t *testing.T) {
	d := []int64{109, 3, 2202, -1, 0, 7, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(2202 * -1)
	if v := <-outc; v.Val != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
	checkTerminated(t, outc)
}

func TestRelativeBasePositionMode(t *testing.T) {
	d := []int64{9, 1, 204, 1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 204 {
		t.Errorf("Output was %v; want 204", v)
	}
	checkTerminated(t, outc)
}

func TestRelativeBaseImmediateMode(t *testing.T) {
	d := []int64{109, 5, 204, 0, 99, 12345}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 12345 {
		t.Errorf("Output was %v; want 12345", v)
	}
	checkTerminated(t, outc)
}

func TestRelativeBaseRelativeMode(t *testing.T) {
	d := []int64{209, 6, 204, 0, 99, 12345, 5}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 12345 {
		t.Errorf("Output was %v; want 12345", v)
	}
	checkTerminated(t, outc)
}

func TestInputPosition(t *testing.T) {
	d := []int64{03, 3, 104, 99, 99}
	outc := run(d, []int64{42})
	defer close(outc)

	if v := <-outc; v.Val != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
	checkTerminated(t, outc)
}

func TestInputRelative(t *testing.T) {
	d := []int64{109, 7, 203, 0, 204, 0, 99, 0}
	outc := run(d, []int64{42})
	defer close(outc)

	if v := <-outc; v.Val != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
	checkTerminated(t, outc)
}

func TestNonBlockingInput(t *testing.T) {
	in := make(chan int64, 10)
	out := make(chan Output, 10)
	data := []int64{03, 12, 04, 12, 1007, 12, 0, 13, 1005, 13, 0, 99, 0, 0, 00}
	e := NewEmulator(data, in, out, true)
	go e.Execute()

	target := int64(1234)
	outputVals := make([]int64, 0, 100)
	for {
		output := <-out
		if output.Done {
			break
		}

		outputVals = append(outputVals, output.Val)
		in <- target
	}

	found := false
	for _, v := range outputVals {
		if v == target {
			found = true
		}
	}

	if !found {
		t.Errorf("Output expected %v, output was %v", target, outputVals)
	}

	close(in)
	close(out)
}

func TestJumpIfTruePosition_False(t *testing.T) {
	d := []int64{05, 9, 6, 104, -1, 99, 104, 42, 99, 0}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
	checkTerminated(t, outc)
}

func TestJumpIfTruePosition_True(t *testing.T) {
	d := []int64{05, 9, 9, 104, -1, 99, 104, 42, 99, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
	checkTerminated(t, outc)
}

func TestJumpIfTrueDirect_False(t *testing.T) {
	d := []int64{1105, 0, 9, 104, -1, 99, 104, 42, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
	checkTerminated(t, outc)
}

func TestJumpIfTrueDirect_True(t *testing.T) {
	d := []int64{1105, 1, 6, 104, -1, 99, 104, 42, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
	checkTerminated(t, outc)
}

func TestJumpIfFalsePosition_False(t *testing.T) {
	d := []int64{06, 9, 10, 104, -1, 99, 104, 42, 99, 0, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
	checkTerminated(t, outc)
}

func TestJumpIfFalsePosition_True(t *testing.T) {
	d := []int64{06, 9, 10, 104, -1, 99, 104, 42, 99, 1, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v.Val != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
	checkTerminated(t, outc)
}

// By this time, assume memory addressing works, focus on functionality.
func TestLessThan(t *testing.T) {
	d := []int64{1107, 1, 2, 13, 1107, 2, 1, 15, 1107, 2, 2, 17, 104, -1, 104, -1, 104, -1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	exp := []int64{1, 0, 0}
	for _, e := range exp {
		if v := <-outc; v.Val != e {
			t.Errorf("Output was %v; want %v", v, e)
		}
	}
	checkTerminated(t, outc)
}

func TestEquals(t *testing.T) {
	d := []int64{1108, 1, 1, 13, 1108, 2, 1, 15, 1108, 2, 2, 17, 104, -1, 104, -1, 104, -1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	exp := []int64{1, 0, 1}
	for _, e := range exp {
		if v := <-outc; v.Val != e {
			t.Errorf("Output was %v; want %v", v, e)
		}
	}
	checkTerminated(t, outc)
}
