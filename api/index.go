package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	project := r.URL.Query().Get("project")
	dataRaw := r.URL.Query().Get("data")
	
	if dataRaw == "" { return }
	pairs := strings.Split(dataRaw, ",")

	// Alphabetical Sort (Repology .rs style)
	sort.Slice(pairs, func(i, j int) bool {
		return strings.Split(pairs[i], ":")[0] < strings.Split(pairs[j], ":")[0]
	})

	const HeaderH = 28
	const RowH = 22
	const LabelW = 150
	const ValueW = 100
	const ColW = LabelW + ValueW

	numPairs := len(pairs)
	numCols := (numPairs + 32) / 33
	if numCols == 0 { numCols = 1 }

	limit := numPairs
	if limit > 33 { limit = 33 }
	
	totalH := HeaderH + (limit * RowH)
	totalW := numCols * ColW

	fmt.Fprintf(w, "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\">", totalW, totalH)
	fmt.Fprint(w, "<linearGradient id=\"g\" x2=\"0\" y2=\"100%\"><stop offset=\"0\" stop-color=\"#bbb\" stop-opacity=\".1\"/><stop offset=\"1\" stop-opacity=\".1\"/></linearGradient>")
	fmt.Fprintf(w, "<g font-family=\"Verdana,Geneva,sans-serif\"><rect width=\"100%%\" height=\"100%%\" fill=\"#555\" rx=\"4\"/>")

	// Title Header
	drawText(w, totalW/2, HeaderH/2, project, "middle", 13, true)

	for i, pair := range pairs {
		col, row := i/33, i%33
		x, y := col*ColW, HeaderH+(row*RowH)
		
		p := strings.Split(pair, ":")
		if len(p) < 2 { continue }

		mgr := p[0]
		valAndState := p[1]
		
		boxCol := "#e05d44"
		displayVer := valAndState
		
		if strings.HasSuffix(valAndState, "-u") {
			boxCol = "#4c1"
			displayVer = strings.TrimSuffix(valAndState, "-u")
		} else if strings.HasSuffix(valAndState, "-n") {
			boxCol = "#e05d44"
			displayVer = strings.TrimSuffix(valAndState, "-n")
		}

		fmt.Fprintf(w, "<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"%s\"/>", x+LabelW, y, ValueW, RowH, boxCol)
		fmt.Fprintf(w, "<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"url(#g)\"/>", x, y, ColW, RowH)
		drawText(w, x+12, y+(RowH/2), mgr, "start", 11, false)
		drawText(w, x+LabelW+(ValueW/2), y+(RowH/2), displayVer, "middle", 11, true)
	}
	fmt.Fprint(w, "</g></svg>")
}

func drawText(w http.ResponseWriter, x, y int, t, a string, s int, b bool) {
	fw := "normal"
	if b { fw = "bold" }
	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" text-anchor=\"%s\" font-size=\"%d\" font-weight=\"%s\" fill=\"#010101\" fill-opacity=\".3\" dominant-baseline=\"central\">%s</text>", x, y+1, a, s, fw, t)
	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" text-anchor=\"%s\" font-size=\"%d\" font-weight=\"%s\" fill=\"#fff\" dominant-baseline=\"central\">%s</text>", x, y, a, s, fw, t)
}