package main

// This grammar is generated using the Grammar.export_go() method and
// should be used with the goleri module.
//
// Source class: SiriGrammar
// Created at: 2020-09-15 15:01:47

import (
	"regexp"

	"github.com/transceptor-technology/goleri"
)

// Element indentifiers
const (
	NoGid                   = iota
	GidAccessExpr           = iota
	GidAccessKeywords       = iota
	GidAfterExpr            = iota
	GidAggregateFunctions   = iota
	GidAlterDatabase        = iota
	GidAlterGroup           = iota
	GidAlterSeries          = iota
	GidAlterServer          = iota
	GidAlterServers         = iota
	GidAlterStmt            = iota
	GidAlterUser            = iota
	GidBeforeExpr           = iota
	GidBetweenExpr          = iota
	GidBoolOperator         = iota
	GidBoolean              = iota
	GidCDifference          = iota
	GidCalcStmt             = iota
	GidCountGroups          = iota
	GidCountPools           = iota
	GidCountSeries          = iota
	GidCountSeriesLength    = iota
	GidCountServers         = iota
	GidCountServersReceived = iota
	GidCountServersSelected = iota
	GidCountShards          = iota
	GidCountShardsSize      = iota
	GidCountStmt            = iota
	GidCountTags            = iota
	GidCountUsers           = iota
	GidCreateGroup          = iota
	GidCreateStmt           = iota
	GidCreateUser           = iota
	GidDropGroup            = iota
	GidDropSeries           = iota
	GidDropServer           = iota
	GidDropShards           = iota
	GidDropStmt             = iota
	GidDropTag              = iota
	GidDropUser             = iota
	GidFAll                 = iota
	GidFCount               = iota
	GidFDerivative          = iota
	GidFDifference          = iota
	GidFFilter              = iota
	GidFFirst               = iota
	GidFInterval            = iota
	GidFLast                = iota
	GidFLimit               = iota
	GidFMax                 = iota
	GidFMean                = iota
	GidFMedian              = iota
	GidFMedianHigh          = iota
	GidFMedianLow           = iota
	GidFMin                 = iota
	GidFPoints              = iota
	GidFPvariance           = iota
	GidFStddev              = iota
	GidFSum                 = iota
	GidFTimeval             = iota
	GidFVariance            = iota
	GidGrantStmt            = iota
	GidGrantUser            = iota
	GidGroupColumns         = iota
	GidGroupName            = iota
	GidGroupTagMatch        = iota
	GidHelpAccess           = iota
	GidHelpAlter            = iota
	GidHelpAlterDatabase    = iota
	GidHelpAlterGroup       = iota
	GidHelpAlterServer      = iota
	GidHelpAlterServers     = iota
	GidHelpAlterUser        = iota
	GidHelpCount            = iota
	GidHelpCountGroups      = iota
	GidHelpCountPools       = iota
	GidHelpCountSeries      = iota
	GidHelpCountServers     = iota
	GidHelpCountShards      = iota
	GidHelpCountUsers       = iota
	GidHelpCreate           = iota
	GidHelpCreateGroup      = iota
	GidHelpCreateUser       = iota
	GidHelpDrop             = iota
	GidHelpDropGroup        = iota
	GidHelpDropSeries       = iota
	GidHelpDropServer       = iota
	GidHelpDropShards       = iota
	GidHelpDropUser         = iota
	GidHelpFunctions        = iota
	GidHelpGrant            = iota
	GidHelpList             = iota
	GidHelpListGroups       = iota
	GidHelpListPools        = iota
	GidHelpListSeries       = iota
	GidHelpListServers      = iota
	GidHelpListShards       = iota
	GidHelpListUsers        = iota
	GidHelpNoaccess         = iota
	GidHelpRevoke           = iota
	GidHelpSelect           = iota
	GidHelpShow             = iota
	GidHelpStmt             = iota
	GidHelpTimeit           = iota
	GidHelpTimezones        = iota
	GidIntExpr              = iota
	GidIntOperator          = iota
	GidKAccess              = iota
	GidKActiveHandles       = iota
	GidKActiveTasks         = iota
	GidKAddress             = iota
	GidKAfter               = iota
	GidKAll                 = iota
	GidKAlter               = iota
	GidKAnd                 = iota
	GidKAs                  = iota
	GidKBackupMode          = iota
	GidKBefore              = iota
	GidKBetween             = iota
	GidKBufferPath          = iota
	GidKBufferSize          = iota
	GidKCount               = iota
	GidKCreate              = iota
	GidKCritical            = iota
	GidKDatabase            = iota
	GidKDbname              = iota
	GidKDbpath              = iota
	GidKDebug               = iota
	GidKDerivative          = iota
	GidKDifference          = iota
	GidKDrop                = iota
	GidKDropThreshold       = iota
	GidKDurationLog         = iota
	GidKDurationNum         = iota
	GidKEnd                 = iota
	GidKError               = iota
	GidKExpirationLog       = iota
	GidKExpirationNum       = iota
	GidKExpression          = iota
	GidKFalse               = iota
	GidKFifoFiles           = iota
	GidKFilter              = iota
	GidKFirst               = iota
	GidKFloat               = iota
	GidKFor                 = iota
	GidKFrom                = iota
	GidKFull                = iota
	GidKGrant               = iota
	GidKGroup               = iota
	GidKGroups              = iota
	GidKHelp                = iota
	GidKIdlePercentage      = iota
	GidKIdleTime            = iota
	GidKIgnoreThreshold     = iota
	GidKInf                 = iota
	GidKInfo                = iota
	GidKInsert              = iota
	GidKInteger             = iota
	GidKIntersection        = iota
	GidKInterval            = iota
	GidKIpSupport           = iota
	GidKLast                = iota
	GidKLength              = iota
	GidKLibuv               = iota
	GidKLimit               = iota
	GidKList                = iota
	GidKListLimit           = iota
	GidKLog                 = iota
	GidKLogLevel            = iota
	GidKMax                 = iota
	GidKMaxOpenFiles        = iota
	GidKMean                = iota
	GidKMedian              = iota
	GidKMedianHigh          = iota
	GidKMedianLow           = iota
	GidKMemUsage            = iota
	GidKMerge               = iota
	GidKMin                 = iota
	GidKModify              = iota
	GidKName                = iota
	GidKNan                 = iota
	GidKNinf                = iota
	GidKNow                 = iota
	GidKNumber              = iota
	GidKOnline              = iota
	GidKOpenFiles           = iota
	GidKOr                  = iota
	GidKPassword            = iota
	GidKPoints              = iota
	GidKPool                = iota
	GidKPools               = iota
	GidKPort                = iota
	GidKPrefix              = iota
	GidKPvariance           = iota
	GidKRead                = iota
	GidKReceivedPoints      = iota
	GidKReindexProgress     = iota
	GidKRevoke              = iota
	GidKSelect              = iota
	GidKSelectPointsLimit   = iota
	GidKSelectedPoints      = iota
	GidKSeries              = iota
	GidKServer              = iota
	GidKServers             = iota
	GidKSet                 = iota
	GidKShards              = iota
	GidKShow                = iota
	GidKSid                 = iota
	GidKSize                = iota
	GidKStart               = iota
	GidKStartupTime         = iota
	GidKStatus              = iota
	GidKStddev              = iota
	GidKString              = iota
	GidKSuffix              = iota
	GidKSum                 = iota
	GidKSymmetricDifference = iota
	GidKSyncProgress        = iota
	GidKTag                 = iota
	GidKTags                = iota
	GidKTeePipeName         = iota
	GidKTimePrecision       = iota
	GidKTimeit              = iota
	GidKTimeval             = iota
	GidKTimezone            = iota
	GidKTo                  = iota
	GidKTrue                = iota
	GidKType                = iota
	GidKUnion               = iota
	GidKUntag               = iota
	GidKUptime              = iota
	GidKUser                = iota
	GidKUsers               = iota
	GidKUsing               = iota
	GidKUuid                = iota
	GidKVariance            = iota
	GidKVersion             = iota
	GidKWarning             = iota
	GidKWhere               = iota
	GidKWhoAmI              = iota
	GidKWrite               = iota
	GidLimitExpr            = iota
	GidListGroups           = iota
	GidListPools            = iota
	GidListSeries           = iota
	GidListServers          = iota
	GidListShards           = iota
	GidListStmt             = iota
	GidListTags             = iota
	GidListUsers            = iota
	GidLogKeywords          = iota
	GidMergeAs              = iota
	GidPoolColumns          = iota
	GidPoolProps            = iota
	GidPrefixExpr           = iota
	GidRComment             = iota
	GidRDoubleqStr          = iota
	GidRFloat               = iota
	GidRGraveStr            = iota
	GidRInteger             = iota
	GidRRegex               = iota
	GidRSingleqStr          = iota
	GidRTimeStr             = iota
	GidRUinteger            = iota
	GidRUuidStr             = iota
	GidRevokeStmt           = iota
	GidRevokeUser           = iota
	GidSTART                = iota
	GidSelectAggregate      = iota
	GidSelectAggregates     = iota
	GidSelectStmt           = iota
	GidSeriesAll            = iota
	GidSeriesColumns        = iota
	GidSeriesMatch          = iota
	GidSeriesName           = iota
	GidSeriesParentheses    = iota
	GidSeriesRe             = iota
	GidSeriesSetopr         = iota
	GidServerColumns        = iota
	GidSetAddress           = iota
	GidSetBackupMode        = iota
	GidSetDropThreshold     = iota
	GidSetExpirationLog     = iota
	GidSetExpirationNum     = iota
	GidSetExpression        = iota
	GidSetIgnoreThreshold   = iota
	GidSetListLimit         = iota
	GidSetLogLevel          = iota
	GidSetName              = iota
	GidSetPassword          = iota
	GidSetPort              = iota
	GidSetSelectPointsLimit = iota
	GidSetTeePipeName       = iota
	GidSetTimezone          = iota
	GidShardColumns         = iota
	GidShowStmt             = iota
	GidStrOperator          = iota
	GidString               = iota
	GidSuffixExpr           = iota
	GidTagColumns           = iota
	GidTagName              = iota
	GidTagSeries            = iota
	GidTimeExpr             = iota
	GidTimeitStmt           = iota
	GidUntagSeries          = iota
	GidUserColumns          = iota
	GidUuid                 = iota
	GidWhereGroup           = iota
	GidWherePool            = iota
	GidWhereSeries          = iota
	GidWhereServer          = iota
	GidWhereShard           = iota
	GidWhereTag             = iota
	GidWhereUser            = iota
)

