//Copyright 2020 Google LLC
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//https://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package zerrors_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/JavierZunzunegui/zerrors"
)

// filename needs updating if this file is renamed.
const fileName = "error_test.go"

const separator = `: `

const frameRegString = `\(error_test\.go\:[1-9][0-9]*\)`

var frameReg = func() *regexp.Regexp {
	r := regexp.MustCompile(frameRegString)
	if exampleFrame := "(" + fileName + ":100)"; !r.MatchString(exampleFrame) {
		panic(fmt.Sprintf("example frame %q not matched the frame format '%s'", exampleFrame, r.String()))
	}
	return r
}()

func TestWrap(t *testing.T) {
	t.Run("nil cases", func(t *testing.T) {
		if err := zerrors.Value(nil); err != nil {
			t.Errorf("Value(nil): expecting nil, got %q", err)
		}

		if frame, ok := zerrors.Frame(nil); ok {
			t.Errorf("Frame(nil): expecting false, got true with frame: %v", frame)
		}

		// Anti-pattern.
		if err := zerrors.New(nil); err != nil {
			t.Errorf("New(nil): expecting nil, got %q", err)
		}

		// Anti-pattern.
		if err := zerrors.Wrap(nil, nil); err != nil {
			t.Errorf("Wrap(nil, nil): expecting nil, got %q", err)
		}
	})

	const baseMsg = "some error"
	baseErr := errors.New(baseMsg)

	t.Run("non-wrapped cases", func(t *testing.T) {
		if err := zerrors.Value(baseErr); err != baseErr {
			t.Errorf("Value(baseErr): expecting baseErr=%q, got %q", baseErr, err)
		}

		if frame, ok := zerrors.Frame(baseErr); ok {
			t.Errorf("Frame(baseErr): expecting false, got true with frame: %v", frame)
		}

		if err := zerrors.Wrap(nil, baseErr); err != nil {
			t.Errorf("Wrap(nil, baseErr): expecting nil, got %q", err)
		}
		// Anti-pattern.
		if err := zerrors.Wrap(baseErr, nil); err != baseErr {
			t.Errorf("Wrap(baseErr, nil): expecting baseErr=%q, got %q", baseErr, err)
		}

		if err := zerrors.SWrap(nil, baseMsg); err != nil {
			t.Errorf("SWrap(nil, baseMsg): expecting nil, got %q", err)
		}
	})

	const basic1 = baseMsg
	detail1 := regexp.MustCompile(baseMsg + " " + frameRegString)
	w1Err := zerrors.New(baseErr)
	s1Err := zerrors.SNew(baseMsg)

	t.Run("first level wrapping", func(t *testing.T) {
		if msg := w1Err.Error(); basic1 != msg {
			t.Errorf("w1Err.Error(): expecting %q, got %q", basic1, msg)
		}
		if msg := zerrors.Detail(w1Err); !detail1.MatchString(msg) {
			t.Errorf("Detail(w1Err): expecting regex '%s', got %q", detail1.String(), msg)
		}
		if err := zerrors.Value(w1Err); err != baseErr {
			t.Errorf("Value(w1Err): expecting baseErr=%q, got %q", baseErr, err)
		}
		if _, ok := zerrors.Frame(w1Err); !ok {
			t.Error("Frame(w1Err): expecting true, got false")
		}
		if err := errors.Unwrap(w1Err); err != nil {
			t.Errorf("Unwrap(w1Err): expecting nil, got %q", err)
		}

		if msg := s1Err.Error(); basic1 != msg {
			t.Errorf("s1Err.Error(): expecting '%q', got %q", basic1, msg)
		}
		if msg := zerrors.Value(s1Err).Error(); msg != baseMsg {
			t.Errorf("Value(SNew(baseMsg)).Error(): expecting err=%q, got %q", baseMsg, msg)
		}
		if _, ok := zerrors.Frame(s1Err); !ok {
			t.Error("Frame(s1Err): expecting true, got false")
		}
		if err := errors.Unwrap(s1Err); err != nil {
			t.Errorf("Unwrap(s1Err): expecting nil, got %q", err)
		}

		if err := zerrors.New(w1Err); err != w1Err {
			t.Errorf("New(w1Err): expecting w1Err=%q, got %q", w1Err, err)
		}

		if err := zerrors.Wrap(nil, w1Err); err != nil {
			t.Errorf("Wrap(nil, w1Err): expecting nil, got %q", err)
		}
		// Anti-pattern.
		if err := zerrors.Wrap(w1Err, nil); err != w1Err {
			t.Errorf("Wrap(w1Err, nil): expecting w1Err=%q, got %q", w1Err, err)
		}
	})

	const secondMsg = "second error"
	secondErr := errors.New(secondMsg)
	const basic2 = secondMsg + separator + basic1
	detail2 := regexp.MustCompile(secondMsg + " " + frameRegString + regexp.QuoteMeta(separator) + detail1.String())
	w2Err := zerrors.Wrap(w1Err, secondErr)
	s2Err := zerrors.SWrap(s1Err, secondMsg)

	t.Run("second level wrapping", func(t *testing.T) {
		if msg := w2Err.Error(); basic2 != msg {
			t.Errorf("w2Err.Error(): expecting '%s', got %q", basic2, msg)
		}
		if msg := zerrors.Detail(w2Err); !detail2.MatchString(msg) {
			t.Errorf("Detail(w2Err): expecting regex '%s', got %q", detail2, msg)
		}
		if err := zerrors.Value(w2Err); err != secondErr {
			t.Errorf("Value(w2Err): expecting secondErr=%q, got %q", secondErr, err)
		}
		if _, ok := zerrors.Frame(w2Err); !ok {
			t.Error("Frame(w2Err): expecting true, got false")
		}
		if err := errors.Unwrap(w2Err); err != w1Err {
			t.Errorf("Unwrap(w2Err): expecting w1Err=%q, got %q", w1Err, err)
		}
		if err := errors.Unwrap(errors.Unwrap(w2Err)); err != nil {
			t.Errorf("Unwrap(Unwrap(w2Err)): expecting nil, got %q", err)
		}

		if msg := s2Err.Error(); basic2 != msg {
			t.Errorf("s2Err.Error(): expecting '%s', got %q", basic2, msg)
		}
		if msg := zerrors.Value(s2Err).Error(); msg != secondMsg {
			t.Errorf("Value(s2Err).Error(): expecting secondMsg=%q, got %q", secondMsg, msg)
		}
		if _, ok := zerrors.Frame(s2Err); !ok {
			t.Error("Frame(s2Err): expecting true, got false")
		}
		if msg := zerrors.Value(errors.Unwrap(s2Err)).Error(); msg != baseMsg {
			t.Errorf("Value(Unwrap(s2Err)).Error(): expecting baseMsg=%q, got %q", baseMsg, msg)
		}
		if err := errors.Unwrap(errors.Unwrap(s2Err)); err != nil {
			t.Errorf("Unwrap(Unwrap(s2Err)): expecting nil, got %q", err)
		}

		// Anti-pattern.
		if msg := zerrors.Detail(zerrors.Wrap(w1Err, zerrors.New(secondErr))); !detail2.MatchString(msg) {
			t.Errorf("Detail(Wrap(w1Err, New(secondErr))): expecting regex '%s', got %q", detail2, msg)
		}
		// Anti-pattern.
		if msg, reg := zerrors.Detail(zerrors.Wrap(baseErr, secondErr)), regexp.MustCompile(secondMsg+" "+frameRegString+regexp.QuoteMeta(separator)+baseMsg); !reg.MatchString(msg) {
			t.Errorf("Detail(Wrap(baseErr, secondErr)): expecting regex '%s', got %q", reg, msg)
		}
	})
}

