package handlers

import (
	"awesoma31/common"
	pb "awesoma31/common/api"
	"log"
	"net/http"
)

type Handler struct {
	authService   pb.AuthServiceClient
	pointsService pb.PointsServiceClient
}

func NewHandler(asc pb.AuthServiceClient, psc pb.PointsServiceClient) *Handler {
	return &Handler{asc, psc}
}

func (h *Handler) MountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /test/auth", h.handleTest)
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle Test")
	t, err := common.ExtractBearerToken(r)
	if err != nil {
		//todo
		common.WriteError(w, http.StatusUnauthorized, "token required")
		return
	}

	authorization, err := h.authService.Authorize(r.Context(), &pb.AuthorizeRequest{
		Token: t,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("authorization: %s\n", authorization)

	pointsPage, err := h.pointsService.GetUserPointsPage(r.Context(), &pb.PointsPageRequest{
		PageParam: "1",
		PageSize:  1,
		Id:        1,
	})
	if err != nil {
		log.Fatal(err)
	}

	common.WriteJson(w, http.StatusOK, pointsPage)
}
