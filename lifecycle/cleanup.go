package lifecycle

type RequiresCleanup interface {
	Cleanup()
}
