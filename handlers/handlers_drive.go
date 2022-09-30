package handlers

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/grandeto/gdrive/auth"
	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/compare"
	"github.com/grandeto/gdrive/constants"
	"github.com/grandeto/gdrive/drive"
	"github.com/grandeto/gdrive/util"
	_ "github.com/joho/godotenv/autoload"
)

var ClientId = os.Getenv("CLIENT_ID")
var ClientSecret = os.Getenv("CLIENT_SECRET")
var OsUser = os.Getenv("OS_USER")

func ListHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).List(drive.ListFilesArgs{
		Out:         os.Stdout,
		MaxFiles:    args.Int64("maxFiles"),
		NameWidth:   args.Int64("nameWidth"),
		Query:       args.String("query"),
		SortOrder:   args.String("sortOrder"),
		SkipHeader:  args.Bool("skipHeader"),
		SizeInBytes: args.Bool("sizeInBytes"),
		AbsPath:     args.Bool("absPath"),
	})
	util.CheckErr(err)
}

func ListChangesHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListChanges(drive.ListChangesArgs{
		Out:        os.Stdout,
		PageToken:  args.String("pageToken"),
		MaxChanges: args.Int64("maxChanges"),
		Now:        args.Bool("now"),
		NameWidth:  args.Int64("nameWidth"),
		SkipHeader: args.Bool("skipHeader"),
	})
	util.CheckErr(err)
}

func DownloadHandler(ctx cli.Context) {
	args := ctx.Args()
	checkDownloadArgs(args)
	err := newDrive(args).Download(drive.DownloadArgs{
		Out:       os.Stdout,
		Id:        args.String("fileId"),
		Force:     args.Bool("force"),
		Skip:      args.Bool("skip"),
		Path:      args.String("path"),
		Delete:    args.Bool("delete"),
		Recursive: args.Bool("recursive"),
		Stdout:    args.Bool("stdout"),
		Progress:  progressWriter(args.Bool("noProgress")),
		Timeout:   durationInSeconds(args.Int64("timeout")),
	})
	util.CheckErr(err)
}

func DownloadQueryHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DownloadQuery(drive.DownloadQueryArgs{
		Out:       os.Stdout,
		Query:     args.String("query"),
		Force:     args.Bool("force"),
		Skip:      args.Bool("skip"),
		Recursive: args.Bool("recursive"),
		Path:      args.String("path"),
		Progress:  progressWriter(args.Bool("noProgress")),
	})
	util.CheckErr(err)
}

func DownloadSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	cachePath := filepath.Join(args.String("configDir"), constants.DefaultCacheFileName)
	err := newDrive(args).DownloadSync(drive.DownloadSyncArgs{
		Out:              os.Stdout,
		Progress:         progressWriter(args.Bool("noProgress")),
		Path:             args.String("path"),
		RootId:           args.String("fileId"),
		DryRun:           args.Bool("dryRun"),
		DeleteExtraneous: args.Bool("deleteExtraneous"),
		Timeout:          durationInSeconds(args.Int64("timeout")),
		Resolution:       conflictResolution(args),
		Comparer:         compare.NewCachedMd5Comparer(cachePath),
	})
	util.CheckErr(err)
}

func DownloadRevisionHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DownloadRevision(drive.DownloadRevisionArgs{
		Out:        os.Stdout,
		FileId:     args.String("fileId"),
		RevisionId: args.String("revId"),
		Force:      args.Bool("force"),
		Stdout:     args.Bool("stdout"),
		Path:       args.String("path"),
		Progress:   progressWriter(args.Bool("noProgress")),
		Timeout:    durationInSeconds(args.Int64("timeout")),
	})
	util.CheckErr(err)
}

func UploadHandler(ctx cli.Context) {
	args := ctx.Args()
	checkUploadArgs(args)
	err := newDrive(args).Upload(drive.UploadArgs{
		Out:         os.Stdout,
		Progress:    progressWriter(args.Bool("noProgress")),
		Path:        args.String("path"),
		Name:        args.String("name"),
		Description: args.String("description"),
		Parents:     args.StringSlice("parent"),
		Mime:        args.String("mime"),
		Recursive:   args.Bool("recursive"),
		Share:       args.Bool("share"),
		Delete:      args.Bool("delete"),
		ChunkSize:   args.Int64("chunksize"),
		Timeout:     durationInSeconds(args.Int64("timeout")),
	})
	util.CheckErr(err)
}

func UploadStdinHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).UploadStream(drive.UploadStreamArgs{
		Out:         os.Stdout,
		In:          os.Stdin,
		Name:        args.String("name"),
		Description: args.String("description"),
		Parents:     args.StringSlice("parent"),
		Mime:        args.String("mime"),
		Share:       args.Bool("share"),
		ChunkSize:   args.Int64("chunksize"),
		Timeout:     durationInSeconds(args.Int64("timeout")),
		Progress:    progressWriter(args.Bool("noProgress")),
	})
	util.CheckErr(err)
}

func UploadSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	cachePath := filepath.Join(args.String("configDir"), constants.DefaultCacheFileName)
	err := newDrive(args).UploadSync(drive.UploadSyncArgs{
		Out:              os.Stdout,
		Progress:         progressWriter(args.Bool("noProgress")),
		Path:             args.String("path"),
		RootId:           args.String("fileId"),
		DryRun:           args.Bool("dryRun"),
		DeleteExtraneous: args.Bool("deleteExtraneous"),
		ChunkSize:        args.Int64("chunksize"),
		Timeout:          durationInSeconds(args.Int64("timeout")),
		Resolution:       conflictResolution(args),
		Comparer:         compare.NewCachedMd5Comparer(cachePath),
	})
	util.CheckErr(err)
}

func UpdateHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Update(drive.UpdateArgs{
		Out:         os.Stdout,
		Id:          args.String("fileId"),
		Path:        args.String("path"),
		Name:        args.String("name"),
		Description: args.String("description"),
		Parents:     args.StringSlice("parent"),
		Mime:        args.String("mime"),
		Progress:    progressWriter(args.Bool("noProgress")),
		ChunkSize:   args.Int64("chunksize"),
		Timeout:     durationInSeconds(args.Int64("timeout")),
	})
	util.CheckErr(err)
}

func InfoHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Info(drive.FileInfoArgs{
		Out:         os.Stdout,
		Id:          args.String("fileId"),
		SizeInBytes: args.Bool("sizeInBytes"),
	})
	util.CheckErr(err)
}

func ImportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Import(drive.ImportArgs{
		Mime:     args.String("mime"),
		Out:      os.Stdout,
		Path:     args.String("path"),
		Parents:  args.StringSlice("parent"),
		Progress: progressWriter(args.Bool("noProgress")),
	})
	util.CheckErr(err)
}

func ExportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Export(drive.ExportArgs{
		Out:        os.Stdout,
		Id:         args.String("fileId"),
		Mime:       args.String("mime"),
		PrintMimes: args.Bool("printMimes"),
		Force:      args.Bool("force"),
	})
	util.CheckErr(err)
}

func ListRevisionsHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListRevisions(drive.ListRevisionsArgs{
		Out:         os.Stdout,
		Id:          args.String("fileId"),
		NameWidth:   args.Int64("nameWidth"),
		SizeInBytes: args.Bool("sizeInBytes"),
		SkipHeader:  args.Bool("skipHeader"),
	})
	util.CheckErr(err)
}

func MkdirHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Mkdir(drive.MkdirArgs{
		Out:         os.Stdout,
		Name:        args.String("name"),
		Description: args.String("description"),
		Parents:     args.StringSlice("parent"),
	})
	util.CheckErr(err)
}

func ShareHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Share(drive.ShareArgs{
		Out:          os.Stdout,
		FileId:       args.String("fileId"),
		Role:         args.String("role"),
		Type:         args.String("type"),
		Email:        args.String("email"),
		Domain:       args.String("domain"),
		Discoverable: args.Bool("discoverable"),
	})
	util.CheckErr(err)
}

func ShareListHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListPermissions(drive.ListPermissionsArgs{
		Out:    os.Stdout,
		FileId: args.String("fileId"),
	})
	util.CheckErr(err)
}

func ShareRevokeHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).RevokePermission(drive.RevokePermissionArgs{
		Out:          os.Stdout,
		FileId:       args.String("fileId"),
		PermissionId: args.String("permissionId"),
	})
	util.CheckErr(err)
}

func DeleteHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Delete(drive.DeleteArgs{
		Out:       os.Stdout,
		Id:        args.String("fileId"),
		Recursive: args.Bool("recursive"),
	})
	util.CheckErr(err)
}

func ListSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListSync(drive.ListSyncArgs{
		Out:        os.Stdout,
		SkipHeader: args.Bool("skipHeader"),
	})
	util.CheckErr(err)
}

func ListRecursiveSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListRecursiveSync(drive.ListRecursiveSyncArgs{
		Out:         os.Stdout,
		RootId:      args.String("fileId"),
		SkipHeader:  args.Bool("skipHeader"),
		PathWidth:   args.Int64("pathWidth"),
		SizeInBytes: args.Bool("sizeInBytes"),
		SortOrder:   args.String("sortOrder"),
	})
	util.CheckErr(err)
}

func DeleteRevisionHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DeleteRevision(drive.DeleteRevisionArgs{
		Out:        os.Stdout,
		FileId:     args.String("fileId"),
		RevisionId: args.String("revId"),
	})
	util.CheckErr(err)
}

func AboutHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).About(drive.AboutArgs{
		Out:         os.Stdout,
		SizeInBytes: args.Bool("sizeInBytes"),
	})
	util.CheckErr(err)
}

func AboutImportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).AboutImport(drive.AboutImportArgs{
		Out: os.Stdout,
	})
	util.CheckErr(err)
}

func AboutExportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).AboutExport(drive.AboutExportArgs{
		Out: os.Stdout,
	})
	util.CheckErr(err)
}

func getOauthClient(args cli.Arguments) (*http.Client, error) {
	if args.String("refreshToken") != "" && args.String("accessToken") != "" {
		util.ExitF("Access token not needed when refresh token is provided")
	}

	if args.String("refreshToken") != "" {
		return auth.NewRefreshTokenClient(ClientId, ClientSecret, args.String("refreshToken")), nil
	}

	if args.String("accessToken") != "" {
		return auth.NewAccessTokenClient(ClientId, ClientSecret, args.String("accessToken")), nil
	}

	configDir := getConfigDir(args)

	if args.String("serviceAccount") != "" {
		serviceAccountPath := util.ConfigFilePath(configDir, args.String("serviceAccount"))
		serviceAccountClient, err := auth.NewServiceAccountClient(serviceAccountPath)
		if err != nil {
			return nil, err
		}
		return serviceAccountClient, nil
	}

	tokenPath := util.ConfigFilePath(configDir, constants.TokenFilename)
	return auth.NewFileSourceClient(ClientId, ClientSecret, tokenPath, authCodePrompt)
}

func getConfigDir(args cli.Arguments) string {
	// Use dir from environment var if present
	if os.Getenv("GDRIVE_CONFIG_DIR") != "" {
		return os.Getenv("GDRIVE_CONFIG_DIR")
	}
	return args.String("configDir")
}

func newDrive(args cli.Arguments) *drive.Drive {
	oauth, err := getOauthClient(args)
	if err != nil {
		util.ExitF("Failed getting oauth client: %s", err.Error())
	}

	client, err := drive.New(oauth)
	if err != nil {
		util.ExitF("Failed getting drive: %s", err.Error())
	}

	return client
}

func authCodePrompt(url string) func() string {
	return func() string {
		authFile := util.GetAuthFileNamePath()

		f, err := os.Create(authFile)

		if err != nil {
			util.ExitF(err.Error())
		}

		defer f.Close()

		s1 := "Authentication needed"
		s2 := "Go to the following url in your browser:"
		s3 := url
		s4 := "Enter verification code here, right after the collon:"

		authMsg := fmt.Sprintf("%s\n%s\n%s\n\n%s", s1, s2, s3, s4)

		f.WriteString(authMsg)

		log.Printf("Check %s for authentication instructions\n\n", authFile)

		// creates a new file watcher
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			util.ExitF(err.Error())
		}

		defer watcher.Close()

		done := make(chan string, 1)

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == authFile && event.Op.String() == "WRITE" {
						ff, _ := os.Open(authFile)
						fileScanner := bufio.NewScanner(ff)

						fileScanner.Split(bufio.ScanLines)

						for fileScanner.Scan() {
							if strings.Contains(fileScanner.Text(), s4) {
								code := strings.Split(fileScanner.Text(), s4)[1]
								done <- code
							}
						}
					}
				case _ = <-watcher.Errors:
					// log.Println("ERROR", err)
				}
			}
		}()

		if err := watcher.Add(authFile); err != nil {
			util.ExitF(err.Error())
		}

		code := <-done

		return code
	}
}

func progressWriter(discard bool) io.Writer {
	if discard {
		return ioutil.Discard
	}
	return os.Stderr
}

func durationInSeconds(seconds int64) time.Duration {
	return time.Second * time.Duration(seconds)
}

func conflictResolution(args cli.Arguments) constants.ConflictResolution {
	keepLocal := args.Bool("keepLocal")
	keepRemote := args.Bool("keepRemote")
	keepLargest := args.Bool("keepLargest")

	if (keepLocal && keepRemote) || (keepLocal && keepLargest) || (keepRemote && keepLargest) {
		util.ExitF("Only one conflict resolution flag can be given")
	}

	if keepLocal {
		return constants.KeepLocal
	}

	if keepRemote {
		return constants.KeepRemote
	}

	if keepLargest {
		return constants.KeepLargest
	}

	return constants.NoResolution
}

func checkUploadArgs(args cli.Arguments) {
	if args.Bool("recursive") && args.Bool("delete") {
		util.ExitF("--delete is not allowed for recursive uploads")
	}

	if args.Bool("recursive") && args.Bool("share") {
		util.ExitF("--share is not allowed for recursive uploads")
	}
}

func checkDownloadArgs(args cli.Arguments) {
	if args.Bool("recursive") && args.Bool("delete") {
		util.ExitF("--delete is not allowed for recursive downloads")
	}
}
