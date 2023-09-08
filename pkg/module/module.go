package module

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var (
	defaultProxy = "http://proxy.golang.org"
)

type Module struct {
	name string

	proxy string
}

func New(name string) *Module {
	return &Module{
		name:  name,
		proxy: defaultProxy,
	}
}

func (m *Module) Versions() ([]string, error) {
	url := m.proxy + "/" + m.name + "/@v/list"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var (
		versions []string
		buf      = bufio.NewReader(resp.Body)
	)

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}

		versions = append(versions, strings.TrimSpace(line))
	}

	return versions, nil
}

func (m *Module) Latest() (string, error) {
	url := m.proxy + "/" + m.name + "/@latest"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var (
		module moduleType
	)

	err = json.Unmarshal(data, &module)
	if err != nil {
		return "", err
	}

	return module.Version, nil
}
