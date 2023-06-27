package messagecontroller

import (
	"chatroom/common/converter"
	commonmodel "chatroom/common/model"
	"chatroom/logger"
	messagemodel "chatroom/message/model"
	messageservice "chatroom/message/service"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type messageController struct {
	validate       *validator.Validate
	messageService messageservice.MessageService
}

type Options struct {
	Validate       *validator.Validate
	MessageService messageservice.MessageService
}

func New(o *Options) *messageController {
	return &messageController{
		validate:       o.Validate,
		messageService: o.MessageService,
	}
}

func (mc *messageController) Messages(r *http.Request, userID int64) (
	*commonmodel.PaginatedRes[[]messagemodel.MessageRes],
	int,
	error,
) {
	page, err := converter.StringToInt(converter.GetParam(r, "page"))
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid params")
		return nil, http.StatusBadRequest, err
	}

	roomID, err := converter.StringToInt64(converter.GetParam(r, "roomId"))
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid params")
		return nil, http.StatusBadRequest, err
	}

	err = mc.validate.Var(page, "gt=0")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// TODO: passing user id
	messages, err := mc.messageService.Messages(userID, page, roomID, 10)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	messagesModel := messagemodel.MapMessagesEntityToModel(*messages)

	return &commonmodel.PaginatedRes[[]messagemodel.MessageRes]{
		Page:  page,
		Limit: 10,
		Data:  messagesModel,
	}, http.StatusOK, nil
}
