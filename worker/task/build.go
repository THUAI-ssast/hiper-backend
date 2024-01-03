package task

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Build(values map[string]interface{}) (status int) {
	idInt, err := strconv.Atoi(values["id"].(string))
	if err != nil {
		log.Fatal(err)
	}
	id := uint(idInt)
	status = 1

	switch values["type"] {
	case "game_logic":
		status = buildGameLogic(id)
	case "ai":
		status = buildAI(id)
	}
	return
}

// 获取任务所需信息
// 起容器，执行任务
// 等待任务完成，获取任务输出，保存与修改相关信息

func buildGameLogic(gameID uint) (status int) {
	// 获取游戏
	game, err := model.GetGameByID(gameID)
	if err != nil {
		log.Fatal(err)
	}
	// 获取游戏逻辑的构建任务
	gameLogic := game.GameLogic
	gameLogicBuild := gameLogic.Build
	// 获取游戏逻辑的构建任务的 Dockerfile
	dockerfile := gameLogicBuild.Dockerfile
	data := []byte(dockerfile)
	dockerfileHash := md5.Sum(data)
	dockerfileMD5 := fmt.Sprintf("%x", dockerfileHash)
	tag := fmt.Sprintf("game-%d-build:%s", gameID, dockerfileMD5)

	// 获取游戏逻辑文件路径
	filePath := fmt.Sprintf("/var/hiper/games/%d/game_logic/gamelogic.zip", gameID)
	// 替换 Dockerfile 中的游戏逻辑文件路径
	dockerfile = strings.Replace(dockerfile, "GAME_LOGIC_PATH", filePath, -1)
	// 创建镜像
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	buildContext := strings.NewReader(dockerfile)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{tag},
	}

	resp, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
		status = 1
		return
	}
	defer resp.Body.Close()
	status = 0
	return
}

func buildAI(aiID uint) (status int) {
	// 获取AI
	ai, err := model.GetAiByID(aiID, true)
	if err != nil {
		log.Fatal(err)
	}
	// 获取AI的构建任务
	aiBuild := ai.Sdk.BuildAi
	// 获取AI的构建任务的 Dockerfile
	dockerfile := aiBuild.Dockerfile
	data := []byte(dockerfile)
	dockerfileHash := md5.Sum(data)
	dockerfileMD5 := fmt.Sprintf("%x", dockerfileHash)
	tag := fmt.Sprintf("ai-%d-build:%s", aiID, dockerfileMD5)
	// 获取AI文件路径
	filePath := fmt.Sprintf("/var/hiper/ais/%d/ai.zip", aiID)
	// 替换 Dockerfile 中的 AI 文件路径
	dockerfile = strings.Replace(dockerfile, "AI_PATH", filePath, -1)
	// 创建镜像
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	buildContext := strings.NewReader(dockerfile)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{tag},
	}

	resp, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer resp.Body.Close()
	return 0
}
