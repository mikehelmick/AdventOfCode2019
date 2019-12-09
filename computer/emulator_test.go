package computer

import "testing"

func run(data []int64, input []int64) chan int64 {
	in := make(chan int64, 10)
	out := make(chan int64, 10)
	done := make(chan bool, 1)

	for _, i := range input {
		in <- i
	}
	e := NewEmulator(data, in, out, done)
	e.Execute()
	<-done

	close(in)
	close(done)
	return out
}

func TestAdditionPosition(t *testing.T) {
	d := []int64{1, 3, 4, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 109 {
		t.Errorf("Output was %v; want 109", v)
	}
}

func TestAdditionImmediate(t *testing.T) {
	d := []int64{1101, 42, 42, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 84 {
		t.Errorf("Output was %v; want 109", v)
	}
}

func TestAdditionRelative(t *testing.T) {
	d := []int64{109, 3, 2201, -1, 0, 7, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 2200 {
		t.Errorf("Output was %v; want 109", v)
	}
}

func TestMultiplicationPosition(t *testing.T) {
	d := []int64{2, 3, 4, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(5 * 104)
	if v := <-outc; v != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
}

func TestMultiplicationImmediate(t *testing.T) {
	d := []int64{1102, 42, 42, 5, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(42 * 42)
	if v := <-outc; v != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
}

func TestMultiplicationRelative(t *testing.T) {
	d := []int64{109, 3, 2202, -1, 0, 7, 104, 0, 99}
	outc := run(d, []int64{})
	defer close(outc)

	expected := int64(2202 * -1)
	if v := <-outc; v != expected {
		t.Errorf("Output was %v; want %v", v, expected)
	}
}

func TestRelativeBasePositionMode(t *testing.T) {
	d := []int64{9, 1, 204, 1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 204 {
		t.Errorf("Output was %v; want 204", v)
	}
}

func TestRelativeBaseImmediateMode(t *testing.T) {
	d := []int64{109, 5, 204, 0, 99, 12345}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 12345 {
		t.Errorf("Output was %v; want 12345", v)
	}
}

func TestRelativeBaseRelativeMode(t *testing.T) {
	d := []int64{209, 6, 204, 0, 99, 12345, 5}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 12345 {
		t.Errorf("Output was %v; want 12345", v)
	}
}

func TestInputPosition(t *testing.T) {
	d := []int64{03, 3, 104, 99, 99}
	outc := run(d, []int64{42})
	defer close(outc)

	if v := <-outc; v != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
}

func TestInputRelative(t *testing.T) {
	d := []int64{109, 7, 203, 0, 204, 0, 99, 0}
	outc := run(d, []int64{42})
	defer close(outc)

	if v := <-outc; v != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
}

func TestJumpIfTruePosition_False(t *testing.T) {
	d := []int64{05, 9, 6, 104, -1, 99, 104, 42, 99, 0}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
}

func TestJumpIfTruePosition_True(t *testing.T) {
	d := []int64{05, 9, 9, 104, -1, 99, 104, 42, 99, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
}

func TestJumpIfTrueDirect_False(t *testing.T) {
	d := []int64{1105, 0, 9, 104, -1, 99, 104, 42, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
}

func TestJumpIfTrueDirect_True(t *testing.T) {
	d := []int64{1105, 1, 6, 104, -1, 99, 104, 42, 99}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
}

func TestJumpIfFalsePosition_False(t *testing.T) {
	d := []int64{06, 9, 10, 104, -1, 99, 104, 42, 99, 0, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != 42 {
		t.Errorf("Output was %v; want 42", v)
	}
}

func TestJumpIfFalsePosition_True(t *testing.T) {
	d := []int64{06, 9, 10, 104, -1, 99, 104, 42, 99, 1, 6}
	outc := run(d, []int64{})
	defer close(outc)

	if v := <-outc; v != -1 {
		t.Errorf("Output was %v; want -1", v)
	}
}

// By this time, assume memory addressing works, focus on functionality.
func TestLessThan(t *testing.T) {
	d := []int64{1107, 1, 2, 13, 1107, 2, 1, 15, 1107, 2, 2, 17, 104, -1, 104, -1, 104, -1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	exp := []int64{1, 0, 0}
	for _, e := range exp {
		if v := <-outc; v != e {
			t.Errorf("Output was %v; want %v", v, e)
		}
	}
}

func TestEquals(t *testing.T) {
	d := []int64{1108, 1, 1, 13, 1108, 2, 1, 15, 1108, 2, 2, 17, 104, -1, 104, -1, 104, -1, 99}
	outc := run(d, []int64{})
	defer close(outc)

	exp := []int64{1, 0, 1}
	for _, e := range exp {
		if v := <-outc; v != e {
			t.Errorf("Output was %v; want %v", v, e)
		}
	}
}
