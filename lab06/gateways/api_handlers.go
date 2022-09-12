package gateways

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"lab06/domain"
	"github.com/emicklei/go-restful/v3"
)
const echoPath = "/echo"
const usersPatch = "/items"

type API struct {
	items map[int]domain.Item
}

func NewAPI() *API {
	return &API{
		items: make(map[int]domain.Item),
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/my-app")
	ws.Route(ws.POST(echoPath).To(api.echoPOSTHandler).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.GET(echoPath).To(api.echoGETHandler).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))

	ws.Route(ws.GET(usersPatch).To(api.getItem).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
	ws.Route(ws.POST(usersPatch).To(api.addItem).Writes(restful.MIME_JSON).Doc("Writes back a json with what you gave it"))
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

func (api *API) addItem(req *restful.Request, resp *restful.Response) {
	usr := &domain.Item{}
	err := req.ReadEntity(usr)
	if err != nil {
		log.Printf("[ERROR] Failed to read item, err=%v", err)
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	_, ok := api.items[usr.GetHash()]
	if ok {
		log.Printf("[ERROR] Item already exists")
		resp.WriteError(http.StatusConflict, fmt.Errorf("Item already exists"))
		return
	}

	x1 := rand.NewSource(time.Now().UnixNano())
	y1 := rand.New(x1)
	val := y1.Intn(200)
	if  _, ok := api.items[val]; !ok {
		usr.Id = val
		api.items[val] = *usr
	}
}

func (api *API) getItem(req *restful.Request, resp *restful.Response) {
	name := req.QueryParameter("name")
	origin_country := req.QueryParameter("origin_country")
	price, err := strconv.Atoi(req.QueryParameter("price"))
	if err != nil {
		log.Printf("[ERROR] Failed to read price")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("item price must be provided"))
		return
	}
	
	if name == "" {
		log.Printf("[ERROR] Failed to read name")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("item name must be provided"))
		return
	}
	if origin_country == "" {
		log.Printf("[ERROR] Failed to read origin_country")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("item origin_country must be provided"))
		return
	}

	quantity, err:= strconv.Atoi(req.QueryParameter("quantity"))
	if err != nil {
		log.Printf("[ERROR] Failed to read quantity")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("item quantity must be provided"))
		return
	}
	
	item := domain.Item{Name: name, Price: price, Quantity: quantity,
		Origin_country: origin_country}

	hash := item.Id
	u, ok := api.items[hash]
	if !ok {
		log.Printf("[ERROR] Item not found")
		resp.WriteError(http.StatusNotFound, fmt.Errorf("Item not found"))
		return
	}
	resp.WriteAsJson(u)
}