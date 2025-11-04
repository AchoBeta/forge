package router

import (
	"errors"
	"forge/interface/def"
	"forge/interface/handler"
	"forge/pkg/log/zlog"
	"forge/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func aiChatServiceErrorToMsgCode(err error) response.MsgCode {
	if err == nil {
		return response.SUCCESS
	}

	if errors.Is(err, errors.New("会话ID不能为空")) {
		return response.CONVERSATION_ID_NOT_NULL
	}
	if errors.Is(err, errors.New("用户ID不能为空")) {
		return response.USER_ID_NOT_NULL
	}
	if errors.Is(err, errors.New("导图ID不能为空")) {
		return response.MAP_ID_NOT_NULL
	}
	if errors.Is(err, errors.New("会话标题不能为空")) {
		return response.CONVERSATION_TITLE_NOT_NULL
	}
	if errors.Is(err, errors.New("该会话不存在")) {
		return response.CONVERSATION_NOT_EXIST
	}
	if errors.Is(err, errors.New("会话权限不足")) {
		return response.AI_CHAT_PERMISSION_DENIED
	}

	return response.COMMON_FAIL
}

// SendMessage 基础ai对话
func SendMessage() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var req def.ProcessUserMessageRequest
		ctx := gCtx.Request.Context()

		if err := gCtx.ShouldBindJSON(&req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_NOT_COMPLETE.Code,
				Message: response.PARAM_NOT_COMPLETE.Msg,
				Data:    def.ProcessUserMessageResponse{Success: false},
			})
		}

		resp, err := handler.GetHandler().SendMessage(ctx, &req)

		zlog.CtxAllInOne(ctx, "send_message", map[string]interface{}{"req": req}, resp, err)

		r := response.NewResponse(gCtx)
		if err != nil {
			msgCode := aiChatServiceErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.ProcessUserMessageResponse{Success: false},
			})
		} else {
			r.Success(resp)
		}
	}
}

// SaveNewConversation 保存新的会话
func SaveNewConversation() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var req def.SaveNewConversationRequest
		ctx := gCtx.Request.Context()

		if err := gCtx.ShouldBindJSON(&req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_NOT_COMPLETE.Code,
				Message: response.PARAM_NOT_COMPLETE.Msg,
				Data:    def.SaveNewConversationResponse{Success: false},
			})
		}

		resp, err := handler.GetHandler().SaveNewConversation(ctx, &req)

		zlog.CtxAllInOne(ctx, "save_new_conversation", map[string]interface{}{"req": req}, resp, err)

		r := response.NewResponse(gCtx)
		if err != nil {
			msgCode := aiChatServiceErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.SaveNewConversationResponse{Success: false},
			})
		} else {
			r.Success(resp)
		}
	}
}

// GetConversationList 获取某导图的所有会话
func GetConversationList() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var req def.GetConversationListRequest
		ctx := gCtx.Request.Context()
		if err := gCtx.ShouldBindQuery(&req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_NOT_COMPLETE.Code,
				Message: response.PARAM_NOT_COMPLETE.Msg,
				Data:    def.GetConversationListResponse{Success: false},
			})
		}

		resp, err := handler.GetHandler().GetConversationList(ctx, &req)

		zlog.CtxAllInOne(ctx, "get_conversation_list", map[string]interface{}{"req": req}, resp, err)

		r := response.NewResponse(gCtx)
		if err != nil {
			msgCode := aiChatServiceErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.GetConversationListResponse{Success: false},
			})
		} else {
			r.Success(resp)
		}
	}
}

func DelConversation() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var req def.DelConversationRequest
		ctx := gCtx.Request.Context()
		if err := gCtx.ShouldBindJSON(&req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_NOT_COMPLETE.Code,
				Message: response.PARAM_NOT_COMPLETE.Msg,
				Data:    def.DelConversationResponse{Success: false},
			})
		}

		resp, err := handler.GetHandler().DelConversation(ctx, &req)
		zlog.CtxAllInOne(ctx, "del_conversation", map[string]interface{}{"req": req}, resp, err)

		r := response.NewResponse(gCtx)
		if err != nil {
			msgCode := aiChatServiceErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.DelConversationResponse{Success: false},
			})
		} else {
			r.Success(resp)
		}
	}
}
