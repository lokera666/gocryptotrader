package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common/file"
)

func TestSendHTTPRequest(t *testing.T) {
	// t.Parallel() not used to maintain code coverage for assigning the default
	// HTTPClient.
	methodPost := "pOst"
	methodGet := "GeT"
	methodDelete := "dEleTe"
	methodGarbage := "ding"

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	_, err := SendHTTPRequest(context.Background(),
		methodGarbage, "https://www.google.com", headers,
		strings.NewReader(""), true,
	)
	if err == nil {
		t.Error("Expected error 'invalid HTTP method specified'")
	}
	_, err = SendHTTPRequest(context.Background(),
		methodPost, "https://www.google.com", headers,
		strings.NewReader(""), true,
	)
	if err != nil {
		t.Error(err)
	}
	_, err = SendHTTPRequest(context.Background(),
		methodGet, "https://www.google.com", headers,
		strings.NewReader(""), true,
	)
	if err != nil {
		t.Error(err)
	}

	err = SetHTTPUserAgent("GCTbot/1337.69 (+http://www.lol.com/)")
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}

	_, err = SendHTTPRequest(context.Background(),
		methodDelete, "https://www.google.com", headers,
		strings.NewReader(""), true,
	)
	if err != nil {
		t.Error(err)
	}
	_, err = SendHTTPRequest(context.Background(),
		methodGet, ":missingprotocolscheme", headers,
		strings.NewReader(""), true,
	)
	if err == nil {
		t.Error("Common HTTPRequest accepted missing protocol")
	}
	_, err = SendHTTPRequest(context.Background(),
		methodGet, "test://unsupportedprotocolscheme", headers,
		strings.NewReader(""), true,
	)
	if err == nil {
		t.Error("Common HTTPRequest accepted invalid protocol")
	}
}

func TestSetHTTPClientWithTimeout(t *testing.T) {
	t.Parallel()
	err := SetHTTPClientWithTimeout(-0)
	if !errors.Is(err, errCannotSetInvalidTimeout) {
		t.Fatalf("received: %v but expected: %v", err, errCannotSetInvalidTimeout)
	}

	err = SetHTTPClientWithTimeout(time.Second * 15)
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}
}

func TestSetHTTPUserAgent(t *testing.T) {
	t.Parallel()
	err := SetHTTPUserAgent("")
	if !errors.Is(err, errUserAgentInvalid) {
		t.Fatalf("received: %v but expected: %v", err, errUserAgentInvalid)
	}

	err = SetHTTPUserAgent("testy test")
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}
}

func TestSetHTTPClient(t *testing.T) {
	t.Parallel()
	err := SetHTTPClient(nil)
	if !errors.Is(err, errHTTPClientInvalid) {
		t.Fatalf("received: %v but expected: %v", err, errHTTPClientInvalid)
	}

	err = SetHTTPClient(new(http.Client))
	if !errors.Is(err, nil) {
		t.Fatalf("received: %v but expected: %v", err, nil)
	}
}

func TestIsEnabled(t *testing.T) {
	t.Parallel()
	expected := "Enabled"
	actual := IsEnabled(true)
	if actual != expected {
		t.Errorf("Expected %s. Actual %s", expected, actual)
	}

	expected = "Disabled"
	actual = IsEnabled(false)
	if actual != expected {
		t.Errorf("Expected %s. Actual %s", expected, actual)
	}
}

