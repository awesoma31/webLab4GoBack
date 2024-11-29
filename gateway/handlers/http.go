package handlers

import (
	"awesoma31/common"
	pb "awesoma31/common/api"
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("POST /auth/reg", h.handleRegister)
	mux.HandleFunc("POST /auth/login", h.handleLogin)
	mux.HandleFunc("GET /api/v1/points/page", h.handleGetPage)
	mux.HandleFunc("POST /api/v1/points/add", h.handleAddPoint)

}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle Test")
	t, err := common.ExtractBearerToken(r)
	if err != nil {
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

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleGetPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle TestGetPage")
	// Extract query parameters
	queryParams := r.URL.Query()

	page := queryParams.Get("page")
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}

	pageSize := queryParams.Get("size")
	if pageSize == "" {
		pageSize = "10" // Default to size 10 if not provided
	}

	// Parse pageSize to int32
	size64, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid page size: "+err.Error())
		return
	}
	size32 := int32(size64)

	t, err := common.ExtractBearerToken(r)
	if err != nil {
		common.WriteError(w, http.StatusUnauthorized, "token required")
		return
	}

	authorization, err := h.authService.Authorize(r.Context(), &pb.AuthorizeRequest{Token: t})
	if err != nil {
		common.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	pointsPage, err := h.pointsService.GetUserPointsPage(r.Context(), &pb.PointsPageRequest{
		PageParam: page,
		PageSize:  size32,
		Id:        authorization.Id,
	})
	if err != nil {
		common.HandleAndWriteGrpcError(w, err)
		return
	}

	common.WriteJson(w, http.StatusOK, pointsPage)
}

func (h *Handler) handleAddPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle TestAddPoint")
	t, err := common.ExtractBearerToken(r)
	if err != nil {
		common.WriteError(w, http.StatusUnauthorized, "token required")
		return
	}

	authorization, err := h.authService.Authorize(r.Context(), &pb.AuthorizeRequest{Token: t})
	if err != nil {
		common.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var pointData pb.PointData

	err = common.ReadJSON(r, &pointData)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	pointResp, err := h.pointsService.AddPoint(r.Context(), &pb.AddPointRequest{
		PointsData:    &pointData,
		Authorization: authorization,
	})
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response := map[string]interface{}{
		"id":      pointResp.Id,
		"x":       pointResp.X,
		"y":       pointResp.Y,
		"r":       pointResp.R,
		"result":  pointResp.Result,
		"ownerId": authorization.Id,
	}

	common.WriteJson(w, http.StatusOK, response)
}
