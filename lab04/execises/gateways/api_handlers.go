package gateways

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"example.com/rest/domain"
	"github.com/emicklei/go-restful/v3"
)

const echoPath = "/echo"
const usersPatch = "/books"

type API struct {
	books map[int]domain.Book
}

func NewAPI() *API {
	return &API{
		books: make(map[int]domain.Book),
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/my-app")
	ws.Route(ws.POST(echoPath).To(api.echoPOSTHandler).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.GET(echoPath).To(api.echoGETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

	ws.Route(ws.GET(usersPatch).To(api.getBook).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.POST(usersPatch).To(api.addBook).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.PATCH(usersPatch).To(api.updateUser).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	// ws.Route(ws.DELETE(usersPatch).To(api.deleteUser).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

}

func (api *API) echoPOSTHandler(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	resp.WriteAsJson(map[string]string{"res": string(data)})
}

func (api *API) echoGETHandler(req *restful.Request, resp *restful.Response) {
	param := req.QueryParameter("echo-param")
	resp.WriteAsJson(map[string]string{
		"res": param,
	})
}

func (api *API) addBook(req *restful.Request, resp *restful.Response) {
	usr := &domain.Book{}
	err := req.ReadEntity(usr)
	if err != nil {
		log.Printf("[ERROR] Failed to read user, err=%v", err)
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	_, ok := api.books[usr.GetHash()]
	if ok {
		log.Printf("[ERROR] Book already exists")
		resp.WriteError(http.StatusConflict, fmt.Errorf("book already exists"))
		return
	}
	api.books[usr.GetHash()] = *usr
}

func (api *API) updateUser(req *restful.Request, resp *restful.Response) {

}

func (api *API) deleteUser(req *restful.Request, resp *restful.Response) {

}

func (api *API) getBook(req *restful.Request, resp *restful.Response) {
	title := req.QueryParameter("title")
	author := req.QueryParameter("author")
	publication_year, err := strconv.Atoi(req.QueryParameter("publication_year"))
	if err != nil {
		log.Printf("[ERROR] Failed to read publication_year")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("book publication_year must be provided"))
		return
	}
	number_of_downloads, err:= strconv.Atoi(req.QueryParameter("number_of_downloads"))
	
	if title == "" {
		log.Printf("[ERROR] Failed to read title")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("book title must be provided"))
		return
	}
	if author == "" {
		log.Printf("[ERROR] Failed to read author")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("book author must be provided"))
		return
	}
	
	if err != nil {
		log.Printf("[ERROR] Failed to read number_of_downloads")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("book number_of_downloads must be provided"))
		return
	}
	
	book := domain.Book{Title: title, Author: author, Publication_year: publication_year,
		Number_of_downloads: number_of_downloads}

	hash := book.GetHash()
	u, ok := api.books[hash]
	if !ok {
		log.Printf("[ERROR] Book not found")
		resp.WriteError(http.StatusNotFound, fmt.Errorf("book not found"))
		return
	}
	resp.WriteAsJson(u)
}