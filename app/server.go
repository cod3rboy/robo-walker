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
	DefaultImageSize string = "200"
	DefaultGridSize  string = "100"
	DefaultDelay     string = "0.1"
	DefaultBGColor   string = "000000FF"
	DefaultFGColor   string = "00FF00FF"
	DefaultPosColor  string = "FF0000FF"
)

const (
	QueryImageSize string = "imgs"
	QueryGridSize  string = "gs"
	QueryDelay     string = "d"
	QueryBGColor   string = "bgc"
	QueryFGColor   string = "fgc"
	QueryPosColor  string = "posc"
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
	fmt.Fprint(w, getEmbeddableGIFImageTag(gif, customOpts.ImageSize))
}

type CustomOptions struct {
	ImageSize       uint
	GridSize        uint
	ForeGroundColor uint32
	BackGroundColor uint32
	PositionColor   uint32
	Delay           float64
}

func executeRender(program string, opts CustomOptions) ([]byte, error) {
	w := robo.NewWorld(int(opts.GridSize))
	if err := w.Run(program); err != nil {
		return nil, err
	}

	renderer := render.NewGIFRendererForSnapshots(
		w.Snapshots(),
		int(math.Ceil(100*opts.Delay)), // 100 = 1sec delay for next frame
	)
	renderOpts := render.DrawOpts{
		FgColor:  render.Color(opts.ForeGroundColor),
		BgColor:  render.Color(opts.BackGroundColor),
		PosColor: render.Color(opts.PositionColor),
		Size:     opts.ImageSize,
	}

	return renderer.Render(renderOpts)
}

func getEmbeddableGIFImageTag(gif []byte, size uint) string {
	gifEncoded := base64.StdEncoding.EncodeToString(gif)
	return fmt.Sprintf("<img src=\"data:image/gif;base64,%s\" alt=\"robot walker program\" width=\"%d\" height=\"%d\" />", gifEncoded, size, size)
}

func getCustomOptsFromQuery(query url.Values) (CustomOptions, error) {
	opts := CustomOptions{}

	delay, err := strconv.ParseFloat(getQueryOrDefault(query, QueryDelay, DefaultDelay), 64)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryDelay)
	}
	opts.Delay = delay

	gridSize, err := strconv.ParseUint(getQueryOrDefault(query, QueryGridSize, DefaultGridSize), 10, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryGridSize)
	}
	opts.GridSize = uint(gridSize)

	imageSize, err := strconv.ParseUint(getQueryOrDefault(query, QueryImageSize, DefaultImageSize), 10, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryImageSize)
	}
	opts.ImageSize = uint(imageSize)

	fgColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryFGColor, DefaultFGColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryFGColor)
	}
	opts.ForeGroundColor = uint32(fgColor)

	bgColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryBGColor, DefaultBGColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryBGColor)
	}
	opts.BackGroundColor = uint32(bgColor)

	posColor, err := strconv.ParseUint(getQueryOrDefault(query, QueryPosColor, DefaultPosColor), 16, 32)
	if err != nil {
		return opts, fmt.Errorf("query param '%s' has invalid value", QueryPosColor)
	}
	opts.PositionColor = uint32(posColor)

	return opts, nil
}

func getQueryOrDefault(query url.Values, queryKey string, fallback string) string {
	if query.Has(queryKey) {
		return query.Get(queryKey)
	}
	return fallback
}
