package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net"
	"net/http"
	"net/http/fcgi"
	"fmt"
	"time"
	"html/template"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
)

// セッション名
var session_name string = "gsid"

// Cookie型のstore情報
var store *sessions.CookieStore

// セッションオブジェクト
var session *sessions.Session

// 構造体
type Data1 struct {
	Count    int
	Msg      string
}

// 主処理
func main(){

	// 構造体を登録
	gob.Register(&Data1{})

	// セッション初期処理
	sessionInit()

	// ルーティング生成
	r := mux.NewRouter()

	// URL別の処理
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		// セッションオブジェクトを取得
		session, _ := store.Get(r, session_name)

		data1, ok := session.Values["data1"].(*Data1)
		if data1 != nil {
			data1.Count++
			data1.Msg = fmt.Sprintf("%d件カウント", data1.Count)
		} else {
			data1 = &Data1{0, "データ無し"}
		}
		fmt.Println(ok)
		fmt.Println(data1)
		session.Values["data1"] = data1

		// 保存
		sessions.Save(r, w)

		// テンプレートを指定
		tmpl := template.Must(template.New("index").ParseFiles("./index.html"))
		tmpl.Execute(w, struct {
			Detail *Data1
		}{
			Detail: data1,
		})

		// logの代わり
		fmt.Print(time.Now())
		fmt.Println(" url = " + r.URL.Path)
	})

	r.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request){

		// セッション初期化
		sessionInit()

		// logの代わり
		fmt.Print(time.Now())
		fmt.Println(" url = " + r.URL.Path)

		// redirect
		http.Redirect(w, r, "/", http.StatusFound)

	})

	// rを割当
	http.Handle("/", r)

	// ポートを割当
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
       		return
   	}

   	fcgi.Serve(l, nil)
}

// セッション用の初期処理
func sessionInit() {

	// 乱数生成
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	str := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	// 新しいstoreとセッションを準備
	store = sessions.NewCookieStore([]byte(str))
	session = sessions.NewSession(store, session_name)

	// セッションの有効範囲を指定
	store.Options = &sessions.Options{
		Domain:   "127.0.0.1",
		Path:     "/",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}

	// log
	fmt.Println("key     data --")
	fmt.Println(str)
	fmt.Println("")
}
