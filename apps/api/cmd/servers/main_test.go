package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/routes"
	pkg "forestwhisperGo/apps/api/pkg/logger"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func TestHealthAndCORS(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.ProxyFix())
	app.Use(pkg.LoggerWithWriter(io.Discard))
	app.Use(middleware.Cors())
	routes.SetupAllRoutes(app)
	app.Get("/", index)

	const origin = "https://example.com"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", origin)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != origin {
		t.Fatalf("allow-origin = %q, want %q", got, origin)
	}

	var body struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		Data map[string]string `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Code != http.StatusOK || body.Msg != "ok" || body.Data["status"] != "ok" {
		t.Fatalf("unexpected health response: %+v", body)
	}

	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", origin)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("OPTIONS / status = %d, want %d", res.StatusCode, http.StatusNoContent)
	}
}

func TestInstallerSelectionFromTOML(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	for _, test := range []struct {
		config string
		want   bool
	}{
		{config: "false", want: false},
		{config: "true", want: true},
	} {
		config := "[server]\ninstall = " + test.config + "\n"
		if err := os.WriteFile("config.toml", []byte(config), 0600); err != nil {
			t.Fatal(err)
		}
		if got := utils.IsInstalled(); got != test.want {
			t.Errorf("IsInstalled() = %t, want %t", got, test.want)
		}
	}
}
