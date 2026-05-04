package handlers

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	api "lab-tracker/internal/delivery/http/gen"
	"lab-tracker/internal/domain"
)

type Service interface {
	ListGroups(ctx context.Context) ([]domain.Group, error)
	ListLabWorks(ctx context.Context, filter domain.LabWorkFilter) (*domain.LabWorkList, error)
	GetLabWork(ctx context.Context, id int64) (*domain.LabWork, error)
	CreateLabWork(ctx context.Context, input domain.LabWorkInput) (*domain.LabWork, error)
	UpdateLabWork(ctx context.Context, id int64, input domain.LabWorkInput) (*domain.LabWork, error)
	DeleteLabWork(ctx context.Context, id int64) error
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password, fullName string, groupID int) (string, error)
}

type Handler struct {
	svc Service
}

func New(s Service) *Handler {
	return &Handler{svc: s}
}

// (POST /api/v1/auth/login)
func (h *Handler) PostApiV1AuthLogin(
	ctx context.Context,
	request api.PostApiV1AuthLoginRequestObject,
) (api.PostApiV1AuthLoginResponseObject, error) {
	token, err := h.svc.Login(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			return api.PostApiV1AuthLogin400JSONResponse(invalidRequest("invalid credentials")), nil

		case domain.ErrUserNotFound:
			return api.PostApiV1AuthLogin401JSONResponse(userNotFound()), nil

		case domain.ErrUnauthorized:
			return api.PostApiV1AuthLogin401JSONResponse(unauthorized()), nil
		}
		return nil, err
	}

	return api.PostApiV1AuthLogin200JSONResponse{Token: token}, nil
}

// (POST /api/v1/auth/register)
func (h *Handler) PostApiV1AuthRegister(
	ctx context.Context,
	request api.PostApiV1AuthRegisterRequestObject,
) (api.PostApiV1AuthRegisterResponseObject, error) {
	token, err := h.svc.Register(
		ctx,
		request.Body.Username,
		request.Body.Password,
		request.Body.FullName,
		request.Body.GroupId,
	)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			return api.PostApiV1AuthRegister409JSONResponse(invalidRequest("user already exists")), nil

		case domain.ErrGroupNotFound:
			return api.PostApiV1AuthRegister404JSONResponse(groupNotFound()), nil
		}
		return nil, err
	}

	return api.PostApiV1AuthRegister201JSONResponse{Token: token}, nil
}

// (GET /api/v1/groups/list)
func (h *Handler) GetApiV1GroupsList(
	ctx context.Context,
	_ api.GetApiV1GroupsListRequestObject,
) (api.GetApiV1GroupsListResponseObject, error) {
	groups, err := h.svc.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]api.Group, 0, len(groups))
	for _, g := range groups {
		result = append(result, api.Group{
			Id:     int(g.ID),
			Name:   g.Name,
			Course: g.Course,
		})
	}

	return api.GetApiV1GroupsList200JSONResponse{Groups: &result}, nil
}

// (POST /api/v1/groups/create)
func (h *Handler) PostApiV1GroupsCreate(
	ctx context.Context,
	request api.PostApiV1GroupsCreateRequestObject,
) (api.PostApiV1GroupsCreateResponseObject, error) {
	// TODO: реализовать
	return api.PostApiV1GroupsCreate201Response{}, nil
}

// (GET /api/v1/labworks)
func (h *Handler) GetApiV1Labworks(
	ctx context.Context,
	request api.GetApiV1LabworksRequestObject,
) (api.GetApiV1LabworksResponseObject, error) {
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}

	limit := 10
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	search := ""
	if request.Params.Search != nil {
		search = *request.Params.Search
	}

	if page < 1 || limit < 1 || limit > 20 {
		return api.GetApiV1Labworks400JSONResponse(invalidRequest("")), nil
	}

	labWorks, err := h.svc.ListLabWorks(ctx, domain.LabWorkFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRequest) {
			return api.GetApiV1Labworks400JSONResponse(invalidRequest("")), nil
		}
		return api.GetApiV1Labworks500JSONResponse(internalError()), nil
	}

	items := make([]api.LabWork, 0, len(labWorks.Items))
	for _, item := range labWorks.Items {
		items = append(items, toAPILabWork(item))
	}

	return api.GetApiV1Labworks200JSONResponse{
		Items:      items,
		Total:      labWorks.Total,
		Page:       labWorks.Page,
		Limit:      labWorks.Limit,
		TotalPages: labWorks.TotalPages,
	}, nil
}

