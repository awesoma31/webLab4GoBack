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
	mux.HandleFunc("POST /test/reg", h.handleRegister)
	mux.HandleFunc("POST /test/login", h.handleTestLogin)
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

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("TestReg")

	var regReq pb.LoginRequest
	err := common.ReadJSON(r, &regReq)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	_, err = h.authService.Register(r.Context(), &regReq)
	if err != nil {
		common.HandleAndWriteGrpcError(w, err)
		return
	}

	log.Printf("User registered successfully: %s", regReq.Username)

	common.WriteJson(w, http.StatusOK, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *Handler) handleTestLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle TestLogin")

	var loginReq pb.LoginRequest
	err := common.ReadJSON(r, &loginReq)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	loginResp, err := h.authService.Login(r.Context(), &loginReq)
	if err != nil {
		common.HandleAndWriteGrpcError(w, err)
		return
	}

	log.Printf("User logged in successfully: %s", loginReq.Username)
	common.WriteJson(w, http.StatusOK, map[string]interface{}{
		"accessToken":  loginResp.AccessToken,
		"refreshToken": loginResp.RefreshToken,
	})
}
