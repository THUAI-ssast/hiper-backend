package task

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/THUAI-ssast/hiper-backend/worker/mq"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
)

func Match(matchID uint) error {
	repository.StartMatchTask(matchID)
	// Get necessary information
	match, err := model.GetMatchByID(matchID, true)
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}
	aiIDs := make([]uint, len(match.Ais))
	for _, ai := range match.Ais {
		prepareImage(repository.AiDomain, repository.RunOperation, ai.ID)
		aiIDs = append(aiIDs, ai.ID)
	}
	prepareImage(repository.GameLogicDomain, repository.RunOperation, match.GameID)

	// Create a network for the match
	networkName := getNetworkName(matchID)
	_, err = cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}

	// Create a named pipe for game logic to communicate with judger
	pipePath := getPipePath(matchID)
	err = syscall.Mkfifo(pipePath, 0666)
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}

	// Run game logic container
	err = runGameLogicContainer(aiIDs, match, pipePath, networkName, matchID)
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}

	var pipe *os.File
	defer cleanMatch(matchID, pipe)

	// Listen to the pipe to execute commands from game logic
	pipe, err = os.Open(pipePath)
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}
	err = executeCommandsFromGameLogic(pipe, match, aiIDs)
	if err != nil {
		repository.EndMatchTask(matchID, model.TaskStateSystemError)
		return err
	}

	repository.EndMatchTask(matchID, model.TaskStateFinished)
	return nil
}

func runGameLogicContainer(aiIDs []uint, match model.Match, pipePath string, networkName string, matchID uint) error {
	/* TODO: 限制以下几个方面，以确保安全性：
	 * CPU、内存占用
	 * 运行时间
	 */
	aiIDsJSON, _ := json.Marshal(aiIDs)
	gameLogicContainerConfig := &container.Config{
		Image: getImage(repository.GameLogicDomain, repository.RunOperation, match.GameID),
		Cmd:   []string{string(aiIDsJSON), match.ExtraInfo},
	}
	binds := getBinds(repository.GameLogicDomain, match.GameID)
	binds = append(binds, fmt.Sprintf("%s:%s", pipePath, "/tmp/pipe"))
	gameLogicHostConfig := &container.HostConfig{
		Binds:       binds,
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(networkName),
	}
	containerName := "game_logic"
	gameLogicResp, err := cli.ContainerCreate(ctx, gameLogicContainerConfig, gameLogicHostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}
	if err = cli.ContainerStart(ctx, gameLogicResp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func executeCommandsFromGameLogic(pipe *os.File, match model.Match, aiIDs []uint) error {
	reader := bufio.NewReader(pipe)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		components := strings.Split(strings.TrimSpace(command), " ")
		commandType := components[0]
		args := components[1:]
		switch commandType {
		case "create_player":
			containerName := args[0]
			aiIndexInt, _ := strconv.Atoi(args[1])
			aiIndex := uint(aiIndexInt)
			createPlayer(containerName, aiIndex, aiIDs, match.ID)
		case "remove_player":
			containerName := args[0]
			removePlayer(containerName)
		case "end_match":
			endInfoJSON := args[0]
			endMatch(endInfoJSON, match)
			return nil
		default:
			continue // ignore unknown commands
		}
	}
}

func createPlayer(containerName string, aiIndex uint, aiIDs []uint, matchID uint) {
	/* TODO: 限制以下几个方面，以确保安全性：
		 * CPU、内存占用
		 * 运行时间
	     * 日志长度
	*/
	aiID := aiIDs[aiIndex]
	containerConfig := &container.Config{
		Image: getImage(repository.AiDomain, repository.RunOperation, aiID),
	}
	hostConfig := &container.HostConfig{
		Binds:       getBinds(repository.AiDomain, aiID),
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(getNetworkName(matchID)),
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Println(err)
	}
	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println(err)
	}
	go waitToAppendPlayerLog(resp.ID, matchID, aiIndex)
}

func waitToAppendPlayerLog(containerID string, matchID uint, aiIndex uint) {
	// Wait for container to exit
	statusCh, _ := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	// Container has exited
	<-statusCh

	// Retrieve container logs
	logOptions := types.ContainerLogsOptions{ShowStderr: true}
	logsReader, err := cli.ContainerLogs(ctx, containerID, logOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer logsReader.Close()
	logs, _ := io.ReadAll(logsReader)
	// Append to corresponding log file
	logPath := fmt.Sprintf("/var/hiper/matches/%d/player_%d.log", matchID, aiIndex)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer logFile.Close()
	header := "====================\n"
	logs = append([]byte(header), logs...)
	logFile.Write(logs)
}

func removePlayer(containerName string) {
	if err := cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
		log.Println(err)
	}
}

// endMatch ends a match that finished successfully
func endMatch(endInfoJSON string, match model.Match) {
	var endInfo EndInfo
	if err := json.Unmarshal([]byte(endInfoJSON), &endInfo); err != nil {
		log.Println(err)
		return
	}
	// Update match
	model.UpdateMatchByID(match.ID, map[string]interface{}{
		"scores": endInfo.Scores,
	})
	if match.MatchType == model.MatchTypePublic {
		replayPath := fmt.Sprintf("/var/hiper/matches/%d/replay.json", match.ID)
		replayFile, err := os.Open(replayPath)
		if err != nil {
			log.Println(err)
			return
		}
		defer replayFile.Close()
		replay, _ := io.ReadAll(replayFile)
		mq.PublishMatchResult(match.ID, string(replay))
	}
}

type EndInfo struct {
	Scores []int
}

// cleanMatch cleans up all resources used by a match
func cleanMatch(matchID uint, pipe *os.File) {
	// Stop all containers in the match network
	networkName := getNetworkName(matchID)
	network, err := cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	for _, curContainer := range network.Containers {
		if err := cli.ContainerStop(ctx, curContainer.Name, container.StopOptions{}); err != nil {
			log.Println(err)
		}
	}
	// Remove network
	if err := cli.NetworkRemove(ctx, networkName); err != nil {
		log.Println(err)
	}
	// Remove pipe
	if pipe != nil {
		pipe.Close()
		pipePath := getPipePath(matchID)
		if err := os.Remove(pipePath); err != nil {
			log.Println(err)
		}
	}
}

func getNetworkName(matchID uint) string {
	return fmt.Sprintf("match-%d", matchID)
}

func getPipePath(matchID uint) string {
	return fmt.Sprintf("/tmp/match-%d-pipe", matchID)
}
