package adapters

import (
	"context"
	"fmt"
	"sync"
)

type FilesApiStub struct {
	lock  sync.Mutex
	files map[string]string
}

func (c *FilesApiStub) UploadFile(ctx context.Context, fileID string, fileContent string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.files == nil {
		c.files = make(map[string]string)
	}

	c.files[fileID] = fileContent

	return nil
}

func (c *FilesApiStub) DownloadFile(ctx context.Context, fileID string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.files == nil {
		c.files = make(map[string]string)
	}

	fileContent, ok := c.files[fileID]
	if !ok {
		return "", fmt.Errorf("file %s not found", fileID)
	}

	return fileContent, nil
}