func TestIsValidCryptoAddress(t *testing.T) {
	t.Parallel()
	b, err := IsValidCryptoAddress("1Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX", "bTC")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if !b {
		t.Errorf("expected address '%s' to be valid", "1Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX")
	}

	b, err = IsValidCryptoAddress("bc1qw508d6qejxtdg4y5r3zarvaly0c5xw7kv8f3t4", "bTC")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if !b {
		t.Errorf("expected address '%s' to be valid", "bc1qw508d6qejxtdg4y5r3zarvaly0c5xw7kv8f3t4")
	}

	b, err = IsValidCryptoAddress("an84characterslonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1569pvx", "bTC")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "an84characterslonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1569pvx")
	}

	b, err = IsValidCryptoAddress("bc1qc7slrfxkknqcq2jevvvkdgvrt8080852dfjewde450xdlk4ugp7szw5tk9", "bTC")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if !b {
		t.Errorf("expected address '%s' to be valid", "bc1qc7slrfxkknqcq2jevvvkdgvrt8080852dfjewde450xdlk4ugp7szw5tk9")
	}

	b, err = IsValidCryptoAddress("0Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX", "btc")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "0Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX")
	}

	b, err = IsValidCryptoAddress("1Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX", "lTc")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "1Mz7153HMuxXTuR2R1t78mGSdzaAtNbBWX")
	}

	b, err = IsValidCryptoAddress("3CDJNfdWX8m2NwuGUV3nhXHXEeLygMXoAj", "ltc")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if !b {
		t.Errorf("expected address '%s' to be valid", "3CDJNfdWX8m2NwuGUV3nhXHXEeLygMXoAj")
	}

	b, err = IsValidCryptoAddress("NCDJNfdWX8m2NwuGUV3nhXHXEeLygMXoAj", "lTc")
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "NCDJNfdWX8m2NwuGUV3nhXHXEeLygMXoAj")
	}

	b, err = IsValidCryptoAddress(
		"0xb794f5ea0ba39494ce839613fffba74279579268",
		"eth",
	)
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if !b {
		t.Errorf("expected address '%s' to be valid", "0xb794f5ea0ba39494ce839613fffba74279579268")
	}

	b, err = IsValidCryptoAddress(
		"xxb794f5ea0ba39494ce839613fffba74279579268",
		"eTh",
	)
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "xxb794f5ea0ba39494ce839613fffba74279579268")
	}

	b, err = IsValidCryptoAddress(
		"xxb794f5ea0ba39494ce839613fffba74279579268",
		"ding",
	)
	if !errors.Is(err, errInvalidCryptoCurrency) {
		t.Errorf("received '%v' expected '%v'", err, errInvalidCryptoCurrency)
	}
	if b {
		t.Errorf("expected address '%s' to be invalid", "xxb794f5ea0ba39494ce839613fffba74279579268")
	}
}

func TestSliceDifference(t *testing.T) {
	t.Parallel()

	assert.ElementsMatch(t, []string{"world", "go"}, SliceDifference([]string{"hello", "world"}, []string{"hello", "go"}))
	assert.ElementsMatch(t, []int64{1, 2, 5, 6}, SliceDifference([]int64{1, 2, 3, 4}, []int64{3, 4, 5, 6}))
	assert.ElementsMatch(t, []float64{1.1, 4.4}, SliceDifference([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3, 4.4}))
	type mixedType struct {
		A string
		B int
	}
	assert.ElementsMatch(t, []mixedType{{"A", 1}, {"D", 4}}, SliceDifference([]mixedType{{"A", 1}, {"B", 2}, {"C", 3}}, []mixedType{{"B", 2}, {"C", 3}, {"D", 4}}))
	assert.ElementsMatch(t, []int{1, 2, 3}, SliceDifference([]int{}, []int{1, 2, 3}))
	assert.ElementsMatch(t, []int{1, 2, 3}, SliceDifference([]int{1, 2, 3}, []int{}))
	assert.Empty(t, SliceDifference([]int{}, []int{}))
}

func TestStringSliceContains(t *testing.T) {
	t.Parallel()
	originalHaystack := []string{"hello", "world", "USDT", "Contains", "string"}
	assert.True(t, StringSliceContains(originalHaystack, "USD"), "Should contain 'USD'")
	assert.False(t, StringSliceContains(originalHaystack, "thing"), "Should not contain 'thing'")
}

