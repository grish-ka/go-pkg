package main

import (
	"fmt"
	"net/http"
	"strings"
)

func tableHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Tell the browser to render HTML
	w.Header().Set("Content-Type", "text/html")

	// 2. Extract URL Parameters
	data := r.URL.Query().Get("data")
	projectTitle := r.URL.Query().Get("project")
	if projectTitle == "" {
		projectTitle = "Package Status"
	}

	// 3. The Style "Package" (CSS)
	// We put this in a raw string to keep it clean.
	fmt.Fprint(w, `
<style>
    body { 
        font-family: -apple-system, sans-serif; 
        margin: 0; padding: 0; 
        background-color: #f6f8fa; 
    }
    .header-bar {
        background-color: #24292f; /* Dark Repology style */
        color: white;
        padding: 20px;
        font-size: 22px;
        font-weight: bold;
        box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    }
    .container { padding: 30px; }
    table { 
        border-collapse: collapse; 
        width: 100%; 
        max-width: 600px; 
        background: white; 
        border: 1px solid #d0d7de;
        border-radius: 6px;
        overflow: hidden;
    }
    td { 
        padding: 12px; 
        border: 1px solid #d0d7de; 
        font-size: 14px;
    }
    .repo-name { 
        background-color: #f6f8fa; 
        font-weight: 600; 
        width: 40%; 
        color: #57606a;
    }
</style>
`)

	// 4. Render the Header "Kitty"
	fmt.Fprintf(w, "<div class='header-bar'>%s</div>\n", projectTitle)

	// 5. Start the Table Container
	fmt.Fprint(w, "<div class='container'>\n<table>\n")

	// 6. Logic Loop
	pairs := strings.Split(data, ",")
	for _, pair := range pairs {
		// Split "alpine:3.18-u" -> ["alpine", "3.18-u"]
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			manager := parts[0]

			// Nested Split for Version and State
			verState := strings.Split(parts[1], "-")
			version := verState[0]
			state := "n" // Default to Not Updated 🔴

			if len(verState) == 2 {
				state = verState[1] // Update to 'u' if present 🟢
			}

			// Color Decision
			bgColor := "#ffeef0" // Soft Red
			if state == "u" {
				bgColor = "#dafbe1" // Soft Green
			}

			// 7. Render the Row
			fmt.Fprintf(w, `<tr>
    <td class="repo-name">%s</td>
    <td style="background-color: %s;">%s</td>
</tr>
`, manager, bgColor, version)
		}
	}

	// 8. Close Tags
	fmt.Fprint(w, "</table>\n</div>")
}

func main() {
	// Traffic Cop (Mux)
	http.HandleFunc("/table", tableHandler)

	fmt.Println("Server online at http://localhost:8080/table")
	http.ListenAndServe(":8080", nil)
}