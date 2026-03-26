package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cod3rboy/robo-walker/render"
	"github.com/cod3rboy/robo-walker/robo"
)

const (
	DefaultWidth    string = "100"
	DefaultHeight   string = "100"
	DefaultDelay    string = "0.1"
	DefaultBGColor  string = "000000FF"
	DefaultFGColor  string = "00FF00FF"
	DefaultPosColor string = "FF0000FF"
)

const (
	QueryWidth    string = "w"
	QueryHeight   string = "h"
	QueryDelay    string = "d"
	QueryBGColor  string = "bgc"
	QueryFGColor  string = "fgc"
	QueryPosColor string = "posc"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{program}", executeProgram)
	serverAddress := ":8880"
	server := http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	log.Printf("server listening at %s\n", serverAddress)
	if err := server.ListenAndServe(); err != nil {
		log.Println("server stopped!")
		log.Println(err)
	}
}

func executeProgram(w http.ResponseWriter, r *http.Request) {
	program := r.PathValue("program")
	w.Header().Set("Content-Type", "text/html")
	program = strings.TrimSpace(program)

	customOpts, err := getCustomOptsFromQuery(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<p style=\"color:red;font-size:32\">%s</p>", err.Error())
		return
	}

	gif, err := executeRender(program, customOpts)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p style=\"color:red;font-size:32\">%s</p>", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, getEmbeddableGIFImageTag(gif, customOpts.Width, customOpts.Height))
}

type CustomOptions struct {
	Width  int
	Height int
	FColor uint32
	BColor uint32
	PColor uint32
	Delay  float64
}

func executeRender(program string, opts CustomOptions) ([]byte, error) {
	w := robo.NewWorld(opts.Width, opts.Height)
	if err := w.Run(program); err != nil {
		return nil, err
	}

	renderer := render.NewGIFRendererForSnapshots(
		w.Snapshots(),
		int(math.Ceil(100*opts.Delay)), // 100 = 1sec delay for next frame
	)
	renderOpts := render.DrawOpts{
		FgColor:  render.Color(opts.FColor),
		BgColor:  render.Color(opts.BColor),
		PosColor: render.Color(opts.PColor),
	}

	return renderer.Render(renderOpts)
}

func getEmbeddableGIFImageTag(gif []byte, width, height int) string {
	gifEncoded := base64.StdEncoding.EncodeToString(gif)
	return fmt.Sprintf("<img src=\"data:image/gif;base64,%s\" alt=\"robot walker program\" width=\"{%d}\" height=\"{%d}\" />", gifEncoded, width, height)
}

func getCustomOptsFromQuery(query url.Values) (CustomOptions, error) {
	opts := CustomOptions{}

	delay, err := strconv.ParseFloat(getQueryOrDefault(query, QueryDelay, DefaultDelay), 64)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryDelay)
	}
	opts.Delay = delay

	width, err := strconv.Atoi(getQueryOrDefault(query, QueryWidth, DefaultWidth))
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryWidth)
	}
	opts.Width = width

	height, err := strconv.Atoi(getQueryOrDefault(query, QueryHeight, DefaultHeight))
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryHeight)
	}
	opts.Height = height

	fgColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryFGColor, DefaultFGColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryFGColor)
	}
	opts.FColor = uint32(fgColor)

	bgColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryBGColor, DefaultBGColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryBGColor)
	}
	opts.BColor = uint32(bgColor)

	posColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryPosColor, DefaultPosColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryPosColor)
	}
	opts.PColor = uint32(posColor)

	return opts, nil
}

func getQueryOrDefault(query url.Values, queryKey string, fallback string) string {
	if query.Has(queryKey) {
		return query.Get(queryKey)
	}
	return fallback
}
