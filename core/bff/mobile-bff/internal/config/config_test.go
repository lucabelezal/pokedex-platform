package config

import "testing"

func TestLoadConfigFavoritesViaCatalogDefaultTrue(t *testing.T) {
	t.Setenv("MOBILE_BFF_PORT", "")
	t.Setenv("POKEMON_CATALOG_SERVICE_URL", "http://localhost:8081")
	t.Setenv("JWT_SECRET", "segredo")
	t.Setenv("FAVORITES_VIA_CATALOG", "")

	cfg := LoadConfig()
	if !cfg.FavoritesViaCatalog {
		t.Error("FavoritesViaCatalog deveria ser true por padrão")
	}
}

func TestLoadConfigFavoritesViaCatalogFalse(t *testing.T) {
	t.Setenv("FAVORITES_VIA_CATALOG", "false")
	cfg := LoadConfig()
	if cfg.FavoritesViaCatalog {
		t.Error("FavoritesViaCatalog deveria ser false")
	}
}

func TestLoadConfigFavoritesViaCatalogVariants(t *testing.T) {
	for _, raw := range []string{"0", "no", "off", "FALSE", "False"} {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("FAVORITES_VIA_CATALOG", raw)
			cfg := LoadConfig()
			if cfg.FavoritesViaCatalog {
				t.Errorf("FAVORITES_VIA_CATALOG=%q deveria desabilitar", raw)
			}
		})
	}
	for _, raw := range []string{"true", "1", "yes", "on", "TRUE"} {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("FAVORITES_VIA_CATALOG", raw)
			cfg := LoadConfig()
			if !cfg.FavoritesViaCatalog {
				t.Errorf("FAVORITES_VIA_CATALOG=%q deveria habilitar", raw)
			}
		})
	}
}
