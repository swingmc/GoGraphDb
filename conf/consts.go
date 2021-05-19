package conf

import "runtime"

var ProjectRootPath = runtime.GOROOT() + "/src/GoGraphDb/"

const (
	InfoLogPath = "./log_file/info_log.txt"
	WarnLogPath = "./log_file/warn_log.txt"
	ErrorLogPath = "./log_file/error_log.txt"
	UndoLogPath = "./log_file/undo_log.txt"
)

const INT_MAX = int(^uint(0) >> 1)
const INT_MIN = ^INT_MAX

const (
	DataDircPath = "./Data/data.txt"
	DataSchemaPath = "./Data/meta_info.txt"
)

const (
	Modify_Nochange = 0
	Modify_Create	= 1
	Modify_Changed  = 2
	Modify_Removed  = 3
)

const Splitor = "::::"

const GcTransactionNum = 1

const (
	COMMAND_CREATE_VERTEX = "CREATE_VERTEX"
	COMMAND_CHANGE_VERTEX = "CHANGE_VERTEX"
	COMMAND_REMOVE_VERTEX = "REMOVE_VERTEX"

	COMMAND_CREATE_EDGE = "CREATE_EDGE"
	COMMAND_CHANGE_EDGE = "CHANGE_EDGE"
	COMMAND_REMOVE_EDGE = "REMOVE_EDGE"
)

const (
	InstructionPointer_NanoTime = 0
	InstructionPointer_ReadableTime = 1
	InstructionPointer_Command = 2
	InstructionPointer_Identifier = 3
	InstructionPointer_JsonObj = 4
)

const (
	InterpreterCommand_StartTransaction = "BEGIN"
	InterpreterCommand_StartReadOnlyTransaction = "BEGIN_READ"
	InterpreterCommand_EndTransaction = "END"

	InterpreterStatus_Empty = 0
	InterpreterStatus_InTransaction = 1
)

const (
	TransactionStatus_Executing = 0;
	TransactionStatus_Complete = 1;
	TransactionStatus_Canceled = 2;
)

const (
	DataReadableStatus_VersionTooLate = 0
	DataReadableStatus_Readable = 1
	DataReadableStatus_Canceled = 2
	DataReadableStatus_Executing = 3
	DataReadableStatus_OneTransaction = 4
)

const (
	DataWriteableStatus_VersionTooLate = 0
	DataWriteableStatus_Writeable = 1
	DataWriteableStatus_Executing = 2
	DataWriteableStatus_OneTransaction = 3
)

const (
	SchemaFile_VertexFlag = "Vertex"
	SchemaFile_EdgeFlag = "Edge"
)

const BaseRetryTime = 500;