func TestStringSliceCompareInsensitive(t *testing.T) {
	t.Parallel()
	originalHaystack := []string{"hello", "WoRld", "USDT", "Contains", "string"}
	assert.False(t, StringSliceCompareInsensitive(originalHaystack, "USD"), "Should not contain 'USD'")
	assert.True(t, StringSliceCompareInsensitive(originalHaystack, "WORLD"), "Should find 'WoRld'")
}

func TestStringSliceContainsInsensitive(t *testing.T) {
	t.Parallel()
	originalHaystack := []string{"bLa", "BrO", "sUp"}
	assert.True(t, StringSliceContainsInsensitive(originalHaystack, "Bla"), "Should contain 'Bla'")
	assert.False(t, StringSliceContainsInsensitive(originalHaystack, "ning"), "Should not contain 'ning'")
}

func TestYesOrNo(t *testing.T) {
	t.Parallel()
	if !YesOrNo("y") {
		t.Error("Common YesOrNo Error.")
	}
	if !YesOrNo("yes") {
		t.Error("Common YesOrNo Error.")
	}
	if YesOrNo("ding") {
		t.Error("Common YesOrNo Error.")
	}
}

func TestEncodeURLValues(t *testing.T) {
	t.Parallel()
	urlstring := "https://www.test.com"
	expectedOutput := `https://www.test.com?env=TEST%2FDATABASE&format=json`
	values := url.Values{}
	values.Set("format", "json")
	values.Set("env", "TEST/DATABASE")

	output := EncodeURLValues(urlstring, values)
	if output != expectedOutput {
		t.Error("common EncodeURLValues error")
	}
}

func TestExtractHost(t *testing.T) {
	t.Parallel()
	address := "localhost:1337"
	addresstwo := ":1337"
	expectedOutput := "localhost"
	actualResult := ExtractHost(address)
	if expectedOutput != actualResult {
		t.Errorf(
			"Expected '%s'. Actual '%s'.", expectedOutput, actualResult)
	}
	actualResultTwo := ExtractHost(addresstwo)
	if expectedOutput != actualResultTwo {
		t.Errorf(
			"Expected '%s'. Actual '%s'.", expectedOutput, actualResult)
	}

	address = "192.168.1.100:1337"
	expectedOutput = "192.168.1.100"
	actualResult = ExtractHost(address)
	if expectedOutput != actualResult {
		t.Errorf(
			"Expected '%s'. Actual '%s'.", expectedOutput, actualResult)
	}
}

func TestExtractPort(t *testing.T) {
	t.Parallel()
	address := "localhost:1337"
	expectedOutput := 1337
	actualResult := ExtractPort(address)
	if expectedOutput != actualResult {
		t.Errorf(
			"Expected '%d'. Actual '%d'.", expectedOutput, actualResult)
	}

	address = "localhost"
	expectedOutput = 80
	actualResult = ExtractPort(address)
	if expectedOutput != actualResult {
		t.Errorf(
			"Expected '%d'. Actual '%d'.", expectedOutput, actualResult)
	}
}

func TestGetURIPath(t *testing.T) {
	t.Parallel()
	// mapping of input vs expected result
	testTable := map[string]string{
		"https://api.pro.coinbase.com/accounts":         "/accounts",
		"https://api.pro.coinbase.com/accounts?a=1&b=2": "/accounts?a=1&b=2",
		"http://www.google.com/accounts?!@#$%;^^":       "",
	}
	for testInput, expectedOutput := range testTable {
		actualOutput := GetURIPath(testInput)
		if actualOutput != expectedOutput {
			t.Errorf("Expected '%s'. Actual '%s'.",
				expectedOutput, actualOutput)
		}
	}
}

func TestGetExecutablePath(t *testing.T) {
	t.Parallel()
	if _, err := GetExecutablePath(); err != nil {
		t.Errorf("Common GetExecutablePath. Error: %s", err)
	}
}

