package gogen

import (
	_ "embed"
	"strings"

	"github.com/hduhelp/go-zero/tools/goctl/api/spec"
	"github.com/hduhelp/go-zero/tools/goctl/config"
	"github.com/hduhelp/go-zero/tools/goctl/util/format"
)

//go:embed middleware.tpl
var middlewareImplementCode string

func genMiddleware(dir string, cfg *config.Config, api *spec.ApiSpec, useGin bool) error {
	middlewares := getMiddleware(api)
	for _, item := range middlewares {
		middlewareFilename := strings.TrimSuffix(strings.ToLower(item), "middleware") + "_middleware"
		filename, err := format.FileNamingFormat(cfg.NamingFormat, middlewareFilename)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		err = genFile(fileGenConfig{
			dir:             dir,
			subdir:          middlewareDir,
			filename:        filename + ".go",
			templateName:    "contextTemplate",
			category:        category,
			templateFile:    middlewareImplementCodeFile,
			builtinTemplate: middlewareImplementCode,
			data: map[string]any{
				"name":   strings.Title(name),
				"useGin": useGin,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
