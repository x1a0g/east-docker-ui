package API

const LOGIN string = "/eapi/v1/login"
const BASE_UPLOAD string = "/eapi/v1/base/up"
const IMAGE_LIST string = "/eapi/v1/image/list"
const IMAGE_DEL string = "/eapi/v1/image/del"
const IMAGE_PULL string = "/eapi/v1/image/pull"
const IMAGE_INFO string = "/eapi/v1/image/info/:id"
const IMAGE_EXPORT string = "/eapi/v1/image/export"
const IMAGE_IMPORT string = "/eapi/v1/image/import"
const IMAGE_SEARCH string = "/eapi/v1/image/search"
const BASE_INDEX string = "/eapi/v1/base/index"
const BASE_STATIC string = "/eapi/v1/base/static"
const BASE_RESOURCE string = "/eapi/v1/base/resusing"
const LOG_TOP5 string = "/eapi/v1/log/top5"
const CON_LIST = "/eapi/v1/con/list"
const CON_DEL = "/eapi/v1/con/del"
const CON_INFO = "/eapi/v1/con/info/:id"
const CON_CREATE = "/eapi/v1/con/create"
const CON_START = "/eapi/v1/con/start"
const CON_STOP = "/eapi/v1/con/stop"
const CON_RESTART = "/eapi/v1/con/restart"
const CON_PAUSE = "/eapi/v1/con/pause"

// repo
const REPO_LIST = "/eapi/v1/repo/list"
const REPO_DEL = "/eapi/v1/repo/del"
const REPO_CREATE = "/eapi/v1/repo/create"
const REPO_INFO = "/eapi/v1/repo/info"
const REPO_UPDATE = "/eapi/v1/repo/update"
const REPO_DOWN = "/eapi/v1/repo/down"
