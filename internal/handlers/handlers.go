package handlers

import (
	"approval-system/config"
	"approval-system/internal/domain"
	docUsecase "approval-system/internal/usecase"
	"approval-system/pkg/request"
	"approval-system/pkg/response"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DocHttpHandler interface {
	CreateDoc(c echo.Context) error
	FindAllDocs(c echo.Context) error
	FindDocByID(c echo.Context) error
	UpdateStaus(c echo.Context) error
}

type docHttpHandler struct {
	cfg        *config.Config
	docUsecase docUsecase.DocUsecaseInterface
}

func NewDocHttpHandler(cfg *config.Config, docUsecase docUsecase.DocUsecaseInterface) DocHttpHandler {
	return &docHttpHandler{cfg, docUsecase}
}

func (h *docHttpHandler) CreateDoc(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(domain.CreateDocumentInput)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.docUsecase.CreateDoc(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *docHttpHandler) FindAllDocs(c echo.Context) error {
	ctx := context.Background()

	res, err := h.docUsecase.FindAllDocs(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *docHttpHandler) FindDocByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid document ID")
	}

	res, err := h.docUsecase.FindDocByID(ctx, id)
	if err != nil {
		return response.ErrResponse(c, http.StatusNotFound, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *docHttpHandler) UpdateStaus(c echo.Context) error {
	ctx := c.Request().Context()
	wrapper := request.ContextWrapper(c)

	req := new(domain.UpdateStatusInput)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.docUsecase.UpdateStatus(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