func TestGetDefaultDataDir(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		dir, ok := os.LookupEnv("APPDATA")
		if !ok {
			t.Fatal("APPDATA is not set")
		}
		dir = filepath.Join(dir, "GoCryptoTrader")
		actualOutput := GetDefaultDataDir(runtime.GOOS)
		if actualOutput != dir {
			t.Fatalf("Unexpected result. Got: %v Expected: %v", actualOutput, dir)
		}
	default:
		var dir string
		usr, err := user.Current()
		if err == nil {
			dir = usr.HomeDir
		} else {
			var err error
			dir, err = os.UserHomeDir()
			if err != nil {
				dir = "."
			}
		}
		dir = filepath.Join(dir, ".gocryptotrader")
		actualOutput := GetDefaultDataDir(runtime.GOOS)
		if actualOutput != dir {
			t.Fatalf("Unexpected result. Got: %v Expected: %v", actualOutput, dir)
		}
	}
}

func TestCreateDir(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		// test for looking up an invalid directory
		err := CreateDir("")
		if err == nil {
			t.Fatal("expected err due to invalid path, but got nil")
		}

		// test for a directory that exists
		dir, ok := os.LookupEnv("TEMP")
		if !ok {
			t.Fatal("LookupEnv failed. TEMP is not set")
		}
		err = CreateDir(dir)
		if err != nil {
			t.Fatalf("CreateDir failed. Err: %v", err)
		}

		// test for creating a directory
		dir, ok = os.LookupEnv("APPDATA")
		if !ok {
			t.Fatal("LookupEnv failed. APPDATA is not set")
		}
		dir = filepath.Join(dir, "GoCryptoTrader", "TestFileASDFG")
		err = CreateDir(dir)
		if err != nil {
			t.Fatalf("CreateDir failed. Err: %v", err)
		}
		err = os.Remove(dir)
		if err != nil {
			t.Fatalf("Failed to remove file. Err: %v", err)
		}
	default:
		err := CreateDir("")
		if err == nil {
			t.Fatal("expected err due to invalid path, but got nil")
		}

		dir := "/home"
		err = CreateDir(dir)
		if err != nil {
			t.Fatalf("CreateDir failed. Err: %v", err)
		}
		var ok bool
		dir, ok = os.LookupEnv("HOME")
		if !ok {
			t.Fatal("LookupEnv of HOME failed")
		}
		dir = filepath.Join(dir, ".gocryptotrader", "TestFileASFG")
		err = CreateDir(dir)
		if err != nil {
			t.Errorf("CreateDir failed. Err: %s", err)
		}
		err = os.Remove(dir)
		if err != nil {
			t.Fatalf("Failed to remove file. Err: %v", err)
		}
	}
}

func TestChangePermission(t *testing.T) {
	t.Parallel()
	testDir := filepath.Join(os.TempDir(), "TestFileASDFGHJ")
	switch runtime.GOOS {
	case "windows":
		err := ChangePermission("*")
		if err == nil {
			t.Fatal("expected an error on non-existent path")
		}
		err = os.Mkdir(testDir, 0o777)
		if err != nil {
			t.Fatalf("Mkdir failed. Err: %v", err)
		}
		err = ChangePermission(testDir)
		if err != nil {
			t.Fatalf("ChangePerm was unsuccessful. Err: %v", err)
		}
		_, err = os.Stat(testDir)
		if err != nil {
			t.Fatalf("os.Stat failed. Err: %v", err)
		}
		err = os.Remove(testDir)
		if err != nil {
			t.Fatalf("os.Remove failed. Err: %v", err)
		}
	default:
		err := ChangePermission("")
		if err == nil {
			t.Fatal("expected an error on non-existent path")
		}
		err = os.Mkdir(testDir, 0o777)
		if err != nil {
			t.Fatalf("Mkdir failed. Err: %v", err)
		}
		err = ChangePermission(testDir)
		if err != nil {
			t.Fatalf("ChangePerm was unsuccessful. Err: %v", err)
		}
		var a os.FileInfo
		a, err = os.Stat(testDir)
		if err != nil {
			t.Fatalf("os.Stat failed. Err: %v", err)
		}
		if a.Mode().Perm() != file.DefaultPermissionOctal {
			t.Fatalf("expected file permissions differ. expecting file.DefaultPermissionOctal got %#o", a.Mode().Perm())
		}
		err = os.Remove(testDir)
		if err != nil {
			t.Fatalf("os.Remove failed. Err: %v", err)
		}
	}
}

