package api

//create contest要import"hiper-backend/mq"，在发送200前要先mq.SendBuildGameMsg(model.Ctx, tempGame.ID)，参见creategame最后几句