// (POST /api/v1/labworks)
func (h *Handler) PostApiV1Labworks(
	ctx context.Context,
	request api.PostApiV1LabworksRequestObject,
) (api.PostApiV1LabworksResponseObject, error) {
	userID, role, ok := authorizedUser(ctx)
	if !ok {
		return api.PostApiV1Labworks401JSONResponse(unauthorized()), nil
	}
	if role != domain.RoleAdmin {
		return api.PostApiV1Labworks403JSONResponse(forbidden()), nil
	}
	if request.Body == nil {
		return api.PostApiV1Labworks400JSONResponse(invalidRequest("invalid request body")), nil
	}

	labWork, err := h.svc.CreateLabWork(ctx, fromAPILabWorkInput(*request.Body, &userID))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRequest):
			return api.PostApiV1Labworks400JSONResponse(invalidRequest("")), nil
		default:
			return api.PostApiV1Labworks500JSONResponse(internalError()), nil
		}
	}

	response := api.PostApiV1Labworks201JSONResponse(toAPILabWork(*labWork))
	return response, nil
}

// (GET /api/v1/labworks/{id})
func (h *Handler) GetApiV1LabworksId(
	ctx context.Context,
	request api.GetApiV1LabworksIdRequestObject,
) (api.GetApiV1LabworksIdResponseObject, error) {
	if request.Id <= 0 {
		return api.GetApiV1LabworksId400JSONResponse(invalidRequest("invalid id")), nil
	}

	labWork, err := h.svc.GetLabWork(ctx, int64(request.Id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLabWorkNotFound):
			return api.GetApiV1LabworksId404JSONResponse(labWorkNotFound()), nil
		case errors.Is(err, domain.ErrInvalidRequest):
			return api.GetApiV1LabworksId400JSONResponse(invalidRequest("")), nil
		default:
			return api.GetApiV1LabworksId500JSONResponse(internalError()), nil
		}
	}

	response := api.GetApiV1LabworksId200JSONResponse(toAPILabWork(*labWork))
	return response, nil
}

// (PUT /api/v1/labworks/{id})
func (h *Handler) PutApiV1LabworksId(
	ctx context.Context,
	request api.PutApiV1LabworksIdRequestObject,
) (api.PutApiV1LabworksIdResponseObject, error) {
	_, role, ok := authorizedUser(ctx)
	if !ok {
		return api.PutApiV1LabworksId401JSONResponse(unauthorized()), nil
	}
	if role != domain.RoleAdmin {
		return api.PutApiV1LabworksId403JSONResponse(forbidden()), nil
	}
	if request.Id <= 0 {
		return api.PutApiV1LabworksId400JSONResponse(invalidRequest("invalid id")), nil
	}
	if request.Body == nil {
		return api.PutApiV1LabworksId400JSONResponse(invalidRequest("invalid request body")), nil
	}

	labWork, err := h.svc.UpdateLabWork(ctx, int64(request.Id), fromAPILabWorkInput(*request.Body, nil))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRequest):
			return api.PutApiV1LabworksId400JSONResponse(invalidRequest("")), nil
		case errors.Is(err, domain.ErrLabWorkNotFound):
			return api.PutApiV1LabworksId404JSONResponse(labWorkNotFound()), nil
		default:
			return api.PutApiV1LabworksId500JSONResponse(internalError()), nil
		}
	}

	response := api.PutApiV1LabworksId200JSONResponse(toAPILabWork(*labWork))
	return response, nil
}

