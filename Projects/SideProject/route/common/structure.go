package common

import (
	"context"
	"database/sql"
	"os"
)

type ArgsUpdateJoinUserInfo struct {
	CTX          context.Context
	Userimgfile  *os.File
	Joinuserinfo *ReqJoinInfo
	DB           *sql.DB
}

type ReqJoinInfo struct {
	UserInfo    JSUser  `json:"user"`
	AccountInfo Account `json:"account"`
}

type Account struct {
	UserID      int64  `json:"user_id"`
	Bank        int    `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

type JSUser struct {
	UserID              int64  `json:"user_id" binding:"required"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	Introduction        string `json:"introduction"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
}
