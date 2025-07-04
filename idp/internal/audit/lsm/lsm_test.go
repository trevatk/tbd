package lsm_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/trevatk/tbd/idp/internal/audit/lsm"
	pb "github.com/trevatk/tbd/lib/protocol/lsm/v1"
)

const (
	key   = "hello"
	value = "world"

	testDir = "testfiles"

	neverExpire = -1
)

type LSMSuite struct {
	suite.Suite
	lsm *lsm.LSM
}

func (suite *LSMSuite) SetupSuite() {
	assert := suite.Assert()

	err := os.Mkdir(testDir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		suite.FailNow("failed to create working directory")
	}

	suite.lsm, err = lsm.New(testDir)
	assert.NoError(err)
}

func (suite *LSMSuite) TestPut() {
	assert := suite.Assert()

	keyvalue := &pb.KeyValue{
		Key:   key,
		Value: []byte(value),
		Ttl:   neverExpire,
	}
	pbbytes, err := proto.Marshal(keyvalue)
	assert.NoError(err)

	err = suite.lsm.Put(key, pbbytes, nil, neverExpire)
	assert.NoError(err)
}

func (suite *LSMSuite) TestGet() {
	assert := suite.Assert()

	err := suite.lsm.Put(key, []byte(value), nil, neverExpire)
	assert.NoError(err)

	vbytes, err := suite.lsm.Get(key)
	assert.NoError(err)

	assert.Equal(value, string(vbytes))
}

func (suite *LSMSuite) TearDownSuite() {
	_ = os.RemoveAll(testDir)
}

func TestLSMSuite(t *testing.T) {
	suite.Run(t, new(LSMSuite))
}
