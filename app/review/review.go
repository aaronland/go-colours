// Command-line tool to generate an HTML page (and associated assets) to review the colour extraction
// for an image using one or more extruders and one or more palettes. The application will spawn a short-lived
// web server to serve the HTML review on a random port number and open its URI in the default browser.
package review

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aaronland/go-colours/extract"
	"github.com/sfomuseum/go-www-show"
)

//go:embed index.html
var index_html string

type TemplateVars struct {
	Images   []*extract.Image
	Palettes []string
}

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	extract_opts := &extract.ExtractOptions{
		ExtruderURIs: opts.ExtruderURIs,
		PaletteURIs:  opts.PaletteURIs,
		Root:         opts.Root,
		Images:       opts.Images,
		AllowRemote:  true,
		CloneImages:  true,
	}

	extract_rsp, err := extract.Extract(ctx, extract_opts)

	if err != nil {
		return fmt.Errorf("Failed to extract images, %w", err)
	}

	if extract_rsp.IsTmpRoot {
		defer os.RemoveAll(extract_rsp.Root)
	}

	//

	index_t, err := template.New("index").Parse(index_html)

	if err != nil {
		return fmt.Errorf("Failed to parse index template, %w", err)
	}

	index_path := filepath.Join(extract_rsp.Root, "index.html")

	index_wr, err := os.OpenFile(index_path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open %s for writing, %w", index_path, err)
	}

	vars := TemplateVars{
		Images:   extract_rsp.Images,
		Palettes: extract_rsp.Palettes,
	}

	err = index_t.Execute(index_wr, vars)

	if err != nil {
		return fmt.Errorf("Failed to encode %s, %w", index_path, err)
	}

	err = index_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close %s after writing, %w", index_path, err)
	}

	//

	mux := http.NewServeMux()

	dir_fs := os.DirFS(extract_rsp.Root)
	http_fs := http.FileServerFS(dir_fs)

	mux.Handle("/", http_fs)

	browser, _ := show.NewBrowser(ctx, "web://")

	show_opts := &show.RunOptions{
		Browser: browser,
		Mux:     mux,
	}

	err = show.RunWithOptions(ctx, show_opts)

	if err != nil {
		return fmt.Errorf("Failed to show results, %w", err)
	}

	return nil
}
