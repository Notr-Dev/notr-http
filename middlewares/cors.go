package middlewares

import notrhttp "github.com/Notr-Dev/notr-http"

func AllowALlOrigins(next notrhttp.Handler) notrhttp.Handler {
	return func(rw notrhttp.Writer, r *notrhttp.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		next(rw, r)
	}
}
