package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pei-pay/bookings_udemy/pkg/config"
	"github.com/pei-pay/bookings_udemy/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate  renders template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// cacheを使わない時、templateCacheをrebuildする
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	// TODO: byteに一回保存する意味は?
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

/*
全てのtemplateを探して、myCacheに保存する
*/
// CreateTemplateCache  creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	/*
		.page.tmplのfileを全て取得
	*/
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	/*
	*
	 */
	for _, page := range pages {

		// pageの名前を抽出
		name := filepath.Base(page)
		// templateSet(ts)を取り出す?
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// layout.tmplを探す
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// 取り出したtemplateSetを"./templates/*.layout.tmpl"に変換する
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		// myCacheに保存する
		/*
			['home.page.tmpl': templateSet, 'about.page.tmpl': templateSet]
		*/
		myCache[name] = ts
	}

	return myCache, nil
}
