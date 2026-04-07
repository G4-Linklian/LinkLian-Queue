package qadeliver

import (
	"worker/models"
	"worker/utils"
	"worker/utils/logger"
)

func HandleLiveRoom(event models.SocketEvent) error {
	logger.Log("Worker: Emitting "+event.Type+" to socket server", "QADeliver - LiveRoom")

	qaLiveId := event.Payload["qa_live_id"]
	sectionId := event.Payload["section_id"]
	liveBy := event.Payload["live_by"]
	liveTitle := event.Payload["live_title"]
	status, _ := event.Payload["status"].(string)
	startedAt := event.Payload["started_at"]
	endedAt := event.Payload["ended_at"]

	if event.Type == "QA_LIVE_STARTED" {
		logger.Log("Preparing Redis for new room...", "QADeliver - LiveRoom")
		// TODO: สั่ง Redis เคลียร์ค่า หรือตั้งค่าเริ่มต้นให้ห้องนี้
		
	} else if event.Type == "QA_LIVE_ENDED" {
		logger.Log("Cleaning up Redis after room ended...", "QADeliver - LiveRoom")
		// TODO: สั่งลบข้อมูลห้องนี้ทิ้งจาก Redis
	}

	err := utils.EmitToSocket(event.Type, map[string]interface{}{
		"qa_live_id": qaLiveId,
		"section_id": sectionId,
		"live_by":    liveBy,
		"live_title": liveTitle,
		"status":     status,
		"started_at": startedAt,
		"ended_at":   endedAt,
	}, "/ws/qa")
	
	if err != nil {
		logger.Error("Emit socket error", "QADeliver - LiveRoom", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: "+event.Type+" Processed Successfully", "QALive")
	return nil
}

func HandleFile(event models.SocketEvent) error {
	logger.Log("Worker: Emitting FILE_CHANGED to socket server", "QADeliver - FileChanged")

	qaLiveId := event.Payload["qa_live_id"]
	postId := event.Payload["post_id"]
	attachmentId := event.Payload["attachment_id"]
	openedAt := event.Payload["opened_at"]

	err := utils.EmitToSocket("FILE_CHANGED", map[string]interface{}{
		"qa_live_id":    qaLiveId,
		"post_id":       postId,
		"attachment_id": attachmentId,
		"opened_at":     openedAt,
	}, "/ws/qa")
	
	if err != nil {
		logger.Error("Emit socket error", "QADeliver - FileChanged", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: FILE_CHANGED Processed Successfully", "QADeliver - FileChanged")
	return nil
}

func HandleQuestion(event models.SocketEvent) error {
	logger.Log("Worker: Emitting "+event.Type+" to socket server", "QADeliver - Question")

	qaLiveId := event.Payload["qa_live_id"]
	qaQuestionId := event.Payload["qa_question_id"]
	question, _ := event.Payload["question"].(string)
	asker := event.Payload["asker"]
	postId := event.Payload["post_id"]
	attachmentId := event.Payload["attachment_id"]
	slideNumber := event.Payload["slide_number"]
	status, _ := event.Payload["status"].(string)
	upvoteCount := event.Payload["upvote_count"]
	createdAt := event.Payload["created_at"]

	err := utils.EmitToSocket(event.Type, map[string]interface{}{
		"qa_live_id":     qaLiveId,
		"qa_question_id": qaQuestionId,
		"question":       question,
		"asker":          asker,
		"post_id":        postId,
		"attachment_id":  attachmentId,
		"slide_number":   slideNumber,
		"status":         status,
		"upvote_count":   upvoteCount,
		"created_at":     createdAt,
	}, "/ws/qa")
	
	if err != nil {
		logger.Error("Emit socket error", "QADeliver - Question", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: "+event.Type+" Processed Successfully", "QADeliver - Question")
	return nil
}

func HandleVote(event models.SocketEvent) error {
	logger.Log("Worker: Emitting QA_UPVOTED to socket server", "QADeliver - Vote")

	qaLiveId := event.Payload["qa_live_id"]
	qaQuestionId := event.Payload["qa_question_id"]
	upvoteCount := event.Payload["upvote_count"]

	err := utils.EmitToSocket("QA_UPVOTED", map[string]interface{}{
		"qa_live_id":     qaLiveId,
		"qa_question_id": qaQuestionId,
		"upvote_count":   upvoteCount,
	}, "/ws/qa")
	
	if err != nil {
		logger.Error("Emit socket error", "QADeliver - Vote", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: QA_UPVOTED Processed Successfully", "QADeliver - Vote")
	return nil
}