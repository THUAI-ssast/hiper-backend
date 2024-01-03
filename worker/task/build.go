package task

import (
	"fmt"
	"io"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func Build(domain repository.DomainType, id uint) error {
	dockerfile, err := getDockerfile(domain, repository.BuildOperation, id)
	if err != nil {
		return err
	}
	if dockerfile == "" {
		return nil
	}

	repository.StartBuildTask(domain, id)

	nameVersioned, err := prepareImage(repository.GameLogicDomain, repository.BuildOperation, id)
	if err != nil {
		repository.EndBuildTask(domain, id, model.TaskStateSystemError, err.Error())
		return err
	}
	/* TODO: 限制以下几个方面，以确保安全性：
		 * CPU、内存占用
		 * 运行时间
	     * 日志长度
	*/
	containerConfig := &container.Config{Image: nameVersioned}
	hostConfig := &container.HostConfig{
		Binds:      getBinds(domain, id),
		AutoRemove: true,
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		repository.EndBuildTask(domain, id, model.TaskStateSystemError, err.Error())
		return err
	}
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		repository.EndBuildTask(domain, id, model.TaskStateSystemError, err.Error())
		return err
	}
	statusCh, _ := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// Container has exited
	<-statusCh

	// Retrieve the container exit code
	respInspect, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		repository.EndBuildTask(domain, id, model.TaskStateSystemError, err.Error())
		return err
	}
	exitCode := respInspect.State.ExitCode
	// Retrieve container logs
	logOptions := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	logsReader, err := cli.ContainerLogs(ctx, resp.ID, logOptions)
	if err != nil {
		repository.EndBuildTask(domain, id, model.TaskStateSystemError, err.Error())
		return err
	}
	defer logsReader.Close()
	logs, _ := io.ReadAll(logsReader)

	// Interpret the exit code and update task state accordingly
	if exitCode != 0 {
		repository.EndBuildTask(domain, id, model.TaskStateInputError, fmt.Sprintf("Docker container exited with non-zero exit code: %d. Logs:\n%s", exitCode, string(logs)))
	}
	repository.EndBuildTask(domain, id, model.TaskStateFinished, fmt.Sprintf("Docker container exited successfully. Logs:\n%s", string(logs)))
	return nil
}
