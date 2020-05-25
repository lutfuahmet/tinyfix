package tinypng

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Response from the TinyPNG API
type Response struct {
	Input   Input
	Output  Output
	Error   string
	Message string
	URL     string
}

// Input size
type Input struct {
	Size int32
}

// Output size, ratio and url
type Output struct {
	Size  int32
	Ratio float64
}

// PopulateFromHTTPResponse populates response based on HTTP response
func (r *Response) PopulateFromHTTPResponse(res *http.Response)(err error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}

	// Get the output URL from the Location header
	r.URL = res.Header.Get("Location")
	return
}

// SaveAs downloads and saves the compressed PNG file
func (r *Response) SaveAs(fn string) (err error) {
	resp, err := http.Get(r.URL)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	out, err := os.Create(fn)
	if err != nil {
		return
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}

// Print a line of statistics
func (r *Response) Print() {
	fmt.Print("Input size: ", r.Input.Size)
	fmt.Print(" Output size: ", r.Output.Size)
	fmt.Println(" Ratio:", r.Output.Ratio)
	fmt.Println("\n", r.URL)
}
