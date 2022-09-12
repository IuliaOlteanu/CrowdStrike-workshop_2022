package api

import (

	"exam-api/domain"
	"sync"
	"fmt"
	"net/http"
	
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
)
var wg sync.WaitGroup
var mu sync.Mutex

func (api *API) createProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	mu.Lock()
	product := &domain.Product{}
	err := req.ReadEntity(product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	id, alreadyExists, err := api.storage.Save(*product)
	if err != nil {
		log.Errorf("Failed to save product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to save product"))
		return
	}

	if alreadyExists {
		log.Infof("Product %s already in store", id)
		_ = resp.WriteError(http.StatusConflict, fmt.Errorf("product already exists"))
		return
	}

	log.Infof("Product %s saved in store", id)

	_ = resp.WriteAsJson(map[string]string{
		"id": id,
	})
	mu.Unlock()
	wg.Done()
}

func (api *API) getProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}

func (api *API) updateProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}

func (api *API) deleteProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}
