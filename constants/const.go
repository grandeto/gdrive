package constants

import "time"

const Name = "gdrive"
const Version = "2.1.1"

const DefaultMaxFiles = 30
const DefaultMaxChanges = 100
const DefaultNameWidth = 40
const DefaultPathWidth = 60
const DefaultUploadChunkSize = 8 * 1024 * 1024
const DefaultTimeout = 5 * 60
const DefaultQuery = "trashed = false and 'me' in owners"
const DefaultShareRole = "reader"
const DefaultShareType = "anyone"

const MinCacheFileSize = 5 * 1024 * 1024

const MaxErrorRetries = 5

const DirectoryMimeType = "application/vnd.google-apps.folder"

const MaxDrawInterval = time.Second * 1
const MaxRateInterval = time.Second * 3

const DefaultIgnoreFile = ".gdriveignore"

type ModTime int

const (
	LocalLastModified ModTime = iota
	RemoteLastModified
	EqualModifiedTime
)

type LargestSize int

const (
	LocalLargestSize LargestSize = iota
	RemoteLargestSize
	EqualSize
)

type ConflictResolution int

const (
	NoResolution ConflictResolution = iota
	KeepLocal
	KeepRemote
	KeepLargest
)

const TimeoutTimerInterval = time.Second * 10

const TokenFilename = "token_v2.json"
const DefaultCacheFileName = "file_cache.json"
const AuthFileName = "gdrive_auth_value.txt"

const HomeDir = "/home"
