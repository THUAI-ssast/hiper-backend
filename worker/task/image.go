package task

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
	"github.com/docker/docker/api/types"
)

func getDockerfile(domain repository.DomainType, operation repository.OperationType, id uint) (dockerfile string, err error) {
	switch domain {
	case repository.GameLogicDomain:
		var g model.Game
		g, err = model.GetGameByID(id, fmt.Sprintf("game_logic_%s_dockerfile", operation))
		if err != nil {
			return
		}
		switch operation {
		case repository.BuildOperation:
			dockerfile = g.GameLogic.Build.Dockerfile
		case repository.RunOperation:
			dockerfile = g.GameLogic.Run.Dockerfile
		}
	case repository.AiDomain:
		var ai model.Ai
		ai, err = model.GetAiByID(id, true)
		if err != nil {
			return
		}
		switch operation {
		case repository.BuildOperation:
			dockerfile = ai.Sdk.BuildAi.Dockerfile
		case repository.RunOperation:
			dockerfile = ai.Sdk.RunAi.Dockerfile
		}
	}
	return
}

func prepareImage(domain repository.DomainType, operation repository.OperationType, id uint) (string, error) {
	dockerfile, err := getDockerfile(domain, operation, id)
	if err != nil {
		return "", err
	}
	dockerfileHash := md5.Sum([]byte(dockerfile))
	name := fmt.Sprintf("%s-%d-%s", domain, id, operation)
	tag := fmt.Sprintf("%x", dockerfileHash)
	nameVersioned := fmt.Sprintf("%s:%x", name, dockerfileHash)
	// Check if image exists
	if _, _, err := cli.ImageInspectWithRaw(ctx, nameVersioned); err != nil {
		// image not exists, build it
		if err = buildImage(dockerfile, name, tag, domain, operation, id); err != nil {
			return "", err
		}
	}
	return nameVersioned, nil
}

// param name and tag: `--tag <name>:<tag>`
func buildImage(dockerfileContent, name string, tag string, domain repository.DomainType, operation repository.OperationType, id uint) error {
	repository.StartBuildImageTask(domain, operation, id)

	// Add Dockerfile to the in-memory tar archive
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)
	tarHeader := &tar.Header{
		Name: "Dockerfile",
		Mode: 0644,
		Size: int64(len(dockerfileContent)),
	}
	if err := tarWriter.WriteHeader(tarHeader); err != nil {
		repository.EndBuildImageTask(domain, operation, id, model.TaskStateSystemError, err.Error())
		return err
	}
	if _, err := tarWriter.Write([]byte(dockerfileContent)); err != nil {
		repository.EndBuildImageTask(domain, operation, id, model.TaskStateSystemError, err.Error())
		return err
	}
	if err := tarWriter.Close(); err != nil {
		repository.EndBuildImageTask(domain, operation, id, model.TaskStateSystemError, err.Error())
		return err
	}

	// Build image
	buildResponse, err := cli.ImageBuild(ctx, tarBuffer, types.ImageBuildOptions{
		Tags: []string{name, name + ":" + tag},
	})
	if err != nil {
		buildOutput, readErr := io.ReadAll(buildResponse.Body)
		if readErr != nil {
			repository.EndBuildImageTask(domain, operation, id, model.TaskStateSystemError, err.Error())
			return err
		}
		// truncate build output
		if len(buildOutput) > 2048 {
			buildOutput = buildOutput[:2048]
		}
		repository.EndBuildImageTask(domain, operation, id, model.TaskStateInputError, string(buildOutput))
	}
	buildResponse.Body.Close()

	repository.EndBuildImageTask(domain, operation, id, model.TaskStateFinished, "")
	return nil
}