func TestAddPaddingOnUpperCase(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Supplied string
		Expected string
	}{
		{
			// empty
		},
		{
			Supplied: "ExpectedHTTPRainbow",
			Expected: "Expected HTTP Rainbow",
		},
		{
			Supplied: "SmellyCatSmellsBad",
			Expected: "Smelly Cat Smells Bad",
		},
		{
			Supplied: "Gronk",
			Expected: "Gronk",
		},
	}

	for x := range testCases {
		if received := AddPaddingOnUpperCase(testCases[x].Supplied); received != testCases[x].Expected {
			t.Fatalf("received '%v' but expected '%v'", received, testCases[x].Expected)
		}
	}
}

func TestErrors(t *testing.T) {
	t.Parallel()

	e1 := errors.New("inconsistent gravity")
	e2 := errors.New("barely marginal interest in your story")
	e3 := errors.New("error making dinner")
	e4 := errors.New("inconsistent gravy")
	e5 := errors.New("add vodka")

	// Nil tests
	assert.NoError(t, AppendError(nil, nil), "Append nil to nil should nil")
	assert.Same(t, AppendError(e1, nil), e1, "Append nil to e1 should e1")
	assert.Same(t, AppendError(nil, e2), e2, "Append e2 to nil should e2")

	// Vanila error tests
	err := AppendError(AppendError(AppendError(nil, e1), e2), e1)
	assert.ErrorContains(t, err, "inconsistent gravity, barely marginal interest in your story, inconsistent gravity", "Should format consistently")
	assert.ErrorIs(t, err, e1, "Should have inconsistent gravity")
	assert.ErrorIs(t, err, e2, "Should be bored by your witty tales")

	err = ExcludeError(err, e2)
	assert.ErrorIs(t, err, e1, "Should still be bored")
	assert.NotErrorIs(t, err, e2, "Should not be an e2")
	me, ok := err.(*multiError)
	if assert.True(t, ok, "Should be a multiError") {
		assert.Len(t, me.errs, 2, "Should only have 2 errors")
	}
	err = ExcludeError(err, e1)
	assert.NoError(t, err, "Error should be empty")
	err = ExcludeError(err, e1)
	assert.NoError(t, err, "Excluding a nil error should be okay")

	// Wrapped error tests
	err = fmt.Errorf("%w: %w", e3, fmt.Errorf("%w: %w", e4, e5))
	assert.ErrorIs(t, ExcludeError(err, e4), e3, "Excluding e4 should retain e3")
	assert.ErrorIs(t, ExcludeError(err, e4), e5, "Excluding e4 should retain the vanilla co-wrapped e5")
	assert.NotErrorIs(t, ExcludeError(err, e4), e4, "e4 should be excluded")
	assert.ErrorIs(t, ExcludeError(err, e5), e3, "Excluding e5 should retain e3")
	assert.ErrorIs(t, ExcludeError(err, e5), e4, "Excluding e5 should retain the vanilla co-wrapped e4")
	assert.NotErrorIs(t, ExcludeError(err, e5), e5, "e5 should be excluded")

	// Hybrid tests
	err = AppendError(fmt.Errorf("%w: %w", e4, e5), e3)
	assert.ErrorIs(t, ExcludeError(err, e4), e3, "Excluding e4 should retain e3")
	assert.ErrorIs(t, ExcludeError(err, e4), e5, "Excluding e4 should retain the vanilla co-wrapped e5")
	assert.NotErrorIs(t, ExcludeError(err, e4), e4, "e4 should be excluded")
	assert.ErrorIs(t, ExcludeError(err, e5), e3, "Excluding e5 should retain e3")
	assert.ErrorIs(t, ExcludeError(err, e5), e4, "Excluding e5 should retain the vanilla co-wrapped e4")
	assert.NotErrorIs(t, ExcludeError(err, e5), e5, "e4 should be excluded")

	// Formatting retention
	err = AppendError(e1, fmt.Errorf("%w: Run out of `%s`: %w", e3, "sausages", e5))
	assert.ErrorIs(t, err, e1, "Should be an e1")
	assert.ErrorIs(t, err, e3, "Should be an e3")
	assert.ErrorIs(t, err, e5, "Should be an e5")
	assert.ErrorContains(t, err, "sausages", "Should know about secret sausages")
}

