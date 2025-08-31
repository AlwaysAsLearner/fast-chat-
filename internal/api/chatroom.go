package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/api/middleware"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/services"
)

type ChatroomHandler struct {
	Service *services.ChatroomService
}

func NewChatroomHandler(cs *services.ChatroomService) *ChatroomHandler {
	return &ChatroomHandler{Service: cs}
}

func RegisterChatroomRoutes(mux *http.ServeMux, h *ChatroomHandler) {
	mux.HandleFunc("/chatrooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// authenticated route
			middleware.AuthMiddleware(http.HandlerFunc(h.CreateChatroom)).ServeHTTP(w, r)
		case http.MethodGet:
			// public route
			h.ListPublicChatrooms(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	// pattern starting with /chatrooms/  all
	mux.HandleFunc("/chatrooms/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path) // ["chatrooms", "123", "join"]
		if len(parts) < 3 {
			notFound(w)
			return
		}

		id, err := parseUint(parts[1])
		if err != nil {
			badRequest(w, "invalid chatroom id")
			return
		}

		switch parts[2] {
		case "join":
			if r.Method != http.MethodPost {
				methodNotAllowed(w)
				return
			}
			middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.JoinChatroom(w, r, id)
			})).ServeHTTP(w, r)
		case "leave":
			if r.Method != http.MethodDelete {
				methodNotAllowed(w)
				return
			}

			middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.LeaveChatroom(w, r, id)
			})).ServeHTTP(w, r)
		default:
			notFound(w)
		}
	})

	mux.Handle("/users/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path) // ["users", "123", "chatrooms"]
		if len(parts) == 3 && parts[2] == "chatrooms" && r.Method == http.MethodGet {
			userId, err := parseUint(parts[1])
			if err != nil {
				badRequest(w, "invalid user id")
				return
			}

			middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ListUserChatrooms(w, r, userId)
			})).ServeHTTP(w, r)
			return
		}
		notFound(w)
	}))

}

// DTOs - Data Transfer Objects
type createChatroomReq struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsPrivate   bool    `json:"is_private"`
}

type ChatroomDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsPrivate   bool    `json:"is_private"`
	OwnerID     uint    `json:"owner_id"`
	MemberCount int     `json:"member_count"`
	CreatedAt   string  `json:"created_at"`
}

// POST /chatrooms
func (ch *ChatroomHandler) CreateChatroom(w http.ResponseWriter, r *http.Request) {
	setJSON(w) // set the response header with json
	var req createChatroomReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		badRequest(w, "name is required")
		return
	}

	// extract userId from the jwt token from the context
	userID, ok := ctxUserID(r)
	if !ok {
		unauthorized(w)
		return
	}

	chatroom, err := ch.Service.CreateChatroom(req.Name, req.Description, req.IsPrivate, userID)

	if err != nil {
		serverError(w, err)
	}

	writeJSON(w, http.StatusCreated, toChatroomDTO(chatroom))
}

// GET /chatrooms?page=&limit=&q=
func (h *ChatroomHandler) ListPublicChatrooms(w http.ResponseWriter, r *http.Request) {
	setJSON(w)

	page, limit := parsePageLimit(r, 1, 20, 200)
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	rooms, total, err := h.Service.ListAllPublicChatrooms(page, limit, q)
	if err != nil {
		serverError(w, err)
		return
	}
	resp := map[string]interface{}{
		"page":      page,
		"limit":     limit,
		"total":     total,
		"chatrooms": toChatroomDTOs(rooms),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ChatroomHandler) JoinChatroom(w http.ResponseWriter, r *http.Request, chatID uint) {
	setJSON(w)

	userID, ok := ctxUserID(r)
	if !ok {
		unauthorized(w)
		return
	}

	if err := h.Service.JoinPublicChatroom(chatID, userID); err != nil {
		serverError(w, err)
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "joined chatroom successfully"})
}

func (h *ChatroomHandler) LeaveChatroom(w http.ResponseWriter, r *http.Request, chatroomID uint) {
	setJSON(w)
	userID, ok := ctxUserID(r)
	if !ok {
		unauthorized(w)
		return
	}

	if err := h.Service.LeaveChatrooom(chatroomID, userID); err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "left"})
}

func (h *ChatroomHandler) ListUserChatrooms(w http.ResponseWriter, r *http.Request, userID uint) {
	setJSON(w)

	if me, ok := ctxUserID(r); !ok || me != userID {
		forbidden(w, "you can only view your own chatrooms")
		return
	}

	rooms, err := h.Service.GetUserJoinedChatrooms(userID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toChatroomDTOs(rooms))
}

/*  **** HELPER FUNCTIONS **** */
func splitPath(path string) []string {
	p := strings.TrimPrefix(path, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func parseUint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	return uint(u64), err
}

func parsePageLimit(r *http.Request, defLimit, defPage, maxLimit int) (int, int) {
	/**
	 1. Parse page and limit, if it does not exist use default values
	 2. When parsing ensure its integer and greater than 0
	 3. Check if the limit exceeds maxLimit, and set that to maxLimit
	 4. return page and limit
	**/
	q := r.URL.Query()
	page := defPage
	limit := defLimit
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= maxLimit {
			limit = n
		}
	}
	return page, limit
}

func ctxUserID(r *http.Request) (uint, bool) {
	// get userid and parse it into uint and return id and ok status
	val := r.Context().Value(middleware.UserIDKey)
	if val == nil {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}

func setJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	// write header status and encode the response in json format
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func badRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func unauthorized(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
}

func forbidden(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusForbidden, map[string]string{"error": msg})
}

func notFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func conflict(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusConflict, map[string]string{"error": msg})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}
func serverError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("server error: %v", err)})
}

func toChatroomDTO(c services.Chatroom) ChatroomDTO {
	return ChatroomDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		IsPrivate:   c.IsPrivate,
		OwnerID:     c.OwnerID,
		MemberCount: c.MemberCount,
		CreatedAt:   c.CreatedAt.String(),
	}
}

func toChatroomDTOs(cs []services.Chatroom) []ChatroomDTO {
	out := make([]ChatroomDTO, 0, len(cs))
	for _, c := range cs {
		out = append(out, toChatroomDTO(c))
	}
	return out
}
