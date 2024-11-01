package common

import (
	"math/rand"
	"sync"
	"time"

	"encoding/json"
	"strconv"
)

const (
	// Unix time for 2023-01-01 00:00:00 UTC in milliseconds
	epoch    = 1672531200000
	nodeBits = 10
	seqBits  = 12
	maxNode  = 1 << nodeBits
	maxSeq   = 1 << seqBits
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var (
	mutex    sync.Mutex
	sequence int
)

func GenerateSnowflake() Snowflake {
	node := seededRand.Intn(maxNode)

	mutex.Lock()
	defer mutex.Unlock()

	currentTime := time.Now().UnixMilli()
	sequence = (sequence + 1) % maxSeq

	id := int64(0)
	id |= (int64(currentTime-epoch) << (nodeBits + seqBits))
	id |= (int64(node) << seqBits)
	id |= (int64(sequence))

	return Snowflake(id)
}

type Snowflake int64

func (ci Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(ci.String())
}

func (ci *Snowflake) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	var sId Snowflake = Snowflake(decode(s))
	*ci = sId
	return nil
}

func (s Snowflake) String() string {
	return strconv.FormatInt(int64(s), 10)
}

func BuildSnowflake(t string) Snowflake {
	return Snowflake(decode(t))
}

func decode(encodedStr string) int64 {
	i, err := strconv.ParseInt(encodedStr, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}
