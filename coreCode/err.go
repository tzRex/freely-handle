package corecode

import "errors"

// 业务常用错误码
var (
	ErrInvalid        = errors.New("invalid.operate")
	ErrColumnTypeFail = errors.New("column.type.fail")
	ErrSearchFail     = errors.New("search.item.error")
	ErrAssert         = errors.New("assert.is.fail")
	ErrExist          = errors.New("data.is.exist")
	ErrRedisNotFind   = errors.New("data.not.find.to.redis")
	ErrRedisSetFail   = errors.New("redis.set.fail")
	ErrIdFail         = errors.New("id.error")
	ErrToken          = errors.New("token.create.error")
	ErrRelationId     = errors.New("relation.id.fail")
	ErrAccess         = errors.New("access.bad")
	ErrNilClass       = errors.New("object.is.nil")
)