// SiriGrammar returns a compiled goleri grammar.
func SiriGrammar() *goleri.Grammar {
	rFloat := goleri.NewRegex(GidRFloat, regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+`))
	rInteger := goleri.NewRegex(GidRInteger, regexp.MustCompile(`^[-+]?[0-9]+`))
	rUinteger := goleri.NewRegex(GidRUinteger, regexp.MustCompile(`^[0-9]+`))
	rTimeStr := goleri.NewRegex(GidRTimeStr, regexp.MustCompile(`^[0-9]+[smhdw]`))
	rSingleqStr := goleri.NewRegex(GidRSingleqStr, regexp.MustCompile(`^(?:'(?:[^']*)')+`))
	rDoubleqStr := goleri.NewRegex(GidRDoubleqStr, regexp.MustCompile(`^(?:"(?:[^"]*)")+`))
	rGraveStr := goleri.NewRegex(GidRGraveStr, regexp.MustCompile(`^(?:`+"`"+`(?:[^`+"`"+`]*)`+"`"+`)+`))
	rUuidStr := goleri.NewRegex(GidRUuidStr, regexp.MustCompile(`^[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}`))
	rRegex := goleri.NewRegex(GidRRegex, regexp.MustCompile(`^(/[^/\\]*(?:\\.[^/\\]*)*/i?)`))
	rComment := goleri.NewRegex(GidRComment, regexp.MustCompile(`^#.*`))
	kAccess := goleri.NewKeyword(GidKAccess, "access", false)
	kActiveHandles := goleri.NewKeyword(GidKActiveHandles, "active_handles", false)
	kActiveTasks := goleri.NewKeyword(GidKActiveTasks, "active_tasks", false)
	kAddress := goleri.NewKeyword(GidKAddress, "address", false)
	kAfter := goleri.NewKeyword(GidKAfter, "after", false)
	kAll := goleri.NewKeyword(GidKAll, "all", false)
	kAlter := goleri.NewKeyword(GidKAlter, "alter", false)
	kAnd := goleri.NewKeyword(GidKAnd, "and", false)
	kAs := goleri.NewKeyword(GidKAs, "as", false)
	kBackupMode := goleri.NewKeyword(GidKBackupMode, "backup_mode", false)
	kBefore := goleri.NewKeyword(GidKBefore, "before", false)
	kBufferSize := goleri.NewKeyword(GidKBufferSize, "buffer_size", false)
	kBufferPath := goleri.NewKeyword(GidKBufferPath, "buffer_path", false)
	kBetween := goleri.NewKeyword(GidKBetween, "between", false)
	kCount := goleri.NewKeyword(GidKCount, "count", false)
	kCreate := goleri.NewKeyword(GidKCreate, "create", false)
	kCritical := goleri.NewKeyword(GidKCritical, "critical", false)
	kDatabase := goleri.NewKeyword(GidKDatabase, "database", false)
	kDbname := goleri.NewKeyword(GidKDbname, "dbname", false)
	kDbpath := goleri.NewKeyword(GidKDbpath, "dbpath", false)
	kDebug := goleri.NewKeyword(GidKDebug, "debug", false)
	kDerivative := goleri.NewKeyword(GidKDerivative, "derivative", false)
	kDifference := goleri.NewKeyword(GidKDifference, "difference", false)
	kDrop := goleri.NewKeyword(GidKDrop, "drop", false)
	kDropThreshold := goleri.NewKeyword(GidKDropThreshold, "drop_threshold", false)
	kDurationLog := goleri.NewKeyword(GidKDurationLog, "duration_log", false)
	kDurationNum := goleri.NewKeyword(GidKDurationNum, "duration_num", false)
	kEnd := goleri.NewKeyword(GidKEnd, "end", false)
	kError := goleri.NewKeyword(GidKError, "error", false)
	kExpression := goleri.NewKeyword(GidKExpression, "expression", false)
	kFalse := goleri.NewKeyword(GidKFalse, "false", false)
	kFifoFiles := goleri.NewKeyword(GidKFifoFiles, "fifo_files", false)
	kFilter := goleri.NewKeyword(GidKFilter, "filter", false)
	kFirst := goleri.NewKeyword(GidKFirst, "first", false)
	kFloat := goleri.NewKeyword(GidKFloat, "float", false)
	kFor := goleri.NewKeyword(GidKFor, "for", false)
	kFrom := goleri.NewKeyword(GidKFrom, "from", false)
	kFull := goleri.NewKeyword(GidKFull, "full", false)
	kGrant := goleri.NewKeyword(GidKGrant, "grant", false)
	kGroup := goleri.NewKeyword(GidKGroup, "group", false)
	kGroups := goleri.NewKeyword(GidKGroups, "groups", false)
	kHelp := goleri.NewChoice(
		GidKHelp,
		true,
		goleri.NewKeyword(NoGid, "help", false),
		goleri.NewToken(NoGid, "?"),
	)
	kIdlePercentage := goleri.NewKeyword(GidKIdlePercentage, "idle_percentage", false)
	kIdleTime := goleri.NewKeyword(GidKIdleTime, "idle_time", false)
	kInf := goleri.NewKeyword(GidKInf, "inf", false)
	kInfo := goleri.NewKeyword(GidKInfo, "info", false)
	kIgnoreThreshold := goleri.NewKeyword(GidKIgnoreThreshold, "ignore_threshold", false)
	kInsert := goleri.NewKeyword(GidKInsert, "insert", false)
	kInteger := goleri.NewKeyword(GidKInteger, "integer", false)
	kIntersection := goleri.NewChoice(
		GidKIntersection,
		false,
		goleri.NewToken(NoGid, "&"),
		goleri.NewKeyword(NoGid, "intersection", false),
	)
	kInterval := goleri.NewKeyword(GidKInterval, "interval", false)
	kIpSupport := goleri.NewKeyword(GidKIpSupport, "ip_support", false)
	kLast := goleri.NewKeyword(GidKLast, "last", false)
	kLength := goleri.NewKeyword(GidKLength, "length", false)
	kLibuv := goleri.NewKeyword(GidKLibuv, "libuv", false)
	kLimit := goleri.NewKeyword(GidKLimit, "limit", false)
	kList := goleri.NewKeyword(GidKList, "list", false)
	kListLimit := goleri.NewKeyword(GidKListLimit, "list_limit", false)
	kLog := goleri.NewKeyword(GidKLog, "log", false)
	kLogLevel := goleri.NewKeyword(GidKLogLevel, "log_level", false)
	kMax := goleri.NewKeyword(GidKMax, "max", false)
	kMaxOpenFiles := goleri.NewKeyword(GidKMaxOpenFiles, "max_open_files", false)
	kMean := goleri.NewKeyword(GidKMean, "mean", false)
	kMedian := goleri.NewKeyword(GidKMedian, "median", false)
	kMedianHigh := goleri.NewKeyword(GidKMedianHigh, "median_high", false)
	kMedianLow := goleri.NewKeyword(GidKMedianLow, "median_low", false)
	kMemUsage := goleri.NewKeyword(GidKMemUsage, "mem_usage", false)
	kMerge := goleri.NewKeyword(GidKMerge, "merge", false)
	kMin := goleri.NewKeyword(GidKMin, "min", false)
	kModify := goleri.NewKeyword(GidKModify, "modify", false)
	kName := goleri.NewKeyword(GidKName, "name", false)
	kNan := goleri.NewKeyword(GidKNan, "nan", false)
	kNinf := goleri.NewSequence(
		GidKNinf,
		goleri.NewToken(NoGid, "-"),
		kInf,
	)
	kNow := goleri.NewKeyword(GidKNow, "now", false)
	kNumber := goleri.NewKeyword(GidKNumber, "number", false)
	kOnline := goleri.NewKeyword(GidKOnline, "online", false)
	kOpenFiles := goleri.NewKeyword(GidKOpenFiles, "open_files", false)
	kOr := goleri.NewKeyword(GidKOr, "or", false)
	kPassword := goleri.NewKeyword(GidKPassword, "password", false)
	kPoints := goleri.NewKeyword(GidKPoints, "points", false)
	kPool := goleri.NewKeyword(GidKPool, "pool", false)
	kPools := goleri.NewKeyword(GidKPools, "pools", false)
	kPort := goleri.NewKeyword(GidKPort, "port", false)
	kPrefix := goleri.NewKeyword(GidKPrefix, "prefix", false)
	kPvariance := goleri.NewKeyword(GidKPvariance, "pvariance", false)
	kRead := goleri.NewKeyword(GidKRead, "read", false)
	kReceivedPoints := goleri.NewKeyword(GidKReceivedPoints, "received_points", false)
	kReindexProgress := goleri.NewKeyword(GidKReindexProgress, "reindex_progress", false)
	kRevoke := goleri.NewKeyword(GidKRevoke, "revoke", false)
	kSelect := goleri.NewKeyword(GidKSelect, "select", false)
	kSelectPointsLimit := goleri.NewKeyword(GidKSelectPointsLimit, "select_points_limit", false)
	kSelectedPoints := goleri.NewKeyword(GidKSelectedPoints, "selected_points", false)
	kSeries := goleri.NewKeyword(GidKSeries, "series", false)
	kServer := goleri.NewKeyword(GidKServer, "server", false)
	kServers := goleri.NewKeyword(GidKServers, "servers", false)
	kSet := goleri.NewKeyword(GidKSet, "set", false)
	kExpirationLog := goleri.NewKeyword(GidKExpirationLog, "expiration_log", false)
	kExpirationNum := goleri.NewKeyword(GidKExpirationNum, "expiration_num", false)
	kShards := goleri.NewKeyword(GidKShards, "shards", false)
	kShow := goleri.NewKeyword(GidKShow, "show", false)
	kSid := goleri.NewKeyword(GidKSid, "sid", false)
	kSize := goleri.NewKeyword(GidKSize, "size", false)
	kStart := goleri.NewKeyword(GidKStart, "start", false)
	kStartupTime := goleri.NewKeyword(GidKStartupTime, "startup_time", false)
	kStatus := goleri.NewKeyword(GidKStatus, "status", false)
	kStddev := goleri.NewKeyword(GidKStddev, "stddev", false)
	kString := goleri.NewKeyword(GidKString, "string", false)
	kSuffix := goleri.NewKeyword(GidKSuffix, "suffix", false)
	kSum := goleri.NewKeyword(GidKSum, "sum", false)
	kSymmetricDifference := goleri.NewChoice(
		GidKSymmetricDifference,
		false,
		goleri.NewToken(NoGid, "^"),
		goleri.NewKeyword(NoGid, "symmetric_difference", false),
	)
	kSyncProgress := goleri.NewKeyword(GidKSyncProgress, "sync_progress", false)
	kTag := goleri.NewKeyword(GidKTag, "tag", false)
	kTags := goleri.NewKeyword(GidKTags, "tags", false)
	kTeePipeName := goleri.NewKeyword(GidKTeePipeName, "tee_pipe_name", false)
	kTimePrecision := goleri.NewKeyword(GidKTimePrecision, "time_precision", false)
	kTimeit := goleri.NewKeyword(GidKTimeit, "timeit", false)
	kTimeval := goleri.NewKeyword(GidKTimeval, "timeval", false)
	kTimezone := goleri.NewKeyword(GidKTimezone, "timezone", false)
	kTo := goleri.NewKeyword(GidKTo, "to", false)
	kTrue := goleri.NewKeyword(GidKTrue, "true", false)
	kType := goleri.NewKeyword(GidKType, "type", false)
	kUnion := goleri.NewChoice(
		GidKUnion,
		false,
		goleri.NewTokens(NoGid, ", |"),
		goleri.NewKeyword(NoGid, "union", false),
	)
	kUntag := goleri.NewKeyword(GidKUntag, "untag", false)
	kUptime := goleri.NewKeyword(GidKUptime, "uptime", false)
	kUser := goleri.NewKeyword(GidKUser, "user", false)
	kUsers := goleri.NewKeyword(GidKUsers, "users", false)
	kUsing := goleri.NewKeyword(GidKUsing, "using", false)
	kUuid := goleri.NewKeyword(GidKUuid, "uuid", false)
	kVariance := goleri.NewKeyword(GidKVariance, "variance", false)
	kVersion := goleri.NewKeyword(GidKVersion, "version", false)
	kWarning := goleri.NewKeyword(GidKWarning, "warning", false)
	kWhere := goleri.NewKeyword(GidKWhere, "where", false)
	kWhoAmI := goleri.NewKeyword(GidKWhoAmI, "who_am_i", false)
	kWrite := goleri.NewKeyword(GidKWrite, "write", false)
	cDifference := goleri.NewChoice(
		GidCDifference,
		false,
		goleri.NewToken(NoGid, "-"),
		kDifference,
	)
	accessKeywords := goleri.NewChoice(
		GidAccessKeywords,
		false,
		kRead,
		kWrite,
		kModify,
		kFull,
		kSelect,
		kShow,
		kList,
		kCount,
		kCreate,
		kInsert,
		kDrop,
		kGrant,
		kRevoke,
		kAlter,
	)
	Boolean := goleri.NewChoice(
		GidBoolean,
		false,
		kTrue,
		kFalse,
	)
	logKeywords := goleri.NewChoice(
		GidLogKeywords,
		false,
		kDebug,
		kInfo,
		kWarning,
		kError,
		kCritical,
	)
	intExpr := goleri.NewPrio(
		GidIntExpr,
		rInteger,
		goleri.NewSequence(
			NoGid,
			goleri.NewToken(NoGid, "("),
			goleri.THIS,
			goleri.NewToken(NoGid, ")"),
		),
		goleri.NewSequence(
			NoGid,
			goleri.THIS,
			goleri.NewTokens(NoGid, "+ - * % /"),
			goleri.THIS,
		),
	)
	string := goleri.NewChoice(
		GidString,
		false,
		rSingleqStr,
		rDoubleqStr,
	)
	timeExpr := goleri.NewPrio(
		GidTimeExpr,
		rTimeStr,
		kNow,
		string,
		rInteger,
		goleri.NewSequence(
			NoGid,
			goleri.NewToken(NoGid, "("),
			goleri.THIS,
			goleri.NewToken(NoGid, ")"),
		),
		goleri.NewSequence(
			NoGid,
			goleri.THIS,
			goleri.NewTokens(NoGid, "+ - * % /"),
			goleri.THIS,
		),
	)
	seriesColumns := goleri.NewList(GidSeriesColumns, goleri.NewChoice(
		NoGid,
		false,
		kName,
		kType,
		kLength,
		kStart,
		kEnd,
		kPool,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	shardColumns := goleri.NewList(GidShardColumns, goleri.NewChoice(
		NoGid,
		false,
		kSid,
		kPool,
		kServer,
		kSize,
		kStart,
		kEnd,
		kType,
		kStatus,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	serverColumns := goleri.NewList(GidServerColumns, goleri.NewChoice(
		NoGid,
		false,
		kAddress,
		kBufferPath,
		kBufferSize,
		kDbpath,
		kIpSupport,
		kLibuv,
		kName,
		kPort,
		kUuid,
		kPool,
		kVersion,
		kOnline,
		kStartupTime,
		kStatus,
		kActiveHandles,
		kActiveTasks,
		kFifoFiles,
		kIdlePercentage,
		kIdleTime,
		kLogLevel,
		kMaxOpenFiles,
		kMemUsage,
		kOpenFiles,
		kReceivedPoints,
		kReindexProgress,
		kSelectedPoints,
		kSyncProgress,
		kTeePipeName,
		kUptime,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	groupColumns := goleri.NewList(GidGroupColumns, goleri.NewChoice(
		NoGid,
		false,
		kExpression,
		kName,
		kSeries,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	userColumns := goleri.NewList(GidUserColumns, goleri.NewChoice(
		NoGid,
		false,
		kName,
		kAccess,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	tagColumns := goleri.NewList(GidTagColumns, goleri.NewChoice(
		NoGid,
		false,
		kName,
		kSeries,
	), goleri.NewToken(NoGid, ","), 1, 0, false)
	poolProps := goleri.NewChoice(
		GidPoolProps,
		false,
		kPool,
		kServers,
		kSeries,
	)
	poolColumns := goleri.NewList(GidPoolColumns, poolProps, goleri.NewToken(NoGid, ","), 1, 0, false)
	boolOperator := goleri.NewTokens(GidBoolOperator, "== !=")
	intOperator := goleri.NewTokens(GidIntOperator, "== != <= >= < >")
	strOperator := goleri.NewTokens(GidStrOperator, "== != <= >= !~ < > ~")
	whereGroup := goleri.NewSequence(
		GidWhereGroup,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				kSeries,
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kExpression,
					kName,
				),
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	whereTag := goleri.NewSequence(
		GidWhereTag,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				kName,
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				kSeries,
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	wherePool := goleri.NewSequence(
		GidWherePool,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				poolProps,
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	whereSeries := goleri.NewSequence(
		GidWhereSeries,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kLength,
					kPool,
				),
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				kName,
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kStart,
					kEnd,
				),
				intOperator,
				timeExpr,
			),
			goleri.NewSequence(
				NoGid,
				kType,
				boolOperator,
				goleri.NewChoice(
					NoGid,
					false,
					kString,
					kInteger,
					kFloat,
				),
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	whereServer := goleri.NewSequence(
		GidWhereServer,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kActiveHandles,
					kActiveTasks,
					kBufferSize,
					kFifoFiles,
					kIdlePercentage,
					kIdleTime,
					kPort,
					kPool,
					kStartupTime,
					kMaxOpenFiles,
					kMemUsage,
					kOpenFiles,
					kReceivedPoints,
					kSelectedPoints,
					kUptime,
				),
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kAddress,
					kBufferPath,
					kDbpath,
					kIpSupport,
					kLibuv,
					kName,
					kUuid,
					kVersion,
					kStatus,
					kReindexProgress,
					kSyncProgress,
					kTeePipeName,
				),
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				kOnline,
				boolOperator,
				Boolean,
			),
			goleri.NewSequence(
				NoGid,
				kLogLevel,
				intOperator,
				logKeywords,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	whereShard := goleri.NewSequence(
		GidWhereShard,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kSid,
					kPool,
					kSize,
				),
				intOperator,
				intExpr,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					true,
					kServer,
					kStatus,
				),
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewChoice(
					NoGid,
					false,
					kStart,
					kEnd,
				),
				intOperator,
				timeExpr,
			),
			goleri.NewSequence(
				NoGid,
				kType,
				boolOperator,
				goleri.NewChoice(
					NoGid,
					false,
					kNumber,
					kLog,
				),
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	whereUser := goleri.NewSequence(
		GidWhereUser,
		kWhere,
		goleri.NewPrio(
			NoGid,
			goleri.NewSequence(
				NoGid,
				kName,
				strOperator,
				string,
			),
			goleri.NewSequence(
				NoGid,
				kAccess,
				intOperator,
				accessKeywords,
			),
			goleri.NewSequence(
				NoGid,
				goleri.NewToken(NoGid, "("),
				goleri.THIS,
				goleri.NewToken(NoGid, ")"),
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kAnd,
				goleri.THIS,
			),
			goleri.NewSequence(
				NoGid,
				goleri.THIS,
				kOr,
				goleri.THIS,
			),
		),
	)
	seriesSetopr := goleri.NewChoice(
		GidSeriesSetopr,
		false,
		kUnion,
		cDifference,
		kIntersection,
		kSymmetricDifference,
	)
	seriesParentheses := goleri.NewSequence(
		GidSeriesParentheses,
		goleri.NewToken(NoGid, "("),
		goleri.THIS,
		goleri.NewToken(NoGid, ")"),
	)
	seriesAll := goleri.NewChoice(
		GidSeriesAll,
		false,
		goleri.NewToken(NoGid, "*"),
		kAll,
	)
	seriesName := goleri.NewRepeat(GidSeriesName, string, 1, 1)
	groupName := goleri.NewRepeat(GidGroupName, rGraveStr, 1, 1)
	tagName := goleri.NewRepeat(GidTagName, rGraveStr, 1, 1)
	seriesRe := goleri.NewRepeat(GidSeriesRe, rRegex, 1, 1)
	uuid := goleri.NewChoice(
		GidUuid,
		false,
		rUuidStr,
		string,
	)
	groupTagMatch := goleri.NewRepeat(GidGroupTagMatch, rGraveStr, 1, 1)
	seriesMatch := goleri.NewPrio(
		GidSeriesMatch,
		goleri.NewList(NoGid, goleri.NewChoice(
			NoGid,
			false,
			seriesAll,
			seriesName,
			groupTagMatch,
			seriesRe,
		), seriesSetopr, 1, 0, false),
		goleri.NewChoice(
			NoGid,
			false,
			seriesAll,
			seriesName,
			groupTagMatch,
			seriesRe,
		),
		seriesParentheses,
		goleri.NewSequence(
			NoGid,
			goleri.THIS,
			seriesSetopr,
			goleri.THIS,
		),
	)
	limitExpr := goleri.NewSequence(
		GidLimitExpr,
		kLimit,
		intExpr,
	)
	beforeExpr := goleri.NewSequence(
		GidBeforeExpr,
		kBefore,
		timeExpr,
	)
	afterExpr := goleri.NewSequence(
		GidAfterExpr,
		kAfter,
		timeExpr,
	)
	betweenExpr := goleri.NewSequence(
		GidBetweenExpr,
		kBetween,
		timeExpr,
		kAnd,
		timeExpr,
	)
	accessExpr := goleri.NewList(GidAccessExpr, accessKeywords, goleri.NewToken(NoGid, ","), 1, 0, false)
	prefixExpr := goleri.NewSequence(
		GidPrefixExpr,
		kPrefix,
		string,
	)
	suffixExpr := goleri.NewSequence(
		GidSuffixExpr,
		kSuffix,
		string,
	)
	fAll := goleri.NewChoice(
		GidFAll,
		false,
		goleri.NewToken(NoGid, "*"),
		kAll,
	)
	fPoints := goleri.NewRepeat(GidFPoints, kPoints, 1, 1)
	fDifference := goleri.NewSequence(
		GidFDifference,
		kDifference,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fDerivative := goleri.NewSequence(
		GidFDerivative,
		kDerivative,
		goleri.NewToken(NoGid, "("),
		goleri.NewList(NoGid, timeExpr, goleri.NewToken(NoGid, ","), 0, 2, false),
		goleri.NewToken(NoGid, ")"),
	)
	fMean := goleri.NewSequence(
		GidFMean,
		kMean,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fMedian := goleri.NewSequence(
		GidFMedian,
		kMedian,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fMedianLow := goleri.NewSequence(
		GidFMedianLow,
		kMedianLow,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fMedianHigh := goleri.NewSequence(
		GidFMedianHigh,
		kMedianHigh,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fSum := goleri.NewSequence(
		GidFSum,
		kSum,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fMin := goleri.NewSequence(
		GidFMin,
		kMin,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fMax := goleri.NewSequence(
		GidFMax,
		kMax,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fCount := goleri.NewSequence(
		GidFCount,
		kCount,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fVariance := goleri.NewSequence(
		GidFVariance,
		kVariance,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fPvariance := goleri.NewSequence(
		GidFPvariance,
		kPvariance,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fStddev := goleri.NewSequence(
		GidFStddev,
		kStddev,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fFirst := goleri.NewSequence(
		GidFFirst,
		kFirst,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fLast := goleri.NewSequence(
		GidFLast,
		kLast,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, timeExpr),
		goleri.NewToken(NoGid, ")"),
	)
	fTimeval := goleri.NewSequence(
		GidFTimeval,
		kTimeval,
		goleri.NewToken(NoGid, "("),
		goleri.NewToken(NoGid, ")"),
	)
	fInterval := goleri.NewSequence(
		GidFInterval,
		kInterval,
		goleri.NewToken(NoGid, "("),
		goleri.NewToken(NoGid, ")"),
	)
	fFilter := goleri.NewSequence(
		GidFFilter,
		kFilter,
		goleri.NewToken(NoGid, "("),
		goleri.NewOptional(NoGid, strOperator),
		goleri.NewChoice(
			NoGid,
			true,
			string,
			rInteger,
			rFloat,
			rRegex,
			kNan,
			kInf,
			kNinf,
		),
		goleri.NewToken(NoGid, ")"),
	)
	fLimit := goleri.NewSequence(
		GidFLimit,
		kLimit,
		goleri.NewToken(NoGid, "("),
		intExpr,
		goleri.NewToken(NoGid, ","),
		goleri.NewChoice(
			NoGid,
			false,
			kMean,
			kMedian,
			kMedianHigh,
			kMedianLow,
			kSum,
			kMin,
			kMax,
			kCount,
			kVariance,
			kPvariance,
			kStddev,
			kFirst,
			kLast,
		),
		goleri.NewToken(NoGid, ")"),
	)
	aggregateFunctions := goleri.NewList(GidAggregateFunctions, goleri.NewChoice(
		NoGid,
		false,
		fAll,
		fLimit,
		fMean,
		fSum,
		fMedian,
		fMedianLow,
		fMedianHigh,
		fMin,
		fMax,
		fCount,
		fVariance,
		fPvariance,
		fStddev,
		fFirst,
		fLast,
		fTimeval,
		fInterval,
		fDifference,
		fDerivative,
		fFilter,
		fPoints,
	), goleri.NewToken(NoGid, "=>"), 1, 0, false)
	selectAggregate := goleri.NewSequence(
		GidSelectAggregate,
		aggregateFunctions,
		goleri.NewOptional(NoGid, prefixExpr),
		goleri.NewOptional(NoGid, suffixExpr),
	)
	selectAggregates := goleri.NewList(GidSelectAggregates, selectAggregate, goleri.NewToken(NoGid, ","), 1, 0, false)
	mergeAs := goleri.NewSequence(
		GidMergeAs,
		kMerge,
		kAs,
		string,
		goleri.NewOptional(NoGid, goleri.NewSequence(
			NoGid,
			kUsing,
			aggregateFunctions,
		)),
	)
	setAddress := goleri.NewSequence(
		GidSetAddress,
		kSet,
		kAddress,
		string,
	)
	setTeePipeName := goleri.NewSequence(
		GidSetTeePipeName,
		kSet,
		kTeePipeName,
		goleri.NewChoice(
			NoGid,
			false,
			kFalse,
			string,
		),
	)
	setBackupMode := goleri.NewSequence(
		GidSetBackupMode,
		kSet,
		kBackupMode,
		Boolean,
	)
	setDropThreshold := goleri.NewSequence(
		GidSetDropThreshold,
		kSet,
		kDropThreshold,
		rFloat,
	)
	setExpression := goleri.NewSequence(
		GidSetExpression,
		kSet,
		kExpression,
		rRegex,
	)
	setIgnoreThreshold := goleri.NewSequence(
		GidSetIgnoreThreshold,
		kSet,
		kIgnoreThreshold,
		Boolean,
	)
	setListLimit := goleri.NewSequence(
		GidSetListLimit,
		kSet,
		kListLimit,
		rUinteger,
	)
	setLogLevel := goleri.NewSequence(
		GidSetLogLevel,
		kSet,
		kLogLevel,
		logKeywords,
	)
	setName := goleri.NewSequence(
		GidSetName,
		kSet,
		kName,
		string,
	)
	setPassword := goleri.NewSequence(
		GidSetPassword,
		kSet,
		kPassword,
		string,
	)
	setPort := goleri.NewSequence(
		GidSetPort,
		kSet,
		kPort,
		rUinteger,
	)
	setSelectPointsLimit := goleri.NewSequence(
		GidSetSelectPointsLimit,
		kSet,
		kSelectPointsLimit,
		rUinteger,
	)
	setTimezone := goleri.NewSequence(
		GidSetTimezone,
		kSet,
		kTimezone,
		string,
	)
	tagSeries := goleri.NewSequence(
		GidTagSeries,
		kTag,
		tagName,
	)
	untagSeries := goleri.NewSequence(
		GidUntagSeries,
		kUntag,
		tagName,
	)
	setExpirationNum := goleri.NewSequence(
		GidSetExpirationNum,
		kSet,
		kExpirationNum,
		timeExpr,
		goleri.NewOptional(NoGid, setIgnoreThreshold),
	)
	setExpirationLog := goleri.NewSequence(
		GidSetExpirationLog,
		kSet,
		kExpirationLog,
		timeExpr,
		goleri.NewOptional(NoGid, setIgnoreThreshold),
	)
	alterDatabase := goleri.NewSequence(
		GidAlterDatabase,
		kDatabase,
		goleri.NewChoice(
			NoGid,
			false,
			setDropThreshold,
			setListLimit,
			setSelectPointsLimit,
			setTimezone,
			setExpirationNum,
			setExpirationLog,
		),
	)
	alterGroup := goleri.NewSequence(
		GidAlterGroup,
		kGroup,
		groupName,
		goleri.NewChoice(
			NoGid,
			false,
			setExpression,
			setName,
		),
	)
	alterServer := goleri.NewSequence(
		GidAlterServer,
		kServer,
		uuid,
		goleri.NewChoice(
			NoGid,
			false,
			setLogLevel,
			setBackupMode,
			setTeePipeName,
			setAddress,
			setPort,
		),
	)
	alterServers := goleri.NewSequence(
		GidAlterServers,
		kServers,
		goleri.NewOptional(NoGid, whereServer),
		goleri.NewChoice(
			NoGid,
			false,
			setLogLevel,
			setTeePipeName,
		),
	)
	alterUser := goleri.NewSequence(
		GidAlterUser,
		kUser,
		string,
		goleri.NewChoice(
			NoGid,
			false,
			setPassword,
			setName,
		),
	)
	alterSeries := goleri.NewSequence(
		GidAlterSeries,
		kSeries,
		seriesMatch,
		goleri.NewOptional(NoGid, whereSeries),
		goleri.NewChoice(
			NoGid,
			false,
			tagSeries,
			untagSeries,
		),
	)
	countGroups := goleri.NewSequence(
		GidCountGroups,
		kGroups,
		goleri.NewOptional(NoGid, whereGroup),
	)
	countTags := goleri.NewSequence(
		GidCountTags,
		kTags,
		goleri.NewOptional(NoGid, whereTag),
	)
	countPools := goleri.NewSequence(
		GidCountPools,
		kPools,
		goleri.NewOptional(NoGid, wherePool),
	)
	countSeries := goleri.NewSequence(
		GidCountSeries,
		kSeries,
		goleri.NewOptional(NoGid, seriesMatch),
		goleri.NewOptional(NoGid, whereSeries),
	)
	countServers := goleri.NewSequence(
		GidCountServers,
		kServers,
		goleri.NewOptional(NoGid, whereServer),
	)
	countServersReceived := goleri.NewSequence(
		GidCountServersReceived,
		kServers,
		kReceivedPoints,
		goleri.NewOptional(NoGid, whereServer),
	)
	countServersSelected := goleri.NewSequence(
		GidCountServersSelected,
		kServers,
		kSelectedPoints,
		goleri.NewOptional(NoGid, whereServer),
	)
	countShards := goleri.NewSequence(
		GidCountShards,
		kShards,
		goleri.NewOptional(NoGid, whereShard),
	)
	countShardsSize := goleri.NewSequence(
		GidCountShardsSize,
		kShards,
		kSize,
		goleri.NewOptional(NoGid, whereShard),
	)
	countUsers := goleri.NewSequence(
		GidCountUsers,
		kUsers,
		goleri.NewOptional(NoGid, whereUser),
	)
	countSeriesLength := goleri.NewSequence(
		GidCountSeriesLength,
		kSeries,
		kLength,
		goleri.NewOptional(NoGid, seriesMatch),
		goleri.NewOptional(NoGid, whereSeries),
	)
	createGroup := goleri.NewSequence(
		GidCreateGroup,
		kGroup,
		groupName,
		kFor,
		rRegex,
	)
	createUser := goleri.NewSequence(
		GidCreateUser,
		kUser,
		string,
		setPassword,
	)
	dropGroup := goleri.NewSequence(
		GidDropGroup,
		kGroup,
		groupName,
	)
	dropTag := goleri.NewSequence(
		GidDropTag,
		kTag,
		tagName,
	)
	dropSeries := goleri.NewSequence(
		GidDropSeries,
		kSeries,
		goleri.NewOptional(NoGid, seriesMatch),
		goleri.NewOptional(NoGid, whereSeries),
		goleri.NewOptional(NoGid, setIgnoreThreshold),
	)
	dropShards := goleri.NewSequence(
		GidDropShards,
		kShards,
		goleri.NewOptional(NoGid, whereShard),
		goleri.NewOptional(NoGid, setIgnoreThreshold),
	)
	dropServer := goleri.NewSequence(
		GidDropServer,
		kServer,
		uuid,
	)
	dropUser := goleri.NewSequence(
		GidDropUser,
		kUser,
		string,
	)
	grantUser := goleri.NewSequence(
		GidGrantUser,
		kUser,
		string,
		goleri.NewOptional(NoGid, setPassword),
	)
	listGroups := goleri.NewSequence(
		GidListGroups,
		kGroups,
		goleri.NewOptional(NoGid, groupColumns),
		goleri.NewOptional(NoGid, whereGroup),
	)
	listTags := goleri.NewSequence(
		GidListTags,
		kTags,
		goleri.NewOptional(NoGid, tagColumns),
		goleri.NewOptional(NoGid, whereTag),
	)
	listPools := goleri.NewSequence(
		GidListPools,
		kPools,
		goleri.NewOptional(NoGid, poolColumns),
		goleri.NewOptional(NoGid, wherePool),
	)
	listSeries := goleri.NewSequence(
		GidListSeries,
		kSeries,
		goleri.NewOptional(NoGid, seriesColumns),
		goleri.NewOptional(NoGid, seriesMatch),
		goleri.NewOptional(NoGid, whereSeries),
	)
	listServers := goleri.NewSequence(
		GidListServers,
		kServers,
		goleri.NewOptional(NoGid, serverColumns),
		goleri.NewOptional(NoGid, whereServer),
	)
	listShards := goleri.NewSequence(
		GidListShards,
		kShards,
		goleri.NewOptional(NoGid, shardColumns),
		goleri.NewOptional(NoGid, whereShard),
	)
	listUsers := goleri.NewSequence(
		GidListUsers,
		kUsers,
		goleri.NewOptional(NoGid, userColumns),
		goleri.NewOptional(NoGid, whereUser),
	)
	revokeUser := goleri.NewSequence(
		GidRevokeUser,
		kUser,
		string,
	)
	alterStmt := goleri.NewSequence(
		GidAlterStmt,
		kAlter,
		goleri.NewChoice(
			NoGid,
			false,
			alterSeries,
			alterUser,
			alterGroup,
			alterServer,
			alterServers,
			alterDatabase,
		),
	)
	calcStmt := goleri.NewRepeat(GidCalcStmt, timeExpr, 1, 1)
	countStmt := goleri.NewSequence(
		GidCountStmt,
		kCount,
		goleri.NewChoice(
			NoGid,
			true,
			countGroups,
			countPools,
			countSeries,
			countServers,
			countServersReceived,
			countServersSelected,
			countShards,
			countShardsSize,
			countUsers,
			countTags,
			countSeriesLength,
		),
	)
	createStmt := goleri.NewSequence(
		GidCreateStmt,
		kCreate,
		goleri.NewChoice(
			NoGid,
			true,
			createGroup,
			createUser,
		),
	)
	dropStmt := goleri.NewSequence(
		GidDropStmt,
		kDrop,
		goleri.NewChoice(
			NoGid,
			false,
			dropGroup,
			dropTag,
			dropSeries,
			dropShards,
			dropServer,
			dropUser,
		),
	)
	grantStmt := goleri.NewSequence(
		GidGrantStmt,
		kGrant,
		accessExpr,
		kTo,
		goleri.NewChoice(
			NoGid,
			false,
			grantUser,
		),
	)
	listStmt := goleri.NewSequence(
		GidListStmt,
		kList,
		goleri.NewChoice(
			NoGid,
			false,
			listSeries,
			listTags,
			listUsers,
			listShards,
			listGroups,
			listServers,
			listPools,
		),
		goleri.NewOptional(NoGid, limitExpr),
	)
	revokeStmt := goleri.NewSequence(
		GidRevokeStmt,
		kRevoke,
		accessExpr,
		kFrom,
		goleri.NewChoice(
			NoGid,
			false,
			revokeUser,
		),
	)
	selectStmt := goleri.NewSequence(
		GidSelectStmt,
		kSelect,
		selectAggregates,
		kFrom,
		seriesMatch,
		goleri.NewOptional(NoGid, whereSeries),
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			false,
			afterExpr,
			betweenExpr,
			beforeExpr,
		)),
		goleri.NewOptional(NoGid, mergeAs),
	)
	showStmt := goleri.NewSequence(
		GidShowStmt,
		kShow,
		goleri.NewList(NoGid, goleri.NewChoice(
			NoGid,
			false,
			kActiveHandles,
			kActiveTasks,
			kBufferPath,
			kBufferSize,
			kDbname,
			kDbpath,
			kDropThreshold,
			kDurationLog,
			kDurationNum,
			kFifoFiles,
			kExpirationLog,
			kExpirationNum,
			kIdlePercentage,
			kIdleTime,
			kIpSupport,
			kLibuv,
			kListLimit,
			kLogLevel,
			kMaxOpenFiles,
			kMemUsage,
			kOpenFiles,
			kPool,
			kReceivedPoints,
			kReindexProgress,
			kSelectedPoints,
			kSelectPointsLimit,
			kServer,
			kStartupTime,
			kStatus,
			kSyncProgress,
			kTeePipeName,
			kTimePrecision,
			kTimezone,
			kUptime,
			kUuid,
			kVersion,
			kWhoAmI,
		), goleri.NewToken(NoGid, ","), 0, 0, false),
	)
	timeitStmt := goleri.NewRepeat(GidTimeitStmt, kTimeit, 1, 1)
	helpStmt := goleri.NewRef()
	START := goleri.NewSequence(
		GidSTART,
		goleri.NewOptional(NoGid, timeitStmt),
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			false,
			selectStmt,
			listStmt,
			countStmt,
			alterStmt,
			createStmt,
			dropStmt,
			grantStmt,
			revokeStmt,
			showStmt,
			calcStmt,
			helpStmt,
		)),
		goleri.NewOptional(NoGid, rComment),
	)
	helpAccess := goleri.NewKeyword(GidHelpAccess, "access", false)
	helpAlterDatabase := goleri.NewKeyword(GidHelpAlterDatabase, "database", false)
	helpAlterGroup := goleri.NewKeyword(GidHelpAlterGroup, "group", false)
	helpAlterServer := goleri.NewKeyword(GidHelpAlterServer, "server", false)
	helpAlterServers := goleri.NewKeyword(GidHelpAlterServers, "servers", false)
	helpAlterUser := goleri.NewKeyword(GidHelpAlterUser, "user", false)
	helpAlter := goleri.NewSequence(
		GidHelpAlter,
		kAlter,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpAlterDatabase,
			helpAlterGroup,
			helpAlterServer,
			helpAlterServers,
			helpAlterUser,
		)),
	)
	helpCountGroups := goleri.NewKeyword(GidHelpCountGroups, "groups", false)
	helpCountPools := goleri.NewKeyword(GidHelpCountPools, "pools", false)
	helpCountSeries := goleri.NewKeyword(GidHelpCountSeries, "series", false)
	helpCountServers := goleri.NewKeyword(GidHelpCountServers, "servers", false)
	helpCountShards := goleri.NewKeyword(GidHelpCountShards, "shards", false)
	helpCountUsers := goleri.NewKeyword(GidHelpCountUsers, "users", false)
	helpCount := goleri.NewSequence(
		GidHelpCount,
		kCount,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpCountGroups,
			helpCountPools,
			helpCountSeries,
			helpCountServers,
			helpCountShards,
			helpCountUsers,
		)),
	)
	helpCreateGroup := goleri.NewKeyword(GidHelpCreateGroup, "group", false)
	helpCreateUser := goleri.NewKeyword(GidHelpCreateUser, "user", false)
	helpCreate := goleri.NewSequence(
		GidHelpCreate,
		kCreate,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpCreateGroup,
			helpCreateUser,
		)),
	)
	helpDropGroup := goleri.NewKeyword(GidHelpDropGroup, "group", false)
	helpDropSeries := goleri.NewKeyword(GidHelpDropSeries, "series", false)
	helpDropServer := goleri.NewKeyword(GidHelpDropServer, "server", false)
	helpDropShards := goleri.NewKeyword(GidHelpDropShards, "shards", false)
	helpDropUser := goleri.NewKeyword(GidHelpDropUser, "user", false)
	helpDrop := goleri.NewSequence(
		GidHelpDrop,
		kDrop,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpDropGroup,
			helpDropSeries,
			helpDropServer,
			helpDropShards,
			helpDropUser,
		)),
	)
	helpFunctions := goleri.NewKeyword(GidHelpFunctions, "functions", false)
	helpGrant := goleri.NewKeyword(GidHelpGrant, "grant", false)
	helpListGroups := goleri.NewKeyword(GidHelpListGroups, "groups", false)
	helpListPools := goleri.NewKeyword(GidHelpListPools, "pools", false)
	helpListSeries := goleri.NewKeyword(GidHelpListSeries, "series", false)
	helpListServers := goleri.NewKeyword(GidHelpListServers, "servers", false)
	helpListShards := goleri.NewKeyword(GidHelpListShards, "shards", false)
	helpListUsers := goleri.NewKeyword(GidHelpListUsers, "users", false)
	helpList := goleri.NewSequence(
		GidHelpList,
		kList,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpListGroups,
			helpListPools,
			helpListSeries,
			helpListServers,
			helpListShards,
			helpListUsers,
		)),
	)
	helpNoaccess := goleri.NewKeyword(GidHelpNoaccess, "noaccess", false)
	helpRevoke := goleri.NewKeyword(GidHelpRevoke, "revoke", false)
	helpSelect := goleri.NewKeyword(GidHelpSelect, "select", false)
	helpShow := goleri.NewKeyword(GidHelpShow, "show", false)
	helpTimeit := goleri.NewKeyword(GidHelpTimeit, "timeit", false)
	helpTimezones := goleri.NewKeyword(GidHelpTimezones, "timezones", false)
	helpStmt.Set(goleri.NewSequence(
		NoGid,
		kHelp,
		goleri.NewOptional(NoGid, goleri.NewChoice(
			NoGid,
			true,
			helpAccess,
			helpAlter,
			helpCount,
			helpCreate,
			helpDrop,
			helpFunctions,
			helpGrant,
			helpList,
			helpNoaccess,
			helpRevoke,
			helpSelect,
			helpShow,
			helpTimeit,
			helpTimezones,
		)),
	))
	return goleri.NewGrammar(START, regexp.MustCompile(`^[a-z_]+`))
}
