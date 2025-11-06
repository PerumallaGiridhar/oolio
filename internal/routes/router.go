package routes

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/PerumallaGiridhar/oolio/internal/response"
	"github.com/PerumallaGiridhar/oolio/internal/routes/order"
	"github.com/PerumallaGiridhar/oolio/internal/routes/product"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MemUsage(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	stats := map[string]string{
		"Alloc":      fmt.Sprintf("%v MiB", m.Alloc/1024/1024),
		"TotalAlloc": fmt.Sprintf("%v MiB", m.TotalAlloc/1024/1024),
		"Sys":        fmt.Sprintf("%v MiB", m.Sys/1024/1024),
		"NumGC":      fmt.Sprintf("%v", m.NumGC),
	}
	response.JSONResponse(w, http.StatusOK, stats)

}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/live"))
	r.Get("/stats", MemUsage)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Mount("/product", product.NewRouter())
		r.Mount("/order", order.NewRouter())
	})

	return r
}
