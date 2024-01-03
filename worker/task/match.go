package task

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func Match(matchID uint) {
	// 获取任务所需信息

	// 起容器，执行任务

	// 等待任务完成，获取任务输出，保存与修改相关信息（含 在 match_result 消息队列中发送消息，如果是 公开对局 的话）
	match, err := model.GetMatchByID(matchID, true)
	if err != nil {
		return
	}
	// 获取对应的赛事信息
	baseContest, err := model.GetBaseContestByID(match.BaseContestID)
	if err != nil {
		return
	}
	// 获取对应的build dockerfile
	game, err := model.GetGameByID(baseContest.GameID)
	if err != nil {
		return
	}
	gameLogic := game.GameLogic
	gameLogicBuild := gameLogic.Build
	dockerfile := gameLogicBuild.Dockerfile
	data := []byte(dockerfile)
	dockerfileHash := md5.Sum(data)
	dockerfileMD5 := fmt.Sprintf("%x", dockerfileHash)
	gameTag := fmt.Sprintf("game-%d-build:%s", baseContest.GameID, dockerfileMD5)

	// 根据tag获取镜像
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	imageExist := false // 镜像是否存在
	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == gameTag {
				imageExist = true
				break
			}
		}
	}

	if !imageExist {
		// 创建镜像
		output := BuildGameLogic(baseContest.GameID)
		if output != 0 {
			return
		}
	}

	for _, ai := range match.Ais {
		// 获取对应的run dockerfile
		sdk, err := model.GetSdkByID(ai.SdkID)
		if err != nil {
			return
		}
		runDockerfile := sdk.RunAi.Dockerfile
		data := []byte(runDockerfile)
		dockerfileHash := md5.Sum(data)
		dockerfileMD5 := fmt.Sprintf("%x", dockerfileHash)
		tag := fmt.Sprintf("sdk-%d-run:%s", ai.SdkID, dockerfileMD5)

		// 根据tag获取镜像
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		images, err := cli.ImageList(ctx, types.ImageListOptions{})
		if err != nil {
			panic(err)
		}
		imageExist := false // 镜像是否存在
		for _, image := range images {
			for _, repoTag := range image.RepoTags {
				if repoTag == tag {
					imageExist = true
					break
				}
			}
		}

		if !imageExist {
			// 创建镜像
			output := buildAI(ai.ID)
			if output != 0 {
				return
			}
		}
	}

	// 创建容器
	// 1. 创建网络
	// 2. 创建容器
	// 3. 启动容器
	// 4. 等待容器结束
	// 5. 获取容器输出
	// 6. 删除容器
	// 7. 删除网络

	// 1. 创建网络
	matchNetworkName := fmt.Sprintf("network-%d", matchID)
	matchNetwork, err := cli.NetworkCreate(ctx, matchNetworkName, types.NetworkCreate{})
	if err != nil {
		panic(err)
	}

	// 2. 创建GAME容器
	gameContainerName := fmt.Sprintf("game-%d", matchID)
	gameContainerConfig := &container.Config{
		Image: gameTag,
	}
	gameHostConfig := &container.HostConfig{}
	gameNetworkingConfig := &network.NetworkingConfig{}
	gameContainerResp, err := cli.ContainerCreate(ctx, gameContainerConfig, gameHostConfig, gameNetworkingConfig, gameContainerName)
	if err != nil {
		panic(err)
	}
	if err := cli.NetworkConnect(ctx, matchNetwork.ID, gameContainerResp.ID, nil); err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, gameContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// 3. 创建AI容器
	for AiIndex, ai := range match.Ais {
		aiContainerName := fmt.Sprintf("ai-%d-%d", matchID, AiIndex)
		//计算tag
		sdk, err := model.GetSdkByID(ai.SdkID)
		if err != nil {
			return
		}
		runDockerfile := sdk.RunAi.Dockerfile
		data := []byte(runDockerfile)
		dockerfileHash := md5.Sum(data)
		dockerfileMD5 := fmt.Sprintf("%x", dockerfileHash)
		tag := fmt.Sprintf("sdk-%d-run:%s", ai.SdkID, dockerfileMD5)
		aiContainerConfig := &container.Config{
			Image: tag,
		}
		aiHostConfig := &container.HostConfig{}
		aiNetworkingConfig := &network.NetworkingConfig{}
		aiContainerResp, err := cli.ContainerCreate(ctx, aiContainerConfig, aiHostConfig, aiNetworkingConfig, aiContainerName)
		if err != nil {
			panic(err)
		}
		if err := cli.NetworkConnect(ctx, matchNetwork.ID, aiContainerResp.ID, nil); err != nil {
			panic(err)
		}
		if err := cli.ContainerStart(ctx, aiContainerResp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}

	// 4. 等待容器结束
	statusCh, errCh := cli.ContainerWait(ctx, gameContainerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
	// 5. 获取容器输出
	gameContainerLogs, err := cli.ContainerLogs(ctx, gameContainerResp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	// 6. 删除容器
	if err := cli.ContainerRemove(ctx, gameContainerResp.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	for AiIndex, ai := range match.Ais {
		aiContainerName := fmt.Sprintf("ai-%d-%d", matchID, AiIndex)
		if err := cli.ContainerRemove(ctx, aiContainerName, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
	}

	// 7. 删除网络
	if err := cli.NetworkRemove(ctx, matchNetwork.ID); err != nil {
		panic(err)
	}

	// 8. 保存与修改相关信息（含 在 match_result 消息队列中发送消息，如果是 公开对局 的话）
	// TODO: 保存与修改相关信息（含 在 match_result 消息队列中发送消息，如果是 公开对局 的话）
}