func TestWrapError_Is(t *testing.T) {
	inErr := errors.New("some error")
	wErr := zerrors.New(inErr)
	outErr := errors.New("wrapping error")
	wErr = zerrors.Wrap(wErr, outErr)
	if !errors.Is(wErr, inErr) {
		t.Error("Is(wErr, inErr): should be true")
	}
	if !errors.Is(wErr, outErr) {
		t.Error("Is(wErr, outErr): should be true")
	}
	if errors.Is(wErr, errors.New("unknown error")) {
		t.Error("Is(Wrap(err1, err2), err3): should be false")
	}
}

type custom1Error struct{ msg string }

func (err custom1Error) Error() string { return err.msg }

type custom2Error struct{ msg string }

func (err custom2Error) Error() string { return err.msg }

func TestWrapError_As(t *testing.T) {
	inErr := custom1Error{"some error"}
	wErr := zerrors.New(inErr)
	outErr := errors.New("wrapping error")
	wErr = zerrors.Wrap(wErr, outErr)
	err1 := custom1Error{}
	if !errors.As(wErr, &err1) {
		t.Error("As(wErr, &err1): should be true")
	} else if msg := err1.Error(); msg != inErr.Error() {
		t.Errorf("As(wErr, &err1): expecting err1.Error()=%q, got %q", inErr.Error(), msg)
	}
	err2 := custom2Error{}
	if errors.As(wErr, &err2) {
		t.Error("As(wErr, &err2): should be false")
	}
}

func TestWrapError_Format(t *testing.T) {
	const s = ` ` + frameRegString

	const baseMsg = "base"
	baseErr := zerrors.SNew(baseMsg)

	const basic1 = baseMsg
	detail1 := regexp.MustCompile(baseMsg + s)

	if s := fmt.Sprintf("%s", baseErr); s != basic1 {
		t.Errorf("fmt.Sprintf(%s, baseErr): expected '%s', got '%s'", "%s", basic1, s)
	}
	if s := fmt.Sprintf("%v", baseErr); s != basic1 {
		t.Errorf("fmt.Sprintf(%s, baseErr): expected '%s', got '%s'", "%v", basic1, s)
	}
	if s, expected := fmt.Sprintf("%q", baseErr), `"`+baseMsg+`"`; s != expected {
		t.Errorf("fmt.Sprintf(%s, baseErr): expected '%s', got '%s'", "%q", expected, s)
	}
	if s := fmt.Sprintf("%+v", baseErr); !detail1.MatchString(s) {
		t.Errorf("fmt.Sprintf(%s, baseErr): expected regex '%s', got '%s'", "%+v", detail1.String(), s)
	}

	const secondMsg = "second error"
	err2 := zerrors.SWrap(baseErr, secondMsg)

	const basic2 = secondMsg + separator + basic1
	detail2 := regexp.MustCompile(secondMsg + " " + frameRegString + regexp.QuoteMeta(separator) + detail1.String())

	if s := fmt.Sprintf("%s", err2); s != basic2 {
		t.Errorf("fmt.Sprintf(%s, err2): expected '%s', got '%s'", "%s", basic2, s)
	}
	if s := fmt.Sprintf("%v", err2); s != basic2 {
		t.Errorf("fmt.Sprintf(%s, err2): expected '%s', got '%s'", "%v", basic2, s)
	}
	if s, expected := fmt.Sprintf("%q", err2), `"`+basic2+`"`; s != expected {
		t.Errorf("fmt.Sprintf(%s, err2): expected '%s', got '%s'", "%q", expected, s)
	}
	if s := fmt.Sprintf("%+v", err2); !detail2.MatchString(s) {
		t.Errorf("fmt.Sprintf(%s, err2): expected regex '%s', got '%s'", "%+v", detail1.String(), s)
	}
}
