package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chatwithsimchar "github.com/tanmaij/zylo/internal/controller/chat_with_sim_char"
	"github.com/tanmaij/zylo/internal/model"
)

type Impl struct {
	ctrl chatwithsimchar.Impl
}

func New(ctrl chatwithsimchar.Impl) Impl {
	return Impl{ctrl: ctrl}
}

type roomRes struct {
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Desc    string `json:"desc"`
	Address string `json:"address"`
}

type messageRes struct {
	Sender    string `json:"sender"`
	SenderAvt string `json:"senderAvatar"`
	Message   string `json:"message"`
}

type getCurrentConversationResponse struct {
	Room     roomRes      `json:"room"`
	Messages []messageRes `json:"messages"`
}

func conversationToResponse(clientUUID string, conv model.Conversation) getCurrentConversationResponse {
	var msgsRs = make([]messageRes, len(conv.Messages))
	for i, msg := range conv.Messages {
		var sender, senderAvt string
		if msg.Role == model.RoleAssistant {
			sender = conv.Character.Name
			senderAvt = conv.Character.AvatarURL
		} else {
			sender = clientUUID
		}

		msgsRs[i] = messageRes{
			Sender:    sender,
			SenderAvt: senderAvt,
			Message:   conv.Messages[i].Content,
		}
	}

	return getCurrentConversationResponse{
		Room: roomRes{
			Name:    conv.Character.Name,
			Avatar:  conv.Character.AvatarURL,
			Address: conv.Character.Address,
			Desc:    conv.Character.Description,
		},
		Messages: msgsRs,
	}
}

func (impl Impl) GetCurrentConversation(w http.ResponseWriter, r *http.Request) error {
	clientUUID := chi.URLParam(r, "client_uuid")

	conv, err := impl.ctrl.GetCurrentConversation(r.Context(), chatwithsimchar.GetCurrentConversation{
		ClientUUID: clientUUID,
	})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(conversationToResponse(clientUUID, conv))
}

type chatRequest struct {
	Message string `json:"message"`
}

func (impl Impl) Chat(w http.ResponseWriter, r *http.Request) error {
	clientUUID := chi.URLParam(r, "client_uuid")

	var rq chatRequest
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		return HandlerError{
			Code:        400,
			Description: "Invalid Body",
		}
	}

	if rq.Message == "" {
		return HandlerError{
			Code:        400,
			Description: "Invalid msg",
		}
	}

	chat, err := impl.ctrl.Chat(r.Context(), chatwithsimchar.ChatInput{
		ClientUUID: clientUUID,
		Message:    rq.Message,
	})
	if err != nil {
		switch err {
		case chatwithsimchar.ErrConvNotFound:
			return HandlerError{
				Code:        404,
				Description: "Conversation not found",
			}
		default:
			return err
		}
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(messageRes{
		Sender:  clientUUID,
		Message: chat.Content,
	})
}