// (DELETE /api/v1/labworks/{id})
func (h *Handler) DeleteApiV1LabworksId(
	ctx context.Context,
	request api.DeleteApiV1LabworksIdRequestObject,
) (api.DeleteApiV1LabworksIdResponseObject, error) {
	_, role, ok := authorizedUser(ctx)
	if !ok {
		return api.DeleteApiV1LabworksId401JSONResponse(unauthorized()), nil
	}
	if role != domain.RoleAdmin {
		return api.DeleteApiV1LabworksId403JSONResponse(forbidden()), nil
	}
	if request.Id <= 0 {
		return api.DeleteApiV1LabworksId400JSONResponse(invalidRequest("invalid id")), nil
	}

	if err := h.svc.DeleteLabWork(ctx, int64(request.Id)); err != nil {
		switch {
		case errors.Is(err, domain.ErrLabWorkNotFound):
			return api.DeleteApiV1LabworksId404JSONResponse(labWorkNotFound()), nil
		case errors.Is(err, domain.ErrInvalidRequest):
			return api.DeleteApiV1LabworksId400JSONResponse(invalidRequest("")), nil
		default:
			return api.DeleteApiV1LabworksId500JSONResponse(internalError()), nil
		}
	}

	return api.DeleteApiV1LabworksId204Response{}, nil
}

// (POST /api/v1/assignments/create)
func (h *Handler) PostApiV1AssignmentsCreate(
	ctx context.Context,
	request api.PostApiV1AssignmentsCreateRequestObject,
) (api.PostApiV1AssignmentsCreateResponseObject, error) {
	// TODO: реализовать
	return api.PostApiV1AssignmentsCreate201Response{}, nil
}

// (POST /api/v1/submissions/create)
func (h *Handler) PostApiV1SubmissionsCreate(
	ctx context.Context,
	request api.PostApiV1SubmissionsCreateRequestObject,
) (api.PostApiV1SubmissionsCreateResponseObject, error) {
	// TODO: реализовать
	return api.PostApiV1SubmissionsCreate201Response{}, nil
}

// (POST /api/v1/grades/set)
func (h *Handler) PostApiV1GradesSet(
	ctx context.Context,
	request api.PostApiV1GradesSetRequestObject,
) (api.PostApiV1GradesSetResponseObject, error) {
	// TODO: реализовать
	return api.PostApiV1GradesSet200Response{}, nil
}

func authorizedUser(ctx context.Context) (uuid.UUID, domain.UserRole, bool) {
	roleValue, ok := ctx.Value("role").(string)
	if !ok || roleValue == "" {
		return uuid.UUID{}, "", false
	}

	userIDValue, ok := ctx.Value("user_id").(string)
	if !ok || userIDValue == "" {
		return uuid.UUID{}, "", false
	}

	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		return uuid.UUID{}, "", false
	}

	return userID, domain.UserRole(roleValue), true
}

func fromAPILabWorkInput(input api.LabWorkInput, createdBy *uuid.UUID) domain.LabWorkInput {
	return domain.LabWorkInput{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Goal:        strings.TrimSpace(input.Goal),
		Equipment:   strings.TrimSpace(input.Equipment),
		Reagents:    strings.TrimSpace(input.Reagents),
		Procedure:   strings.TrimSpace(input.Procedure),
		FilePath:    input.FilePath,
		CreatedBy:   createdBy,
	}
}

func toAPILabWork(labWork domain.LabWork) api.LabWork {
	createdAt := time.Time{}
	if parsed, err := time.Parse(time.RFC3339, labWork.CreatedAt); err == nil {
		createdAt = parsed
	}

	return api.LabWork{
		Id:          int(labWork.ID),
		Title:       labWork.Title,
		Description: labWork.Description,
		Goal:        labWork.Goal,
		Equipment:   labWork.Equipment,
		Reagents:    labWork.Reagents,
		Procedure:   labWork.Procedure,
		FilePath:    labWork.FilePath,
		CreatedAt:   createdAt,
	}
}
