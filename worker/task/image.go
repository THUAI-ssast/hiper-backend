package task

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
	"github.com/docker/docker/api/types"
)

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
