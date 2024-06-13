package main

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"imageconverter/frontend"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/h2non/bimg"
)

const MAX_FORM_SIZE = 32 << 20

var ctx = context.Background()

//go:embed assets/*
var assets embed.FS

func main() {
	// config := GetConfig()

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		frontend.HomePage().Render(ctx, w)
	})
	router.HandleFunc("POST /convert", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(MAX_FORM_SIZE); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		defer file.Close()
		t := r.FormValue("type")
		b64, err := Convert(file, t)
		if err != nil {
			fmt.Fprintf(w, "Something went wrong: %v", err.Error())
			return
		}
		frontend.ImagePreview(b64).Render(ctx, w)
	})
	router.Handle("GET /assets/", Gzip(http.FileServer(http.FS(assets))))

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Starting server at http://127.0.0.1:3000")
	server.ListenAndServe()
}
func Convert(file multipart.File, t string) (string, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	img := bimg.NewImage(bytes)
	bimgType := GetImageType(t)
	converted, err := img.Convert(bimgType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:image/%v;base64,%v", strings.ToLower(t), base64.StdEncoding.EncodeToString(converted)), nil
}

func GetImageType(t string) bimg.ImageType {
	switch t {
	case "PNG":
		return bimg.PNG
	case "WEBP":
		return bimg.WEBP
	case "JPEG":
		return bimg.JPEG
	case "AVIF":
		return bimg.AVIF
	case "TIFF":
		return bimg.TIFF
	default:
		return bimg.PNG
	}
}
