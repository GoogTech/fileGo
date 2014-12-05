package godw

import (
	"bufio"
	"fmt"
	z "github.com/nutzam/zgo"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 主
func main() {

	// 设置CPU核心数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 设置日志的结构
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	// -------------------------------------------------------- //

	http.Handle("/css/", http.FileServer(http.Dir("template")))

	http.Handle("/js/", http.FileServer(http.Dir("template")))

	http.Handle("/files/", http.FileServer(http.Dir("template")))

	http.Handle("/images/", http.FileServer(http.Dir("template")))

	// -------------------------------------------------------- //

	http.HandleFunc("/", index)

	http.HandleFunc("/addfile.go", addfile)

	http.HandleFunc("/rmfile.go", rmfile)

	http.HandleFunc("/filelist.go", filelist)

	http.HandleFunc("/upload.go", upload)

	http.HandleFunc("/download.go", download)

	// -------------------------------------------------------- //

	// 建立监听
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// 踢出错误
		log.Panic(err)
	}

}

// 上传文件接口
func upload(w http.ResponseWriter, r *http.Request) {

	// 加锁,写入
	if "POST" == r.Method {

		file, multi, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()

		// 判断文件是否存在
		if z.Exists(fmt.Sprintf("files/%s", multi.Filename)) {
			// 返回错误信息
			http.Error(w, fmt.Sprintf("WARN: [%s] file exists ...", multi.Filename), 500)
			return
		}

		f, err := os.Create(fmt.Sprintf("files/%s", multi.Filename))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	}

	// 重定向
	http.Redirect(w, r, "/filelist.go", http.StatusFound)

	// 返回
	return

}

// 下载文件接口
func download(w http.ResponseWriter, r *http.Request) {

	// 解析参数
	r.ParseForm()

	// 获取文件名称
	fname := z.Trim(r.FormValue("f"))

	// 添加头信息
	w.Header().Set("Content-Type", "multipart/form-data")

	// 添加头信息
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))

	// 判断安装包是否存在
	if !z.Exists(fmt.Sprintf("files/%s", fname)) {
		http.Error(w, fmt.Sprintf("WARN: [%s] file not exists ...", fname), 500)
		return
	}

	// 写入文件流
	z.FileRF(fmt.Sprintf("files/%s", fname), func(f *os.File) {
		_, err := io.Copy(w, bufio.NewReader(f))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	// 返回
	return

}

// 主页
func index(w http.ResponseWriter, r *http.Request) {

	// 解析参数
	r.ParseForm()

	// cookie
	//if _, err := r.Cookie("username"); err != nil {
	//
	if _, ok := r.Form["admin"]; ok {
		// cookie
		cookie := http.Cookie{Name: "username", Value: "admin", Expires: time.Now().Add(24 * time.Hour)}
		// cookie
		http.SetCookie(w, &cookie)
	}
	//
	if _, ok := r.Form["danoo"]; ok {
		// cookie
		cookie := http.Cookie{Name: "username", Value: "danoo", Expires: time.Now().Add(24 * time.Hour)}
		// cookie
		http.SetCookie(w, &cookie)
	}
	//}

	// 解析主页面
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		// 输出错误信息
		http.Error(w, err.Error(), 500)
		return
	}

	// 执行
	t.Execute(w, nil)

	// 返回
	return

}

// 添加文件
func addfile(w http.ResponseWriter, r *http.Request) {

	// cookie
	if _, err := r.Cookie("username"); err != nil {
		// 重定向
		http.Redirect(w, r, "/filelist.go", http.StatusFound)
		// 返回
		return
	}

	// 解析主页面
	t, err := template.ParseFiles("template/files/addfile.html")
	if err != nil {
		// 输出错误信息
		http.Error(w, err.Error(), 500)
		return
	}

	// 执行
	t.Execute(w, t)

	// 返回
	return

}

// 删除文件
func rmfile(w http.ResponseWriter, r *http.Request) {

	// cookie
	if cookie, err := r.Cookie("username"); err != nil {
		// 重定向
		http.Redirect(w, r, "/filelist.go", http.StatusFound)
		// 返回
		return
	} else {
		// cookie
		if cookie.Value != "admin" {
			// 重定向
			http.Redirect(w, r, "/filelist.go", http.StatusFound)
			// 返回
			return
		}
	}

	// 解析参数
	r.ParseForm()

	// 获取文件名称
	fname := z.Trim(r.FormValue("f"))

	// 判断安装包是否存在
	if z.Exists(fmt.Sprintf("files/%s", fname)) && !z.IsBlank(fname) {
		// 删除
		z.Fremove(fmt.Sprintf("files/%s", fname))
	}

	// 重定向
	http.Redirect(w, r, "/filelist.go", http.StatusFound)

	// 返回
	return

}

/*
	这里偷个懒
	应该将文件信息记录到数据库或者文件中
	我这个每次都去扫描，浪费资源
*/

type I struct {
	Id   int
	Name string
	Size string
	Date string
	Stat string
}

type D struct {
	// 文件列表
	Files []*I
}

// 构造
func NewD() *D {
	d := new(D)
	d.Files = make([]*I, 0)
	return d
}

// 文件列表
func filelist(w http.ResponseWriter, r *http.Request) {

	// 管理员
	var admin string

	// cookie
	if cookie, err := r.Cookie("username"); err == nil {
		// 权限
		if cookie.Value == "admin" {
			// 管理员
			admin = "admin"
		}
	}

	// 解析参数
	r.ParseForm()

	// 获取文件名称
	fname := z.Trim(r.FormValue("f"))

	// 创建返回对象
	d := NewD()

	// ID
	var id int

	// 遍历本地文件
	filepath.Walk("files", func(ph string, f os.FileInfo, err error) error {
		// 文件不存在
		if f == nil {
			return nil
		}
		// 跳过文件夹
		if f.IsDir() {
			return nil
		}
		// 判断文件是否存在
		if z.IsBlank(fname) {
			// 累加
			id++
			// 记录文件
			d.Files = append(d.Files, &I{id, f.Name(), fmt.Sprintf("%d", f.Size()), f.ModTime().String(), admin})
		} else {
			// 检查包含
			if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(fname)) {
				// 累加
				id++
				// 记录文件
				d.Files = append(d.Files, &I{id, f.Name(), fmt.Sprintf("%d", f.Size()), f.ModTime().String(), admin})
			}
		}
		// 返回
		return nil
	})

	// 解析主页面
	t, err := template.ParseFiles("template/files/filelist.html")
	if err != nil {
		// 输出错误信息
		http.Error(w, err.Error(), 500)
		return
	}

	// 执行
	t.Execute(w, d)

	// 返回
	return

}
