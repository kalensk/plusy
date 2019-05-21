package redisdb

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func (r *Redis) SaveOffset(offset int64) error {
	r.log.Debugf("Saving Telegram offset %d", offset)
	_, err := r.conn.Do("SET", r.plusyOffsetKey(), offset)
	if err != nil {
		// TODO: how to not repeat the command both above nad in the error message. Be sure to fix everywhere...
		return errors.Wrapf(err, "redis command '%s %s %d' failed", "SET", r.plusyOffsetKey(), offset)
	}

	return nil
}

func (r *Redis) getOffset() (int64, error) {
	reply, err := redis.Int64(r.conn.Do("GET", r.plusyOffsetKey()))
	if err == redis.ErrNil {
		return -1, errors.Wrapf(err, "redis command '%s %s' returned nil", "GET", r.plusyOffsetKey())
	}

	if err != nil {
		return -1, errors.Wrapf(err, "redis command '%s %s' failed", "GET", r.plusyOffsetKey())
	}

	if reply != r.lastOffsetPrinted {
		r.log.Debug("Returned telegram offset: ", reply)
		r.lastOffsetPrinted = reply
	}

	return reply, nil
}

func (r *Redis) GetNextOffset() (int64, error) {
	offset, err := r.getOffset()
	if err != nil {
		return -1, errors.Wrapf(err, "redis failed to get next offset")
	}

	return offset + 1, nil
}
