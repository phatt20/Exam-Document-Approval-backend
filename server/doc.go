package server

import (
	repository "approval-system/internal/Repository"
	"approval-system/internal/handlers"
	"approval-system/internal/usecase"
)

func (s *server) docService() {

	repo := repository.NewDocRepository(s.postgres)

	docUsecase := usecase.NewDocUsecase(repo)

	docHandler := handlers.NewDocHttpHandler(s.cfg, docUsecase)

	docGroup := s.app.Group("/doc_v1")

	docGroup.POST("/doc/create", docHandler.CreateDoc)
	docGroup.GET("/doc", docHandler.FindAllDocs)
	docGroup.GET("/doc/:id", docHandler.FindDocByID)
	docGroup.PUT("/doc/update-staus", docHandler.UpdateStaus)
}
