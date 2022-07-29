package gateway

import (
	"fmt"
	"lab05/domain"
	"net/http"
	"encoding/json"

	"github.com/emicklei/go-restful/v3"

	log "github.com/sirupsen/logrus"
)

const (
	booksPath      = "/books"
	booksPathStore = "/books/store/{id}"
)

type API struct {
	books   map[int]domain.Book
	storage domain.Storage
}

func NewAPI(storage domain.Storage) *API {
	return &API{
		books:   make(map[int]domain.Book),
		storage: storage}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/book-app")
	//ws.Route(ws.POST(booksPathStore).To(api.addBook).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Adding a new book in the database"))
	ws.Route(ws.GET(booksPathStore).To(api.getBook).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Getting a book from database"))
	ws.Route(ws.PUT(booksPathStore).To(api.addBook).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Adding a new book in the database"))
	//ws.Route(ws.PATCH(booksPathStore).To(api.saveBook).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Partial changing a new book in the database"))
	//ws.Route(ws.PUT(booksPathStore).To(api.saveBook).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Deleting a new book in the database"))
}

func (api *API) addBook(req *restful.Request, resp *restful.Response) {
	fmt.Println("add book")
	book := &domain.Book{}
	err := req.ReadEntity(book)
	id := req.PathParameter("id")
	if err != nil {
		log.WithError(err).Error("Failed to parse book json")
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		log.Error("error occured")
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("error occured"))
		return
	}
	api.storage.WriteContent(id, string(bookJSON))
	log.Info("Book added successfully in database")
}

func (api *API) getBook(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")

	b, err:= api.storage.GetContent(id)
	if err != nil {
		log.Error("Book not found")
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("book not found"))
		return
	}

	err = resp.WriteAsJson(b)
	if err != nil {
		log.WithError(err).Error("Failed to write response")
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}
