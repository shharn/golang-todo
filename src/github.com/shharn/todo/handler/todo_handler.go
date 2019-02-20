package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shharn/todo/model"
	"github.com/shharn/todo/util"
	"github.com/shharn/todo/service"
)

type IndexFileFinder interface {
	GetAsByteArray() ([]byte, error)
}

type LocalFileSystemFinder struct {
	path string
}

func (f LocalFileSystemFinder) GetAsByteArray() ([]byte, error) {
	data, err := ioutil.ReadFile(f.path)
	return data, err
}

func NewLocalFileSystemFinder(path string) *LocalFileSystemFinder {
	return &LocalFileSystemFinder{
		path: path,
	}
}

func WithFinder(handler func (http.ResponseWriter, *http.Request, IndexFileFinder), finder IndexFileFinder) http.HandlerFunc {
	return func (w http.ResponseWriter, rq *http.Request) {
		handler(w, rq, finder)
	}
}

func IndexHandler(w http.ResponseWriter, rq *http.Request, finder IndexFileFinder) {
	if rq.URL.Path != "/" {
		http.NotFound(w, rq)
		return;
	}

	data, err := finder.GetAsByteArray()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%+v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

var todoService = service.NewSimpleTodoService()

func TodosHandler(w http.ResponseWriter, rq *http.Request) {
	var errRecord error
	switch (rq.Method) {
	case "GET":
		queryParams := rq.URL.Query()
		pageNo := sanitizePageNo(queryParams["pageNo"])
		todos, totalCount := todoService.GetTodosByPageNo(pageNo)
		dto := model.TodosDto{Todos: todos, TotalCount: totalCount}
		if bytes, err := util.NewJSONSerializer().Serialize(dto); err != nil {
			errRecord = errors.WithStack(err)
		} else {
			w.Write(bytes)
		}
	case "POST":
		var todo model.Todo
		if err := util.NewJSONSerializer().Deserialize(rq.Body, &todo); err != nil {
			errRecord = errors.WithStack(err)
			break
		}
		if err := todoService.CreateTodo(todo); err != nil {
			errRecord = err
		}
	case "PATCH":
		var todo model.Todo
		if err := util.NewJSONSerializer().Deserialize(rq.Body, &todo); err != nil {
			errRecord = errors.WithStack(err)
			break
		}
		if err := todoService.UpdateTodo(todo); err != nil {
			errRecord = err
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if errRecord != nil {
		logError(errRecord)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func sanitizePageNo(source []string) int {
	if len(source) < 1 {
		return 1
	}

	var (
		rawPageNo string
		sanitized int
	)
	rawPageNo = source[0]
	sanitized, err := strconv.Atoi(rawPageNo)
	if err != nil || sanitized < 1 {
		sanitized = 1
	}
	return sanitized
}

func logError(err error) {
	log.Printf("%+v", err)
}