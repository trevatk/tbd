package keyvalue_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/trevatk/tbd/lib/keyvalue"

	pb "github.com/trevatk/tbd/lib/protocol/lsm/v1"
)

type LSMSuite struct {
	suite.Suite
	lsm *keyvalue.LSM
}

func (suite *LSMSuite) SetupSuite() {

	assert := suite.Assert()

	err := os.Mkdir("testfiles", os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		suite.FailNow("failed to create working directory")
	}

	suite.lsm, err = keyvalue.New("testfiles")
	assert.NoError(err)
}

func (suite *LSMSuite) TestPut() {

	assert := suite.Assert()

	keyvalue := &pb.KeyValue{
		Key:   "helloworld",
		Value: []byte("helloworld"),
		Ttl:   -1, // forever
	}
	pbbytes, err := proto.Marshal(keyvalue)
	assert.NoError(err)

	err = suite.lsm.Put("helloworld", pbbytes, nil, -1)
	assert.NoError(err)
}

func (suite *LSMSuite) TestGet() {

	assert := suite.Assert()

	err := suite.lsm.Put("helloworld", []byte("helloworld"), nil, -1)
	assert.NoError(err)

	vbytes, err := suite.lsm.Get("helloworld")
	assert.NoError(err)

	assert.Equal("helloworld", string(vbytes))
}

func (suite *LSMSuite) TearDownSuite() {
	os.RemoveAll("testfiles")
}

func TestLSMSuite(t *testing.T) {
	suite.Run(t, new(LSMSuite))
}
