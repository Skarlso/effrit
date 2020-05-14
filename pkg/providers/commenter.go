package providers

// Commenter will comments on a PR which is received by the server.
type Commenter interface {
	Comment(owner string, repo string, number int) error
}
