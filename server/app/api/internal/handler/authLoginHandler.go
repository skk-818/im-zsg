package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/app/api/internal/logic"
	"server/app/api/internal/svc"
	"server/app/api/internal/types"
)

func authLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuthLoginLogic(r.Context(), svcCtx)
		resp, err := l.AuthLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
