package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/minhdanh/vault-template/pkg/template"
)

type vaultConfig struct {
	token, endpoint string
}

type rendererConfig struct {
	vault vaultConfig
}

type renderer struct {
	vaultRenderer *template.VaultTemplateRenderer
}

func NewRenderer(cfg rendererConfig) (*renderer, error) {
	var vaultRenderer *template.VaultTemplateRenderer

	if cfg.vault.token != "" && cfg.vault.endpoint != "" {
		var err error
		vaultRenderer, err = template.NewVaultTemplateRenderer(cfg.vault.token, cfg.vault.endpoint)

		if err != nil {
			return nil, err
		}
	} else {
		panic("Error: Vault endpoint or token is incorrect.")
	}

	return &renderer{
		vaultRenderer: vaultRenderer,
	}, nil
}

func (r *renderer) renderSingleFile(inputFilePath, outputFilePath string) (err error) {
	inputContent, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		return
	}

	renderedContent := string(inputContent)
	if r.vaultRenderer != nil {
		renderedContent, err = r.vaultRenderer.RenderTemplate(string(inputContent))

		if err != nil {
			return
		}
	}

	if outputFilePath == "-" {
		fmt.Printf("%v", renderedContent)
		return
	}

	// make output path
	outputDirectory := filepath.Dir(outputFilePath)
	err = os.MkdirAll(outputDirectory, 0755)

	if err != nil {
		return
	}
	outputFile, err := os.Create(outputFilePath)

	if err != nil {
		return
	}

	defer func() {
		err = outputFile.Close()
	}()

	_, err = outputFile.Write([]byte(renderedContent))
	fmt.Println("Created file " + outputFilePath)
	return
}
