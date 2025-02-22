package hdfs

import (
	"os"
	"syscall"
)

const (
	fileNotFoundException        = "java.io.FileNotFoundException"
	permissionDeniedException    = "org.apache.hadoop.security.AccessControlException"
	pathIsNotEmptyDirException   = "org.apache.hadoop.fs.PathIsNotEmptyDirectoryException"
	fileAlreadyExistsException   = "org.apache.hadoop.fs.FileAlreadyExistsException"
	alreadyBeingCreatedException = "org.apache.hadoop.hdfs.protocol.AlreadyBeingCreatedException"
	invalidPathException         = "org.apache.hadoop.fs.InvalidPathException"
	safeModeException            = "org.apache.hadoop.hdfs.server.namenode.SafeModeException"
	dSQuotaExceededException     = "org.apache.hadoop.hdfs.protocol.DSQuotaExceededException"
	nSQuotaExceededException     = "org.apache.hadoop.hdfs.protocol.NSQuotaExceededException"
	parentNotDirectoryException  = "org.apache.hadoop.fs.ParentNotDirectoryException"
	UnresolvedLinkException      = "org.apache.hadoop.fs.UnresolvedLinkException"
	NotReplicatedYetException    = "org.apache.hadoop.hdfs.server.namenode.NotReplicatedYetException"
	illegalArgumentException     = "org.apache.hadoop.HadoopIllegalArgumentException"
	notALeaderException          = "org.apache.hadoop.ipc.NotALeaderException"
	javaIOException              = "java.io.IOException"
)

// Error represents a remote java exception from an HDFS namenode or datanode.
type Error interface {
	// Method returns the RPC method that encountered an error.
	Method() string
	// Desc returns the long form of the error code (for example ERROR_CHECKSUM).
	Desc() string
	// Exception returns the java exception class name (for example
	// java.io.FileNotFoundException).
	Exception() string
	// Message returns the full error message, complete with java exception
	// traceback.
	Message() string
}

func interpretCreateException(err error) error {
	if remoteErr, ok := err.(Error); ok && remoteErr.Exception() == alreadyBeingCreatedException {
		return os.ErrExist
	}

	return interpretException(err)
}

func interpretException(err error) error {
	var exception string
	if remoteErr, ok := err.(Error); ok {
		exception = remoteErr.Exception()
	}

	switch exception {
	case fileNotFoundException:
		return os.ErrNotExist
	case permissionDeniedException:
		return os.ErrPermission
	case pathIsNotEmptyDirException:
		return syscall.ENOTEMPTY
	case fileAlreadyExistsException:
		return os.ErrExist
	case invalidPathException:
		return syscall.ENOENT
	case safeModeException:
		return syscall.EROFS
	case dSQuotaExceededException:
		return syscall.EDQUOT
	case nSQuotaExceededException:
		return syscall.EDQUOT
	case parentNotDirectoryException:
		return syscall.ENOENT
	case UnresolvedLinkException:
		return syscall.ENOLINK
	case NotReplicatedYetException:
		return syscall.EPROTO // Protocol Error
	case illegalArgumentException:
		return os.ErrInvalid
	case javaIOException:
		// In HopsFS all RuntimeExceptions are
		// caught by the TX request handler and
		// then thrown as IOException.
		return os.ErrInvalid
	case notALeaderException:
		return syscall.EAGAIN
	default:
		return err
	}
}
