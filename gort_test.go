package gort

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGort(t *testing.T) {
	router := NewRouter()

	t.Run("Router", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/users/:id", func(ctx *Context) {
			ctx.WriteString("hello")
		})

		if router.routes.root == nil {
			t.Error("unexpected nil root")
			return
		}

		if router.routes.root.children == nil {
			t.Error("unexpected nil children")
			return
		}

		if len(router.routes.root.children) != 1 {
			t.Error("expected 1 child")
			return
		}

		if router.routes.root.children["users"] == nil {
			t.Error("unexpected nil child")
			return
		}

		route := router.routes.find("/users/foo")

		if route == nil {
			t.Error("unexpected nil route")
			return
		}
	})

	t.Run("Store", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/store/:key", func(ctx *Context) {
			value, ok := ctx.Store.Get(ctx.Params["key"])
			if ok {
				ctx.JSON(value)
				return
			}
			key := ctx.Params["key"]

			ctx.Store.Set(key, ctx.Request.RemoteAddr)
			ctx.JSON("ok")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/store/foo")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "\"ok\"" {
			t.Error("expected response body to be \"ok\"")
			return
		}

		value, ok := router.store.Get("foo")
		if !ok {
			t.Error("expected ok to be true")
			return
		}

		if value == nil {
			t.Error("unexpected nil value")
			return
		}
	})

	t.Run("Param", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/users/:id", func(ctx *Context) {
			ctx.WriteString(ctx.Param("id"))
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/users/foo")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "foo" {
			t.Error("expected response body to be foo")
			return
		}
	})

	t.Run("SetHeader", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/header", func(ctx *Context) {
			ctx.SetHeader("X-Test", "test")
			ctx.WriteString("hello")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/header")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		if res.Header.Get("X-Test") != "test" {
			t.Error("expected X-Test header to be test")
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})

	t.Run("SetHeaders", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/headers", func(ctx *Context) {
			ctx.SetHeaders(map[string]string{
				"X-Test": "test",
			})
			ctx.WriteString("hello")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/headers")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		if res.Header.Get("X-Test") != "test" {
			t.Error("expected X-Test header to be test")
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})

	t.Run("SetCookie", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/cookie", func(ctx *Context) {
			ctx.SetCookie(&http.Cookie{
				Name:  "test",
				Value: "test",
			})
			ctx.WriteString("hello")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/cookie")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		cookies := res.Cookies()

		if len(cookies) != 1 {
			t.Error("expected 1 cookie")
			return
		}

		if cookies[0].Name != "test" {
			t.Error("expected cookie name to be test")
			return
		}

		if cookies[0].Value != "test" {
			t.Error("expected cookie value to be test")
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})

	t.Run("SetStatus", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/status", func(ctx *Context) {
			ctx.SetStatus(http.StatusCreated)
			ctx.WriteString("hello")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/status")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("expected status code to be %d, got %d", http.StatusCreated, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})

	t.Run("GetHeader", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/get-header", func(ctx *Context) {
			ctx.WriteString(ctx.GetHeader("X-Test"))
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/get-header", nil)
		if err != nil {
			t.Error(err)
			return
		}

		req.Header.Set("X-Test", "test")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "test" {
			t.Error("expected response body to be test")
			return
		}
	})

	t.Run("Send", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/send", func(ctx *Context) {
			ctx.Send([]byte("hello"))
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/send")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})

	t.Run("Redirect", func(t *testing.T) {
		router.AddRoute(http.MethodGet, "/redirect", func(ctx *Context) {
			ctx.Redirect("/redirected")
		})

		router.AddRoute(http.MethodGet, "/redirected", func(ctx *Context) {
			ctx.SetStatus(http.StatusFound)
			ctx.SetHeader("Location", "/redirected")
			ctx.WriteString("hello")
		})

		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/redirect")
		if err != nil {
			t.Error(err)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusFound {
			t.Errorf("expected status code to be %d, got %d", http.StatusFound, res.StatusCode)
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(data) != "hello" {
			t.Error("expected response body to be hello")
			return
		}
	})
}

func BenchmarkRouter(b *testing.B) {
	router := NewRouter()

	users := map[string]string{
		"foo": "bar",
	}

	b.ResetTimer()
	router.AddRoute(http.MethodGet, "/users/:id", func(ctx *Context) {
		id := ctx.Param("id")
		user, ok := users[id]
		if !ok {
			ctx.NotFound()
			return
		}

		ctx.JSON(user)
	})

	ts := httptest.NewServer(router)

	for i := 0; i < b.N; i++ {
		http.Get(ts.URL + "/users/foo")
	}

	ts.Close()
}
