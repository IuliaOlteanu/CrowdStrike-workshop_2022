package api

import (
	"exam-api/domain"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
)

func (api *API) createProductMemorySingle(req *restful.Request, resp *restful.Response) {
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

}

func (api *API) getProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	product, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	_ = resp.WriteAsJson(product)
}

func (api *API) updateProductMemorySingle(req *restful.Request, resp *restful.Response) {
	productdiff := &domain.ProductDiff{}
	err := req.ReadEntity(productdiff)
	if err != nil {
		log.Errorf("Failed to read productdiff, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	product, exists, err := api.storage.Get(productdiff.ID)
	id := productdiff.ID
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	if productdiff.Diff.Price != 0 {
		product.Price = productdiff.Diff.Price
	}
	if productdiff.Diff.Stock != 0 {
		product.Stock = productdiff.Diff.Stock
	}
	if productdiff.Diff.Tags != nil {
		product.Tags = productdiff.Diff.Tags
	}

	_ = resp.WriteAsJson(product)

}

func (api *API) deleteProductMemorySingle(req *restful.Request, resp *restful.Response) {

	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	_, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	// product = domain.Product{}

	ok, err := api.storage.Delete(id)
	if !ok {
		log.Infof("Product hasn't been deleted", id)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("product not deleted"))
		return
	}
	

}
