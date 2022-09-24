package loader

import (
	"fmt"

	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/constants"
	"github.com/grandeto/gdrive/handlers"
	"github.com/grandeto/gdrive/util"
)

var DefaultConfigDir = util.GetDefaultConfigDir()

func LoadGlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:         "configDir",
			Patterns:     []string{"-c", "--config"},
			Description:  fmt.Sprintf("Application path, default: %s", DefaultConfigDir),
			DefaultValue: DefaultConfigDir,
		},
		cli.StringFlag{
			Name:        "refreshToken",
			Patterns:    []string{"--refresh-token"},
			Description: "Oauth refresh token used to get access token (for advanced users)",
		},
		cli.StringFlag{
			Name:        "accessToken",
			Patterns:    []string{"--access-token"},
			Description: "Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)",
		},
		cli.StringFlag{
			Name:        "serviceAccount",
			Patterns:    []string{"--service-account"},
			Description: "Oauth service account filename, used for server to server communication without user interaction (filename path is relative to config dir)",
		},
	}
}

func LoadHandlers(globalFlags []cli.Flag) []*cli.Handler {
	return []*cli.Handler{
		&cli.Handler{
			Pattern:     "[global] list [options]",
			Description: "List files",
			Callback:    handlers.ListHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "maxFiles",
						Patterns:     []string{"-m", "--max"},
						Description:  fmt.Sprintf("Max files to list, default: %d", constants.DefaultMaxFiles),
						DefaultValue: constants.DefaultMaxFiles,
					},
					cli.StringFlag{
						Name:         "query",
						Patterns:     []string{"-q", "--query"},
						Description:  fmt.Sprintf(`Default query: "%s". See https://developers.google.com/drive/search-parameters`, constants.DefaultQuery),
						DefaultValue: constants.DefaultQuery,
					},
					cli.StringFlag{
						Name:        "sortOrder",
						Patterns:    []string{"--order"},
						Description: "Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy",
					},
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", constants.DefaultNameWidth),
						DefaultValue: constants.DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "absPath",
						Patterns:    []string{"--absolute"},
						Description: "Show absolute path to file (will only show path from first parent)",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] download [options] <fileId>",
			Description: "Download file or directory",
			Callback:    handlers.DownloadHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skip",
						Patterns:    []string{"-s", "--skip"},
						Description: "Skip existing files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Download directory recursively, documents will be skipped",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.BoolFlag{
						Name:        "delete",
						Patterns:    []string{"--delete"},
						Description: "Delete remote file when download is successful",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "stdout",
						Patterns:    []string{"--stdout"},
						Description: "Write file content to stdout",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] download query [options] <query>",
			Description: "Download all files and directories matching query",
			Callback:    handlers.DownloadQueryHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skip",
						Patterns:    []string{"-s", "--skip"},
						Description: "Skip existing files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Download directories recursively, documents will be skipped",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] upload [options] <path>",
			Description: "Upload file or directory",
			Callback:    handlers.UploadHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Upload directory recursively",
						OmitValue:   true,
					},
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "name",
						Patterns:    []string{"--name"},
						Description: "Filename",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.BoolFlag{
						Name:        "share",
						Patterns:    []string{"--share"},
						Description: "Share file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "delete",
						Patterns:    []string{"--delete"},
						Description: "Delete local file when upload is successful",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", constants.DefaultUploadChunkSize),
						DefaultValue: constants.DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] upload - [options] <name>",
			Description: "Upload file from stdin",
			Callback:    handlers.UploadStdinHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", constants.DefaultUploadChunkSize),
						DefaultValue: constants.DefaultUploadChunkSize,
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.BoolFlag{
						Name:        "share",
						Patterns:    []string{"--share"},
						Description: "Share file",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] update [options] <fileId> <path>",
			Description: "Update file, this creates a new revision of the file",
			Callback:    handlers.UpdateHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "name",
						Patterns:    []string{"--name"},
						Description: "Filename",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", constants.DefaultUploadChunkSize),
						DefaultValue: constants.DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] info [options] <fileId>",
			Description: "Show file info",
			Callback:    handlers.InfoHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Show size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] mkdir [options] <name>",
			Description: "Create directory",
			Callback:    handlers.MkdirHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id of created directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "Directory description",
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share [options] <fileId>",
			Description: "Share file or directory",
			Callback:    handlers.ShareHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringFlag{
						Name:         "role",
						Patterns:     []string{"--role"},
						Description:  fmt.Sprintf("Share role: owner/writer/commenter/reader, default: %s", constants.DefaultShareRole),
						DefaultValue: constants.DefaultShareRole,
					},
					cli.StringFlag{
						Name:         "type",
						Patterns:     []string{"--type"},
						Description:  fmt.Sprintf("Share type: user/group/domain/anyone, default: %s", constants.DefaultShareType),
						DefaultValue: constants.DefaultShareType,
					},
					cli.StringFlag{
						Name:        "email",
						Patterns:    []string{"--email"},
						Description: "The email address of the user or group to share the file with. Requires 'user' or 'group' as type",
					},
					cli.StringFlag{
						Name:        "domain",
						Patterns:    []string{"--domain"},
						Description: "The name of Google Apps domain. Requires 'domain' as type",
					},
					cli.BoolFlag{
						Name:        "discoverable",
						Patterns:    []string{"--discoverable"},
						Description: "Make file discoverable by search engines",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "revoke",
						Patterns:    []string{"--revoke"},
						Description: "Delete all sharing permissions (owner roles will be skipped)",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share list <fileId>",
			Description: "List files permissions",
			Callback:    handlers.ShareListHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share revoke <fileId> <permissionId>",
			Description: "Revoke permission",
			Callback:    handlers.ShareRevokeHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] delete [options] <fileId>",
			Description: "Delete file or directory",
			Callback:    handlers.DeleteHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Delete directory and all it's content",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync list [options]",
			Description: "List all syncable directories on drive",
			Callback:    handlers.ListSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync content [options] <fileId>",
			Description: "List content of syncable directory",
			Callback:    handlers.ListRecursiveSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringFlag{
						Name:        "sortOrder",
						Patterns:    []string{"--order"},
						Description: "Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy",
					},
					cli.IntFlag{
						Name:         "pathWidth",
						Patterns:     []string{"--path-width"},
						Description:  fmt.Sprintf("Width of path column, default: %d, minimum: 9, use 0 for full width", constants.DefaultPathWidth),
						DefaultValue: constants.DefaultPathWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync download [options] <fileId> <path>",
			Description: "Sync drive directory to local directory",
			Callback:    handlers.DownloadSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "keepRemote",
						Patterns:    []string{"--keep-remote"},
						Description: "Keep remote file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLocal",
						Patterns:    []string{"--keep-local"},
						Description: "Keep local file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLargest",
						Patterns:    []string{"--keep-largest"},
						Description: "Keep largest file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "deleteExtraneous",
						Patterns:    []string{"--delete-extraneous"},
						Description: "Delete extraneous local files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "dryRun",
						Patterns:    []string{"--dry-run"},
						Description: "Show what would have been transferred",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync upload [options] <path> <fileId>",
			Description: "Sync local directory to drive",
			Callback:    handlers.UploadSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "keepRemote",
						Patterns:    []string{"--keep-remote"},
						Description: "Keep remote file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLocal",
						Patterns:    []string{"--keep-local"},
						Description: "Keep local file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLargest",
						Patterns:    []string{"--keep-largest"},
						Description: "Keep largest file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "deleteExtraneous",
						Patterns:    []string{"--delete-extraneous"},
						Description: "Delete extraneous remote files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "dryRun",
						Patterns:    []string{"--dry-run"},
						Description: "Show what would have been transferred",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", constants.DefaultUploadChunkSize),
						DefaultValue: constants.DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] changes [options]",
			Description: "List file changes",
			Callback:    handlers.ListChangesHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "maxChanges",
						Patterns:     []string{"-m", "--max"},
						Description:  fmt.Sprintf("Max changes to list, default: %d", constants.DefaultMaxChanges),
						DefaultValue: constants.DefaultMaxChanges,
					},
					cli.StringFlag{
						Name:         "pageToken",
						Patterns:     []string{"--since"},
						Description:  fmt.Sprintf("Page token to start listing changes from"),
						DefaultValue: "1",
					},
					cli.BoolFlag{
						Name:        "now",
						Patterns:    []string{"--now"},
						Description: fmt.Sprintf("Get latest page token"),
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", constants.DefaultNameWidth),
						DefaultValue: constants.DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision list [options] <fileId>",
			Description: "List file revisions",
			Callback:    handlers.ListRevisionsHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", constants.DefaultNameWidth),
						DefaultValue: constants.DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision download [options] <fileId> <revId>",
			Description: "Download revision",
			Callback:    handlers.DownloadRevisionHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "stdout",
						Patterns:    []string{"--stdout"},
						Description: "Write file content to stdout",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", constants.DefaultTimeout),
						DefaultValue: constants.DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision delete <fileId> <revId>",
			Description: "Delete file revision",
			Callback:    handlers.DeleteRevisionHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] import [options] <path>",
			Description: "Upload and convert file to a google document, see 'about import' for available conversions",
			Callback:    handlers.ImportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Mime type of imported file",
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] export [options] <fileId>",
			Description: "Export a google document",
			Callback:    handlers.ExportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Mime type of exported file",
					},
					cli.BoolFlag{
						Name:        "printMimes",
						Patterns:    []string{"--print-mimes"},
						Description: "Print available mime types for given file",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about [options]",
			Description: "Google drive metadata, quota usage",
			Callback:    handlers.AboutHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Show size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about import",
			Description: "Show supported import formats",
			Callback:    handlers.AboutImportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about export",
			Description: "Show supported export formats",
			Callback:    handlers.AboutExportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "version",
			Description: "Print application version",
			Callback:    handlers.PrintVersion,
		},
		&cli.Handler{
			Pattern:     "help",
			Description: "Print help",
			Callback:    handlers.PrintHelp,
		},
		&cli.Handler{
			Pattern:     "help <command>",
			Description: "Print command help",
			Callback:    handlers.PrintCommandHelp,
		},
		&cli.Handler{
			Pattern:     "help <command> <subcommand>",
			Description: "Print subcommand help",
			Callback:    handlers.PrintSubCommandHelp,
		},
	}
}
