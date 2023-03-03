package controller

import (
	"net/http"
	"trelloBE/db"
	"trelloBE/schema"
	"trelloBE/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// /board/all
func GetBoards(w http.ResponseWriter, r *http.Request) {

}

// /board/create
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	userDetails := r.Context().Value("userDetails").(map[string]interface{})

	var userId primitive.ObjectID = userDetails["id"].(primitive.ObjectID)

	createdBoard, err := db.InsertBoard(userId)
	if err != nil {
		util.WriteApiMessage(w, "Failed to create board", 0, false)
		util.Log.Println("Failed to create board, err:", err)
		return
	}

	createdBoardResponse := schema.Board{
		Id:        createdBoard.Id,
		UserId:    createdBoard.UserId,
		MemberIds: createdBoard.MemberIds,
	}
	util.WriteDataToResponse(w, createdBoardResponse)
}
