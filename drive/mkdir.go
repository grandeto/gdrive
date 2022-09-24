package drive

import (
	"fmt"
	"io"

	"github.com/grandeto/gdrive/constants"
	"google.golang.org/api/drive/v3"
)

type MkdirArgs struct {
	Out         io.Writer
	Name        string
	Description string
	Parents     []string
}

func (self *Drive) Mkdir(args MkdirArgs) error {
	f, err := self.mkdir(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(args.Out, "Directory %s created\n", f.Id)
	return nil
}

func (self *Drive) mkdir(args MkdirArgs) (*drive.File, error) {
	dstFile := &drive.File{
		Name:        args.Name,
		Description: args.Description,
		MimeType:    constants.DirectoryMimeType,
	}

	// Set parent folders
	dstFile.Parents = args.Parents

	// Create directory
	f, err := self.service.Files.Create(dstFile).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory: %s", err)
	}

	return f, nil
}