func TestParseStartEndDate(t *testing.T) {
	t.Parallel()
	pt := time.Date(1999, 1, 1, 0, 0, 0, 0, time.Local)
	ft := time.Date(2222, 1, 1, 0, 0, 0, 0, time.Local)
	et := time.Date(2020, 1, 1, 1, 0, 0, 0, time.Local)
	nt := time.Time{}

	err := StartEndTimeCheck(nt, nt)
	if !errors.Is(err, ErrDateUnset) {
		t.Errorf("received %v, expected %v", err, ErrDateUnset)
	}

	err = StartEndTimeCheck(et, nt)
	if !errors.Is(err, ErrDateUnset) {
		t.Errorf("received %v, expected %v", err, ErrDateUnset)
	}

	err = StartEndTimeCheck(et, zeroValueUnix)
	if !errors.Is(err, ErrDateUnset) {
		t.Errorf("received %v, expected %v", err, ErrDateUnset)
	}

	err = StartEndTimeCheck(zeroValueUnix, et)
	if !errors.Is(err, ErrDateUnset) {
		t.Errorf("received %v, expected %v", err, ErrDateUnset)
	}

	err = StartEndTimeCheck(et, et)
	if !errors.Is(err, ErrStartEqualsEnd) {
		t.Errorf("received %v, expected %v", err, ErrStartEqualsEnd)
	}

	err = StartEndTimeCheck(et, pt)
	if !errors.Is(err, ErrStartAfterEnd) {
		t.Errorf("received %v, expected %v", err, ErrStartAfterEnd)
	}

	err = StartEndTimeCheck(ft, ft.Add(time.Hour))
	if !errors.Is(err, ErrStartAfterTimeNow) {
		t.Errorf("received %v, expected %v", err, ErrStartAfterTimeNow)
	}

	err = StartEndTimeCheck(pt, et)
	if !errors.Is(err, nil) {
		t.Errorf("received %v, expected %v", err, nil)
	}
}

