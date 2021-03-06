package git

import (
	"gopkg.in/src-d/go-git.v4/fixtures"
	"gopkg.in/src-d/go-git.v4/plumbing"

	. "gopkg.in/check.v1"
)

type RevListSuite struct {
	BaseSuite
	r *Repository
}

var _ = Suite(&RevListSuite{})

const (
	initialCommit = "b029517f6300c2da0f4b651b8642506cd6aaf45d"
	secondCommit  = "b8e471f58bcbca63b07bda20e428190409c2db47"

	someCommit            = "918c48b83bd081e863dbe1b80f8998f058cd8294"
	someCommitBranch      = "e8d3ffab552895c19b9fcf7aa264d277cde33881"
	someCommitOtherBranch = "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"
)

// Created using: git log --graph --oneline --all
//
// Basic fixture repository commits tree:
//
// * 6ecf0ef vendor stuff
// | * e8d3ffa some code in a branch
// |/
// * 918c48b some code
// * af2d6a6 some json
// *   1669dce Merge branch 'master'
// |\
// | *   a5b8b09 Merge pull request #1
// | |\
// | | * b8e471f Creating changelog
// | |/
// * | 35e8510 binary file
// |/
// * b029517 Initial commit

func (s *RevListSuite) SetUpTest(c *C) {
	r, err := NewFilesystemRepository(fixtures.Basic().One().DotGit().Base())
	c.Assert(err, IsNil)
	s.r = r
}

// ---
// | |\
// | | * b8e471f Creating changelog
// | |/
// * | 35e8510 binary file
// |/
// * b029517 Initial commit
func (s *RevListSuite) TestRevListObjects(c *C) {
	revList := map[string]bool{
		"b8e471f58bcbca63b07bda20e428190409c2db47": true, // second commit
		"c2d30fa8ef288618f65f6eed6e168e0d514886f4": true, // init tree
		"d3ff53e0564a9f87d8e84b6e28e5060e517008aa": true, // CHANGELOG
	}

	initCommit, err := s.r.Commit(plumbing.NewHash(initialCommit))
	c.Assert(err, IsNil)
	secondCommit, err := s.r.Commit(plumbing.NewHash(secondCommit))
	c.Assert(err, IsNil)

	localHist, err := RevListObjects(s.r, []*Commit{initCommit}, nil)
	c.Assert(err, IsNil)

	remoteHist, err := RevListObjects(s.r, []*Commit{secondCommit}, localHist)
	c.Assert(err, IsNil)

	for _, h := range remoteHist {
		c.Assert(revList[h.String()], Equals, true)
	}
	c.Assert(len(remoteHist), Equals, len(revList))
}

func (s *RevListSuite) TestRevListObjectsReverse(c *C) {
	initCommit, err := s.r.Commit(plumbing.NewHash(initialCommit))
	c.Assert(err, IsNil)

	secondCommit, err := s.r.Commit(plumbing.NewHash(secondCommit))
	c.Assert(err, IsNil)

	localHist, err := RevListObjects(s.r, []*Commit{secondCommit}, nil)
	c.Assert(err, IsNil)

	remoteHist, err := RevListObjects(s.r, []*Commit{initCommit}, localHist)
	c.Assert(err, IsNil)

	c.Assert(len(remoteHist), Equals, 0)
}

func (s *RevListSuite) TestRevListObjectsSameCommit(c *C) {
	commit, err := s.r.Commit(plumbing.NewHash(secondCommit))
	c.Assert(err, IsNil)

	localHist, err := RevListObjects(s.r, []*Commit{commit}, nil)
	c.Assert(err, IsNil)

	remoteHist, err := RevListObjects(s.r, []*Commit{commit}, localHist)
	c.Assert(err, IsNil)

	c.Assert(len(remoteHist), Equals, 0)
}

// * 6ecf0ef vendor stuff
// | * e8d3ffa some code in a branch
// |/
// * 918c48b some code
// -----
func (s *RevListSuite) TestRevListObjectsNewBranch(c *C) {
	someCommit, err := s.r.Commit(plumbing.NewHash(someCommit))
	c.Assert(err, IsNil)

	someCommitBranch, err := s.r.Commit(plumbing.NewHash(someCommitBranch))
	c.Assert(err, IsNil)

	someCommitOtherBranch, err := s.r.Commit(plumbing.NewHash(someCommitOtherBranch))
	c.Assert(err, IsNil)

	localHist, err := RevListObjects(s.r, []*Commit{someCommit}, nil)
	c.Assert(err, IsNil)

	remoteHist, err := RevListObjects(
		s.r, []*Commit{someCommitBranch, someCommitOtherBranch}, localHist)
	c.Assert(err, IsNil)

	revList := map[string]bool{
		"a8d315b2b1c615d43042c3a62402b8a54288cf5c": true, // init tree
		"cf4aa3b38974fb7d81f367c0830f7d78d65ab86b": true, // vendor folder
		"9dea2395f5403188298c1dabe8bdafe562c491e3": true, // foo.go
		"e8d3ffab552895c19b9fcf7aa264d277cde33881": true, // branch commit
		"dbd3641b371024f44d0e469a9c8f5457b0660de1": true, // init tree
		"7e59600739c96546163833214c36459e324bad0a": true, // README
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5": true, // otherBranch commit
	}

	for _, h := range remoteHist {
		c.Assert(revList[h.String()], Equals, true)
	}
	c.Assert(len(remoteHist), Equals, len(revList))
}