func TestGetAssertError(t *testing.T) {
	err := GetTypeAssertError("*[]string", float64(0))
	if err.Error() != "type assert failure from float64 to *[]string" {
		t.Fatal(err)
	}

	err = GetTypeAssertError("<nil>", nil)
	if err.Error() != "type assert failure from <nil> to <nil>" {
		t.Fatal(err)
	}

	err = GetTypeAssertError("bruh", struct{}{})
	if !errors.Is(err, ErrTypeAssertFailure) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrTypeAssertFailure)
	}

	err = GetTypeAssertError("string", struct{}{})
	if err.Error() != "type assert failure from struct {} to string" {
		t.Errorf("unexpected error message: %v", err)
	}

	err = GetTypeAssertError("string", struct{}{}, "bidSize")
	if err.Error() != "type assert failure from struct {} to string for: bidSize" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMatchesEmailPattern(t *testing.T) {
	success := MatchesEmailPattern("someone semail")
	if success {
		t.Error("MatchesEmailPattern() unexpected test validation result")
	}
	success = MatchesEmailPattern("someone esemail@gmail")
	if success {
		t.Error("MatchesEmailPattern() unexpected test validation result")
	}
	success = MatchesEmailPattern("123@gmail")
	if !success {
		t.Error("MatchesEmailPattern() unexpected test validation result")
	}
	success = MatchesEmailPattern("someonesemail@email.com")
	if !success {
		t.Error("MatchesEmailPattern() unexpected test validation result")
	}
}

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()
	sample, err := GenerateRandomString(5, NumberCharacters)
	if err != nil {
		t.Errorf("GenerateRandomString()  %v", err)
	}
	value, err := strconv.Atoi(sample)
	if len(sample) != 5 || err != nil || value < 0 {
		t.Error("GenerateRandomString() unexpected test validation result")
	}
	sample, err = GenerateRandomString(5)
	if err != nil {
		t.Errorf("GenerateRandomString()  %v", err)
	}
	values, err := strconv.ParseInt(sample, 10, 64)
	if len(sample) != 5 || err != nil || values < 0 {
		t.Error("GenerateRandomString() unexpected test validation result")
	}
	_, err = GenerateRandomString(1, "")
	if err == nil {
		t.Errorf("GenerateRandomString() expecting %s, but found %v", "invalid characters, character must not be empty", err)
	}
	sample, err = GenerateRandomString(0, "")
	if err != nil && !strings.Contains(err.Error(), "invalid length") {
		t.Errorf("GenerateRandomString()  %v", err)
	}
	if sample != "" {
		t.Error("GenerateRandomString() unexpected test validation result")
	}
}

// TestErrorCollector exercises the error collector
func TestErrorCollector(t *testing.T) {
	e := CollectErrors(4)
	for i := range 4 {
		go func() {
			if i%2 == 0 {
				e.C <- errors.New("Collected error")
			} else {
				e.C <- nil
			}
			e.Wg.Done()
		}()
	}
	v := e.Collect()
	errs, ok := v.(*multiError)
	require.True(t, ok, "Must return a multiError")
	assert.Len(t, errs.Unwrap(), 2, "Should have 2 errors")
}

// TestBatch ensures the Batch function does not regress into common behavioural faults if implementation changes
func TestBatch(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b := Batch(s, 3)
	require.Len(t, b, 4)
	assert.Len(t, b[0], 3)
	assert.Len(t, b[3], 1)

	b[0][0] = 42
	assert.Equal(t, 1, s[0], "Changing the batches must not change the source")

	require.NotPanics(t, func() { Batch(s, -1) }, "Must not panic on negative batch size")
	done := make(chan any, 1)
	go func() { done <- Batch(s, 0) }()
	require.Eventually(t, func() bool { return len(done) > 0 }, time.Second, time.Millisecond, "Batch 0 must not hang")

	for _, i := range []int{-1, 0, 50} {
		b = Batch(s, i)
		require.Lenf(t, b, 1, "A batch size of %v should produce a single batch", i)
		assert.Lenf(t, b[0], len(s), "A batch size of %v should produce a single batch", i)
	}
}

type A int

func (a A) String() string {
	return strconv.Itoa(int(a))
}

func TestSortStrings(t *testing.T) {
	assert.Equal(t, []A{1, 2, 5, 6}, SortStrings([]A{6, 2, 5, 1}))
}

func TestCounter(t *testing.T) {
	t.Parallel()
	c := Counter{n: -5}
	require.Equal(t, int64(1), c.IncrementAndGet())
	require.Equal(t, int64(2), c.IncrementAndGet())
}

// 683185328	         1.787 ns/op	       0 B/op	       0 allocs/op
func BenchmarkCounter(b *testing.B) {
	c := Counter{}
	for i := 0; i < b.N; i++ {
		c.IncrementAndGet()
	}
}
